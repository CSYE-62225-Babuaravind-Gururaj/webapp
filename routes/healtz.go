package routes

import (
	"io"
	"net/http"
	"sync"
	"time"

	"cloud-proj/health-check/logs"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// rate limiter structure
var healthzRateLimiter = struct {
	sync.RWMutex
	hitCounts map[string]int
	lastReset time.Time
}{
	hitCounts: make(map[string]int),
	lastReset: time.Now(),
}

func rateLimitWarning(logger zerolog.Logger, c *gin.Context) {
	healthzRateLimiter.Lock()
	defer healthzRateLimiter.Unlock()

	// Reset counts every minute
	if time.Since(healthzRateLimiter.lastReset) > time.Minute {
		healthzRateLimiter.hitCounts = make(map[string]int)
		healthzRateLimiter.lastReset = time.Now()
	}

	// Increment the hit count for the /healthz endpoint
	healthzRateLimiter.hitCounts["/healthz"]++

	// Log a warning if the /healthz endpoint is hit more than 5 times in a minute
	if healthzRateLimiter.hitCounts["/healthz"] > 5 {
		logger.Warn().
			Str("method", "GET").
			Str("path", "/healthz").
			Int("hitCount", healthzRateLimiter.hitCounts["/healthz"]).
			Msg("Rate limit warning: /healthz endpoint is being hit too frequently")
	}
}

func RouteHealthz(db *gorm.DB) gin.HandlerFunc {
	logger := logs.CreateLogger()

	return func(c *gin.Context) {
		// Add rate limit check and potential warning log
		rateLimitWarning(logger, c)

		if reqBody, err := io.ReadAll(c.Request.Body); err != nil || len(reqBody) != 0 || c.Request.URL.RawQuery != "" {
			c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
			c.Status(http.StatusBadRequest)
			return
		}

		sqlDB, err := db.DB()
		if err != nil {
			c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
			c.Status(http.StatusBadRequest)
			return
		}

		if err := sqlDB.Ping(); err != nil {
			c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
			c.Status(http.StatusServiceUnavailable)
			return
		}

		c.Header("Content-Type", "text/plain")
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.String(http.StatusOK, "CI/CD check done")

		// Log the successful health check
		logger.Info().
			Str("method", "GET").
			Str("path", "/healthz").
			Int("status", 200).
			Msg("Health check successful")
	}
}
