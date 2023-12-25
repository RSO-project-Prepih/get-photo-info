package prometheus

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	DBConnectionAttempts = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "db_connection_attempts",
			Help: "Number of attempts to connect to the database",
		},
	)

	DBConnectionErrors = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "db_connection_errors",
			Help: "Number of errors when connecting to the database",
		},
	)

	HTTPRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path"},
	)
)

func GetMetrics() http.Handler {
	log.Println("Starting the prometheus metrics...")
	prometheus.MustRegister(DBConnectionAttempts)
	prometheus.MustRegister(DBConnectionErrors)
	prometheus.MustRegister(HTTPRequestDuration)

	return promhttp.Handler()
}
