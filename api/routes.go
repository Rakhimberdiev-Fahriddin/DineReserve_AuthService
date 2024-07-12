package api

import (
	_ "auth-service/api/docs"
	"auth-service/api/handler"
	"auth-service/api/middleware"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

// @title Auth Service API
// @version 1.0
// @description This is a sample server for Auth Service.
// @host localhost:8081
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @schemes http
func Routes(handle *handler.Handler) *gin.Engine {
	router := gin.Default()

	// Swagger endpointini sozlash
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Use(middleware.LoggerMiddleware())

	router.POST("auth/register", handle.RegisterHandler)
	router.POST("auth/login", handle.LoginHandler)
	router.GET("auth/refresh_token", handle.RefreshToken)

	return router
}
