package main

import (
	"cloud-proj/health-check/config"
	"cloud-proj/health-check/database"
	"cloud-proj/health-check/logs" // Ensure this is updated for zerolog.
	"cloud-proj/health-check/middleware"
	"cloud-proj/health-check/routes"
	"cloud-proj/health-check/utils"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Setup zerolog global logger to use console writer for human-friendly output.
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})

	// Assuming CreateLogger now returns a *zerolog.Logger tailored for your application's needs.

	logger := logs.CreateLogger()

	config.LoadEnv()

	database.InitDB()

	// No need for logger.Sync() as in zap, zerolog does not buffer by default.

	logger.Info().Msg("Hello from Zerolog!")

	utils.InitPubSubClient()

	router := gin.Default()

	// If middleware and routes need the logger, consider passing it as an argument or using a context.
	router.Use(middleware.CheckMethodAndPath)

	router.GET("/healthz", routes.RouteHealthz(database.DB))

	authenticatedAndVerified := router.Group("/")
	authenticatedAndVerified.Use(middleware.BasicAuth(), middleware.UserVerificationMiddleware())
	{
		authenticatedAndVerified.GET("/v1/user/self", routes.GetUserRoute)
		authenticatedAndVerified.PUT("/v1/user/self", routes.UpdateUserRoute)
	}

	router.POST("/v1/user", routes.CreateUserRoute)

	router.NoRoute(middleware.HandleNoRoute)

	router.GET("/v1/user/verify", routes.VerifyUserRoute)

	router.GET("/v1/user/check", routes.GetUserRoute)

	router.Run(":8080")
}
