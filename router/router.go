// router/router.go
package router

import (
	"cloud-proj/health-check/middleware"
	"cloud-proj/health-check/routes"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RouterSetup(db *gorm.DB) *gin.Engine {
	router := gin.Default()
	router.Use(middleware.CheckMethodAndPath)

	router.GET("/healthz", routes.RouteHealthz(db))

	router.GET("/v1/user/self", middleware.BasicAuth(), routes.GetUserRoute)

	router.PUT("/v1/user/self", middleware.BasicAuth(), routes.UpdateUserRoute)

	router.POST("/v1/user", routes.CreateUserRoute)

	router.NoRoute(middleware.HandleNoRoute)

	return router
}
