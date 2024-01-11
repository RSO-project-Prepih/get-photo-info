package health

import (
	"log"
	"net/http"
	"time"

	"github.com/RSO-project-Prepih/get-photo-info/database"
	"github.com/RSO-project-Prepih/get-photo-info/prometheus"
	"github.com/heptiolabs/healthcheck"
	_ "github.com/lib/pq"
)

// GetHealth checks godoc
// @Summary Get get health checks
// @Description get the health checks of the service
// @Tags health
// @Accept  json
// @Produce  json
// @Success 200 {object} string
// @Router /live [get]
// @Router /ready [get]
func HealthCheckHandler() (http.HandlerFunc, http.HandlerFunc) {
	log.Println("Starting the health check...")
	// Create a health instance.
	health := healthcheck.NewHandler()

	health.AddLivenessCheck("goroutine-threshold", healthcheck.GoroutineCountCheck(100))

	// Add a readiness check for a database.
	health.AddReadinessCheck("database", func() error {
		startTiem := time.Now()
		db := database.NewDBConnection()
		defer db.Close()

		err := db.Ping()
		duration := time.Since(startTiem)
		prometheus.HTTPRequestDuration.WithLabelValues("database").Observe(duration.Seconds())

		return err
	})

	return health.LiveEndpoint, health.ReadyEndpoint
}
