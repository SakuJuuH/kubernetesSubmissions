package main

import (
	"fmt"
	"os"
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

func main() {
	var port = os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	router := gin.Default()

	var pongCounter int64 = 0

	router.GET("/pingpong", func(c *gin.Context) {
		count := atomic.AddInt64(&pongCounter, 1)
		c.JSON(200, gin.H{
			"pong": count,
		})
	})

	err := router.Run(":" + port)

	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		os.Exit(1)
	}

}
