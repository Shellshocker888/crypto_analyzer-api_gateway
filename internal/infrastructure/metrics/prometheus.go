package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

var (
	httpRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"path", "status"},
	)

	httpDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: prometheus.DefBuckets,
		}, []string{"path"},
	)

	limitedRequests = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "rate_limited_requests",
			Help: "Total number of rate-limited requests",
		},
	)

	// Глобальный registry
	Registry = prometheus.NewRegistry()
)

// Инициализация — один раз при старте приложения
func InitMetrics() {
	Registry.MustRegister(httpRequests, httpDuration, limitedRequests)
}

// Инкремент запросов
func IncRequest(path string, status int) {
	httpRequests.WithLabelValues(path, strconv.Itoa(status)).Inc()
}

// Наблюдение длительности
func ObserveDuration(path string, d time.Duration) {
	httpDuration.WithLabelValues(path).Observe(d.Seconds())
}

// Инкремент rate-limited
func IncRateLimited() {
	limitedRequests.Inc()
}
