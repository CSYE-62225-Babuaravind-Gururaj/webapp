package main

import (
	"cloud-proj/health-check/config"
	"cloud-proj/health-check/database"
	"cloud-proj/health-check/logs"
	"cloud-proj/health-check/routes"

	"cloud-proj/health-check/middleware"

	"github.com/gin-gonic/gin"
)

func main() {

	logger := logs.CreateLogger()

	config.LoadEnv()
	database.InitDB()

	defer logger.Sync() //flush the logger

	logger.Info("Hello from Zap!")

	router := gin.Default()
	router.Use(middleware.CheckMethodAndPath)

	router.GET("/healthz", routes.RouteHealthz(database.DB))

	router.GET("/v1/user/self", middleware.BasicAuth(), routes.GetUserRoute)

	router.PUT("/v1/user/self", middleware.BasicAuth(), routes.UpdateUserRoute)

	router.POST("/v1/user", routes.CreateUserRoute)

	router.NoRoute(middleware.HandleNoRoute)

	router.Run(":8080")

}
