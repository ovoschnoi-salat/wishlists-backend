package middlewares

import (
	"net/netip"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var DefaultMetricPath = "/metrics"

type Prometheus struct {
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewPrometheus() *Prometheus {
	factory := promauto.With(prometheus.DefaultRegisterer)
	p := &Prometheus{
		requestCounter: factory.NewCounterVec(
			prometheus.CounterOpts{
				Name: "api_requests_total",
				Help: "Total number of HTTP requests processed by the API.",
			},
			[]string{"method", "path", "status"},
		),
		requestDuration: factory.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "api_request_duration_seconds",
				Help:    "Histogram of response latencies for API requests.",
				Buckets: []float64{.01, .03, .05, .1, .3},
			},
			[]string{"method", "path"},
		),
	}
	return p
}

func (p *Prometheus) Middleware(c *gin.Context) {
	if c.FullPath() == DefaultMetricPath {
		if IsPrivate(c.RemoteIP()) {
			c.Next()
		} else {
			c.AbortWithStatus(404)
		}
		return
	}
	start := time.Now()

	c.Next()

	elapsed := float64(time.Since(start)) / float64(time.Second)
	status := strconv.Itoa(c.Writer.Status())
	url, method := getUrlAndMethodFromContext(c)

	p.requestCounter.WithLabelValues(method, url, status).Inc()
	p.requestDuration.WithLabelValues(method, url).Observe(elapsed)
}

func getUrlAndMethodFromContext(c *gin.Context) (url string, method string) {
	method = c.Request.Method
	path := c.FullPath()
	if path == "" {
		if method == "OPTIONS" {
			return "cors", method
		}
		return "unknown", method
	}
	return path, method
}

func (p *Prometheus) RegisterGinPrometheusHandler(e *gin.Engine) {
	h := promhttp.Handler()
	e.GET(DefaultMetricPath, gin.WrapH(h))
}

func IsPrivate(ip string) bool {
	ipAddr, err := netip.ParseAddr(ip)
	if err != nil {
		return false
	}
	return ipAddr.IsPrivate()
}
