package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func CorsMiddleware(c *gin.Context) {
	var allowedOrigins = os.Getenv("ALLOWED_ORIGINS")

	c.Header("Access-Control-Allow-Origin", allowedOrigins)
	c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
	c.Header("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.88 Safari/537.36")

	if c.Request.Method == http.MethodOptions {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}

	c.Next()
}
