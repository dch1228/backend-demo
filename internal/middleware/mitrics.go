package middleware

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/duchenhao/backend-demo/internal/conf"
)

var (
	// 处理中的请求数
	apiRequestsInFlight prometheus.Gauge

	// 请求量
	apiRequestTotal *prometheus.CounterVec

	// 请求
	apiRequestHistogram *prometheus.HistogramVec
)

func init() {
	apiRequestsInFlight = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "api_request_in_flight",
			Help: "gauge of requests currently being served",
		},
	)

	apiRequestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_request_total",
			Help: "api request counter",
		},
		[]string{"path", "code", "method"},
	)

	apiRequestHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "api_request_duration_milliseconds",
			Help:    "api request duration ms histogram",
			Buckets: []float64{100, 200, 300, 500, 1000, 3000},
		},
		[]string{"path", "code", "method"},
	)

	prometheus.MustRegister(
		apiRequestsInFlight,
		apiRequestTotal,
		apiRequestHistogram,
	)
}

func RequestMetrics(g *gin.Engine) gin.HandlerFunc {
	metricsAccounts := gin.Accounts{conf.Core.MetricsUsername: conf.Core.MetricsPassword}
	g.GET("/metrics", gin.BasicAuth(metricsAccounts), gin.WrapH(promhttp.Handler()))
	return func(ctx *gin.Context) {
		now := time.Now()
		apiRequestsInFlight.Inc()
		defer apiRequestsInFlight.Dec()
		ctx.Next()

		code := strconv.Itoa(ctx.Writer.Status())
		method := strings.ToLower(ctx.Request.Method)
		path := ctx.FullPath()

		apiRequestTotal.WithLabelValues(path, code, method).Inc()
		duration := time.Since(now).Nanoseconds() / int64(time.Millisecond)
		apiRequestHistogram.WithLabelValues(path, code, method).Observe(float64(duration))
	}
}
