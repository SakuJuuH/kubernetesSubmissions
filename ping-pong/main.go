package main

import (
	"fmt"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

func main() {
	var port = os.Getenv("PORT")

	if port == "" {
		port = "3001"
	}

	router := gin.Default()

	var pongCounter int64 = 0

	router.GET("/pingpong", func(c *gin.Context) {
		count := atomic.AddInt64(&pongCounter, 1)

		c.JSON(http.StatusOK, gin.H{
			"pong": count,
		})
	})

	router.GET("/pings", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"pongs": pongCounter,
		})
	})

	err := router.Run(":" + port)

	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		os.Exit(1)
	}
}
