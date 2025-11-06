package server

import (
	"github.com/a1ostudio/nova/internal/model"
	"github.com/a1ostudio/nova/internal/pkg/resp"

	"github.com/gin-gonic/gin"
)

const version = "1.0.0"

// Healthcheck
//
//	@Summary		状态检测
//	@Description	检测当前 API 服务状态
//	@Tags			Common
//	@Success		200	{object}	resp.Result[model.Healthcheck]	"返回状态信息"
//	@Router			/v1/healthcheck [get]
func (server Server) healthcheck(ctx *gin.Context) {
	data := &model.Healthcheck{
		Status: "available",
		System: model.System{
			Environment: server.config.Env,
			Version:     version,
		},
	}

	resp.Success(ctx, data)
}
