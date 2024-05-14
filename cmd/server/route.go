package server

import (
	auth_controller "gin-boilerplate/internal/controller/auth"
	auth_usecase "gin-boilerplate/internal/usecase/auth"
	"gin-boilerplate/middlewares"

	doc "gin-boilerplate/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitHandler() *gin.Engine {
	// Init gin
	r := gin.Default()

	// Add middlewares
	r.Use(middlewares.LogMiddleware())
	r.Use(middlewares.CORSMiddleware())
	r.Use(middlewares.ResponseMiddleware())

	// Add routes

	// Health check
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})

	authController := auth_controller.NewController(auth_usecase.AuthUsecase())

	authGroup := r.Group("/auth")
	authGroup.POST("/signup", authController.SignUp)
	authGroup.POST("/signin", authController.SignIn)

	// Swagger
	doc.SwaggerInfo.Title = "gin-boilerplate API"
	doc.SwaggerInfo.BasePath = "/"

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return r
}
