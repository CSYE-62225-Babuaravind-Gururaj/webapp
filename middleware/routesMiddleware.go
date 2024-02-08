package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckMethodAndPath(c *gin.Context) {
	if c.Request.Method != "GET" && c.Request.URL.Path == "/healthz" {
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.AbortWithStatus(http.StatusMethodNotAllowed)
		return
	}
	// if c.Request.Method == "POST" && c.Request.URL.Path != "/v1/user" {
	// 	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	// 	c.AbortWithStatus(http.StatusNotFound)
	// 	return
	// }
	c.Next()
}

func HandleNoRoute(c *gin.Context) {
	log.Printf("Request Method: %s, Path: %s", c.Request.Method, c.Request.URL.Path)
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Status(http.StatusNotFound)
}
