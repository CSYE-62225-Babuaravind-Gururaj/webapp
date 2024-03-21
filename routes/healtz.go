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

		if err := sqlDB.Ping(); err != nil {
			c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
			c.Status(http.StatusServiceUnavailable)
			return
		}

		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Status(http.StatusOK)

		// Log the incoming request
		logger.Info().
			Str("method", "GET").
			Str("path", "/healthz").
			Int("status", 200).
			Msg("Health check successful")
	}
}
