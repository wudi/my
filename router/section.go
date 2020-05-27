package router

import (
	"github.com/gin-gonic/gin"
	"my/api"
	"my/middleware"
)

func InitSectionRouter(Router *gin.RouterGroup) {
	ApiRouter := Router.Group("section").Use(middleware.JWTAuth())
	{
		ApiRouter.GET("/section/list", api.Sections)
	}
}