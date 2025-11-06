package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	db "github.com/a1ostudio/nova/db/sqlc"
	"github.com/a1ostudio/nova/internal/config"
	"github.com/a1ostudio/nova/internal/controller"
	"github.com/a1ostudio/nova/internal/logger"
	"github.com/a1ostudio/nova/internal/middleware"
	"github.com/a1ostudio/nova/internal/pkg/resp"
	"github.com/a1ostudio/nova/internal/pkg/token"
	"github.com/a1ostudio/nova/internal/pkg/validation"

	docs "github.com/a1ostudio/nova/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"go.uber.org/zap"
)

type Server struct {
	config      config.Config
	store       db.Store
	tokenMaker  token.Maker
	router      *gin.Engine
	redis       *redis.Client
	controllers []controller.RegisterRoutes
	httpServer  *http.Server // 保存 HTTP 服务器引用
}

func NewServer(config config.Config, store db.Store, redis *redis.Client) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	// Register controllers

	server := &Server{
		config:      config,
		store:       store,
		tokenMaker:  tokenMaker,
		redis:       redis,
		controllers: []controller.RegisterRoutes{},
	}

	// 注册 validation
	validation.NewValidation()

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.New()

	router.NoRoute(resp.WrapNotFoundError())

	router.HandleMethodNotAllowed = true
	router.NoMethod(resp.WrapMethodNotAllowedError())

	// CORS middleware
	cors := cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool { // 允许的前端地址
			u, err := url.Parse(origin)
			if err != nil {
				return false
			}
			host := u.Hostname() // 只取主机名，不含端口
			return strings.HasSuffix(host, server.config.Domain)
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // 允许的方法
		AllowHeaders:     []string{"Origin", "Content-Type"},                  // 允许的请求头
		ExposeHeaders:    []string{"Content-Length"},                          // 允许前端获取的响应头
		AllowCredentials: true,                                                // 允许携带 Cookie
		MaxAge:           12 * time.Hour,
	})

	router.MaxMultipartMemory = 8 << 20 // 8MB 限制上传文件大小
	router.Use(cors)

	// middlewares
	router.Use(logger.LoggerMiddleware())
	router.Use(middleware.RecoverPanic())

	go middleware.CleanupClients(1*time.Minute, 5*time.Minute)
	// Rate limiting middleware
	router.Use(middleware.RateLimitByIPMiddleware(server.config.LimitRate, server.config.LimitBurst))

	if server.config.Env != config.Dev {
		// 5s 超时
		router.Use(middleware.Timeout(5 * time.Second))
	}

	docs.SwaggerInfo.Version = "v1.0.0"

	nova := router.Group("nova")
	{
		v1 := nova.Group("v1")
		{
			v1.GET("healthcheck", server.healthcheck)

			for _, controller := range server.controllers {
				controller.RegisterRoutes(v1)
			}
		}

		if server.config.Env != config.Prod {
			router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		}
	}

	server.router = router
}

func (server *Server) Start() error {
	addr := fmt.Sprintf(":%d", server.config.Port)
	server.httpServer = &http.Server{
		Addr:    addr,
		Handler: server.router,
	}

	logger.L().Info("Starting server...", zap.String("addr", addr))

	err := server.httpServer.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

// Shutdown 优雅关闭 HTTP 服务器
func (server *Server) Shutdown(ctx context.Context) error {
	if server.httpServer == nil {
		return nil
	}

	logger.L().Info("Shutting down HTTP server...")
	return server.httpServer.Shutdown(ctx)
}

// Deprecated: 请使用 Start 和 Shutdown 方法组合替代本方法。
func (server *Server) StartWithGracefulShutdown() {
	addr := fmt.Sprintf(":%d", server.config.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: server.router,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		// Notify 监听 SIGINT 和 SIGTERM 信号
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		logger.L().Info("Shutdown server...", zap.String("signal", s.String()))

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		shutdownError <- srv.Shutdown(ctx)
	}()

	logger.L().Info("Starting server...", zap.String("addr", addr))

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		logger.L().Fatal("ListenAndServe", zap.Error(err))
		return
	}

	err = <-shutdownError
	if err != nil {
		logger.L().Fatal("Server shutdown failed", zap.Error(err))
		return
	}

	logger.L().Info("Server shutdown successfully")
}
