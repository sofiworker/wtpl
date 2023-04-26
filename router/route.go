package router

import (
	"github.com/gin-gonic/gin"
	"wtpl/api"
	v1 "wtpl/api/v1"
	"wtpl/conf"
	"wtpl/middleware"
)

func NewRoute() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	switch conf.GetConfig().App.Mode {
	case "prod":
	case "dev":
		gin.SetMode(gin.DebugMode)
	}
	engine := gin.New()

	engine.Use(middleware.TokenLimiterHandler(), middleware.LoggerHandler([]string{}...))

	engine.GET("/ping", api.Ping)

	userGroup := engine.Group("/user")
	{
		userGroup.POST("/sign_in", middleware.WrapJsonHandler(v1.SignIn))
	}
	return engine
}
