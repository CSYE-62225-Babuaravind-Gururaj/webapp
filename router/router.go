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

	authenticatedAndVerified := router.Group("/")
	authenticatedAndVerified.Use(middleware.BasicAuth(), middleware.UserVerificationMiddleware())
	{
		authenticatedAndVerified.GET("/v2/user/self", routes.GetUserRoute)
		authenticatedAndVerified.PUT("/v2/user/self", routes.UpdateUserRoute)
	}

	router.POST("/v2/user", routes.CreateUserRoute)

	router.NoRoute(middleware.HandleNoRoute)

	router.GET("/v2/user/verify", routes.VerifyUserRoute)

	return router
}
