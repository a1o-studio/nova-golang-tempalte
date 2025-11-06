package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	db "github.com/a1ostudio/nova/db/sqlc"
	_ "github.com/a1ostudio/nova/docs"
	"github.com/a1ostudio/nova/internal/config"
	"github.com/a1ostudio/nova/internal/logger"
	"github.com/a1ostudio/nova/internal/server"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

//	@title			nova
//	@description	API documentation for nova Project
//	@host			local.a1o.studio:4000
//	@BasePath		/nova

func main() {
	config := mustLoadConfig()
	logger.NewLogger(config.Env)
	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()

	connPool := mustConnectDB(ctx, config.DBSource)
	runDBMigration(config.MigrationURL, config.DBSource)

	redisClient := newRedisClient(config)
	store := db.NewStore(connPool)

	server := mustNewServer(config, store, redisClient)
	startHTTPServer(server)

	waitForShutdown(ctx, server)
}

func mustLoadConfig() config.Config {
	config, err := config.LoadConfig(".")
	if err != nil {
		logger.L().Fatal("cannot load config", zap.Error(err))
	}
	return config
}

func mustConnectDB(ctx context.Context, dbSource string) *pgxpool.Pool {
	connPool, err := pgxpool.New(ctx, dbSource)
	if err != nil {
		logger.L().Fatal("cannot connect to db", zap.Error(err))
	}
	return connPool
}

func newRedisClient(config config.Config) *redis.Client {
	redisAddr := fmt.Sprintf("0.0.0.0:%d", config.RedisPort)
	return redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: config.RedisPassword,
	})
}

func mustNewServer(config config.Config, store db.Store, redisClient *redis.Client) *server.Server {
	server, err := server.NewServer(config, store, redisClient)
	if err != nil {
		logger.L().Fatal("cannot create server", zap.Error(err))
	}
	return server
}

func startHTTPServer(server *server.Server) {
	go func() {
		if err := server.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.L().Fatal("HTTP server failed", zap.Error(err))
		}
	}()
}

func waitForShutdown(ctx context.Context, server *server.Server) {
	<-ctx.Done()
	logger.L().Info("Shutting down gracefully...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.L().Error("HTTP server shutdown failed", zap.Error(err))
		} else {
			logger.L().Info("HTTP server stopped gracefully")
		}
	}()
	wg.Wait()
	logger.L().Info("Application stopped")
}

func runDBMigration(migrationURL, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		logger.L().Warn("migration setup failed, skipping", zap.Error(err))
		return
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		logger.L().Warn("migration failed, skipping", zap.Error(err))
		return
	}

	logger.L().Info("database migration completed")
}
