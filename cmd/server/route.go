package server

import (
	"context"
	"habbit-tracker/consts"
	auth_controller "habbit-tracker/internal/controller/auth"
	auth_usecase "habbit-tracker/internal/usecase/auth"
	"habbit-tracker/middlewares"
	"habbit-tracker/pkg/logger"
	"habbit-tracker/pkg/metrics"

	doc "habbit-tracker/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap/zapcore"
)

func InitHandler() *gin.Engine {
	// Init gin
	r := gin.Default()
	//init all dependencies
	metrics, err := metrics.NewMetrics(consts.MetricServiceName)
	if err != nil {
		logger.Fatal(context.Background(), "Fail to start metrics %v", zapcore.Field{Key: "error", Type: zapcore.StringType, String: err.Error()})
	}

	metrics.SetSkipPath([]string{
		"/swagger/*any",
		"/metrics",
		"/",
	})
	r.GET("/metrics", func(c *gin.Context) {
		metrics.GinServeHandler(c)
	})

	r.Use(metrics.GinMetricsHttpMiddleware())
	// Add middlewares
	r.Use(middlewares.RecoverMiddleware())
	r.Use(middlewares.LogMiddleware())
	r.Use(middlewares.ResponseMiddleware())
	r.Use(middlewares.MiddlewareTracing())
	r.Use(middlewares.CORSMiddleware())

	// Add routes

	// Health check
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})

	InitAuthRoute(r)

	// Swagger
	doc.SwaggerInfo.Title = "Habbit Tracker API"
	doc.SwaggerInfo.BasePath = "/"

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return r
}

func InitAuthRoute(r *gin.Engine) {

	authController := auth_controller.NewController(auth_usecase.AuthUsecase())

	authGroup := r.Group("/auth")
	authGroup.POST("/signup", authController.SignUp)
	authGroup.POST("/signin", authController.SignIn)
}
