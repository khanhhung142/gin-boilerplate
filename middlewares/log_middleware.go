package middlewares

import (
	"habbit-tracker/consts"
	"habbit-tracker/pkg/idutil"
	"habbit-tracker/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LogMiddleware is a middleware to log the request
// When error occurs, we set the error to the gin context, so that we can log the error in this middleware
func LogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Next()
		latency := time.Since(t)
		var err error
		if errStr, ok := c.Get(consts.GinErrorKey); ok {
			castErr, isErr := errStr.(consts.CustomError)
			if !isErr {
				castErr = consts.CodeErrorUnknown
			}
			err = castErr
		}
		if c.Request.RequestURI != "/" {
			logger.Info(c, "request",
				zap.String("traceId", c.GetString(consts.TraceKey)),
				zap.String("time", t.Format("2006/01/02 - 15:04:05 -07")),
				zap.String("method", c.Request.Method),
				zap.String("uri", c.Request.RequestURI),
				zap.String("latency", latency.String()),
				zap.String("ip", c.ClientIP()),
				zap.Any("body", c.Request.Body),
				zap.Error(err),
			)
		}
	}
}

// MiddlewareTracing is a middleware to add traceId to the context
// This traceId will be used to trace the request in the whole system

func MiddlewareTracing() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceId := idutil.ULIDNow()
		// SetCtx for this app
		c.Set(consts.TraceKey, traceId)

		// SetCtx for calling gRPC
		// c.SetCtx(metadata.NewIncomingContext(c.GetCtx(), metadata.Pairs(tracekeyGrpc, traceId)))
		c.Next()
	}
}
