package metrics

import (
	"context"
	"habbit-tracker/consts"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	skipPaths = map[string]bool{
		"/metrics":     true,
		"/favicon.ico": true,
	}

	ProcessDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "app",
		Subsystem: "kafka",
		Name:      "kafka_msg_consuming_duration_ms",
		Help:      "Kafka consumer message processing duration in ms",
	}, []string{"service_name", "topic"})

	ErrorCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "app",
		Subsystem: "kafka",
		Name:      "kafka_msg_consuming_error_count",
		Help:      "Kafka consuming errors count",
	}, []string{"service_name", "topic", "err_type"})

	SuccessCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "app",
		Subsystem: "kafka",
		Name:      "kafka_msg_consuming_success_count",
		Help:      "Kafka consuming success count",
	}, []string{"service_name", "topic"})
)

// Prometheus Metrics struct
type PrometheusMetrics struct {
	ServiceName                    string
	ByCid                          bool
	RequestTotal                   prometheus.Counter
	RequestsTotalByPath            *prometheus.CounterVec
	RequestsTotalByApplicationCode *prometheus.CounterVec
	RequestDuration                *prometheus.HistogramVec

	ProcessDuration *prometheus.HistogramVec
	SuccessCount    *prometheus.CounterVec
	ErrorCount      *prometheus.CounterVec
}

func NewMetrics(serviceName string) (IMetrics, error) {
	metr := PrometheusMetrics{
		ServiceName: serviceName,
		RequestTotal: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "http_request_total",
			},
		),
		RequestsTotalByPath: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
			},
			[]string{"service_name", "http_status", "method", "path", "cid"},
		),
		RequestsTotalByApplicationCode: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_by_application_code_total",
			},
			[]string{"service_name", "method", "path", "resp_code", "resp_code_msg", "cid"},
		),
		RequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "http_request_duration_seconds",
			},
			[]string{"service_name", "http_status", "method", "path"},
		),
		ProcessDuration: ProcessDuration,
		SuccessCount:    SuccessCount,
		ErrorCount:      ErrorCount,
	}

	if err := prometheus.Register(collectors.NewBuildInfoCollector()); err != nil {
		return nil, err
	}

	prometheus.MustRegister(
		metr.RequestTotal,
		metr.RequestsTotalByPath,
		metr.RequestDuration,
		metr.RequestsTotalByApplicationCode,
		ProcessDuration,
		ErrorCount,
		SuccessCount,
	)

	return &metr, nil

}

func (metr *PrometheusMetrics) WrapHandler() http.Handler {
	return promhttp.Handler()
}

func (metr *PrometheusMetrics) GinServeHandler(r *gin.Context) {

	promhttp.Handler().ServeHTTP(r.Writer, r.Request)
}

func (metr *PrometheusMetrics) IncHits(httpStatus int, method, path string, cids ...string) {
	cid := ""
	if len(cids) > 0 {
		cid = cids[0]
	}

	metr.RequestTotal.Inc()
	metr.RequestsTotalByPath.WithLabelValues(metr.ServiceName, strconv.Itoa(httpStatus), method, path, cid).Inc()
}

func (metr *PrometheusMetrics) IncRequestByAppCode(respCode int, respCodeMsg string, method, path string, cids ...string) {

	cid := ""
	if len(cids) > 0 {
		cid = cids[0]
	}

	metr.RequestsTotalByApplicationCode.
		WithLabelValues(metr.ServiceName, method, path, strconv.Itoa(respCode), respCodeMsg, cid).Inc()

}

func (metr *PrometheusMetrics) ObserveResponseTime(status int, method, path string, observeTime float64) {
	metr.RequestDuration.WithLabelValues(metr.ServiceName, strconv.Itoa(status), method, path).Observe(observeTime)
}

func (metr *PrometheusMetrics) isSkipPath(path string) bool {
	return skipPaths[path]
}

func (metr *PrometheusMetrics) SetSkipPath(paths []string) {
	for _, val := range paths {
		skipPaths[val] = true
	}
}

func (metr *PrometheusMetrics) UnSetSkipPath(paths []string) {
	for _, val := range paths {
		delete(skipPaths, val)
	}
}

func (s *PrometheusMetrics) GinMetricsHttpMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		// before request
		t := time.Now()
		c.Next()

		// after request
		endpoint := c.Request.URL.Path
		statusCode := c.Writer.Status()
		if statusCode == http.StatusNotFound {
			endpoint = c.Request.URL.Path
		}

		if !s.isSkipPath(endpoint) {

			cid := getCidFromReq(c)
			method := c.Request.Method
			latency := time.Since(t)

			s.ObserveResponseTime(statusCode, method, endpoint, latency.Seconds())
			s.IncHits(statusCode, method, endpoint, cid)

			//tracking resp code
			if err, ok := c.Get(consts.GinErrorKey); ok {
				err := err.(consts.CustomError)
				s.IncRequestByAppCode(err.Code, err.Message, method, endpoint, cid)
				return
			}

			if _, ok := c.Get(consts.GinResponseKey); ok {
				s.IncRequestByAppCode(0, "Success", method, endpoint, cid)
				return
			}

		}
	}
}

func (s *PrometheusMetrics) GrpcMetricsInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		startTime := time.Now()
		// Call the handler to process the request
		resp, err = handler(ctx, req)
		var path string
		if info != nil {
			path = info.FullMethod
		}

		latency := time.Since(startTime).Seconds()

		var statusCode int
		if err != nil {
			st, ok := status.FromError(err)
			if ok {
				statusCode = int(st.Code())
			} else {
				statusCode = -1
			}
		} else {
			statusCode = int(codes.OK)
		}

		s.ObserveResponseTime(statusCode, "gRPC", path, latency)
		s.IncHits(statusCode, "gRPC", path, "")

		return resp, err
	}
}

func (metr PrometheusMetrics) logConsumingDuration(topic string, observeTime float64) {
	metr.ProcessDuration.WithLabelValues(metr.ServiceName, topic).Observe(observeTime)
}

func (metr PrometheusMetrics) TopicLogSuccess(topic string) {
	metr.SuccessCount.WithLabelValues(metr.ServiceName, topic).Inc()
}

func (metr PrometheusMetrics) TopicLogError(topic string, errorType string) {
	metr.ErrorCount.WithLabelValues(metr.ServiceName, topic, errorType).Inc()
}

func (metr PrometheusMetrics) TopicKafkaEventHandler(topic string, err error, observeTime float64) {

	metr.logConsumingDuration(topic, observeTime)

	if err != nil {
		metr.TopicLogError(topic, "failed")
		return
	}
	metr.TopicLogSuccess(topic)
}
