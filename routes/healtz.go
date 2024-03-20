package routes

import (
	"io"
	"log/slog"
	"net/http"

	"cloud-proj/health-check/logs"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap/exp/zapslog"
	"gorm.io/gorm"
)

func RouteHealthz(db *gorm.DB) gin.HandlerFunc {

	logger := logs.CreateLogger()

	defer logger.Sync()

	sl := slog.New(zapslog.NewHandler(logger.Core(), nil))

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

		//Checks DB connectivity (testing connection using the ORM)
		if err := sqlDB.Ping(); err != nil {
			c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
			c.Status(http.StatusServiceUnavailable)
			return
		}

		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Status(http.StatusOK)
		sl.Info(
			"incoming request",
			slog.String("method", "GET"),
			slog.String("path", "/healthz"),
			slog.Int("status", 200),
		)
	}
}
