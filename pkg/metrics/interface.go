package metrics

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

// App Metrics interface
type IMetrics interface {

	// IncHits increments the number of hits for a given status code, method, path and cid
	IncHits(httpStatus int, method, path string, cid ...string)
	// IncRequestByAppCode increments the number of requests for a given response code, method, path and cid
	IncRequestByAppCode(respCode int, respCodeMsg string, method, path string, cid ...string)
	// IncRequestByAppCode increments the number of requests for a given response code, method, path and cid
	ObserveResponseTime(status int, method, path string, observeTime float64)
	// SetSkipPath sets the paths to skip metrics
	SetSkipPath(paths []string)
	// UnSetSkipPath unsets the paths to skip metrics
	UnSetSkipPath(paths []string)
	// WrapHandler wraps the http handler
	WrapHandler() http.Handler
	// GinServeHandler wraps the http handler
	GinServeHandler(r *gin.Context)
	// GfMetricsHttpMiddleware returns the middleware for metrics
	GinMetricsHttpMiddleware() gin.HandlerFunc
	// GrpcMetricsInterceptor returns the middleware for metrics
	GrpcMetricsInterceptor() grpc.UnaryServerInterceptor

	// for kafka metrics
	TopicLogError(topic string, errorType string)
	TopicLogSuccess(topic string)
	TopicKafkaEventHandler(topic string, err error, observeTime float64)
}
