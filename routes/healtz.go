package routes

import (
	"io"
	"net/http"

	"cloud-proj/health-check/logs"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RouteHealthz(db *gorm.DB) gin.HandlerFunc {
	logger := logs.CreateLogger()

	// Note: No need for logger.Sync() with zerolog in typical usage, as it does not buffer.

	return func(c *gin.Context) {
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

		// Checks DB connectivity (testing connection using the ORM)
		if err := sqlDB.Ping(); err != nil {
			c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
			c.Status(http.StatusServiceUnavailable)
			return
		}

		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Status(http.StatusOK)

		// Log the incoming request using zerolog
		logger.Info().
			Str("method", "GET").
			Str("path", "/healthz").
			Int("status", 200).
			Msg("incoming request")
	}
}
