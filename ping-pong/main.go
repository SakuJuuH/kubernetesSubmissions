package main

import (
	"fmt"
	"os"
	"strconv"
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

		err := writeCountToFile(count)
		if err != nil {
			fmt.Printf("Error writing count to file: %v\n", err)
		}

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

func writeCountToFile(count int64) error {
	err := os.MkdirAll("/usr/src/app/files", 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	file, err := os.OpenFile("/usr/src/app/files/ping_count.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(strconv.FormatInt(count, 10))
	if err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}

	return nil
}
