package controller

import "github.com/gin-gonic/gin"

type RegisterRoutes interface {
	RegisterRoutes(router *gin.RouterGroup)
}
