package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var randomString = uuid.New().String()
var timestamp string

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		timestamp = getCurrentTimestamp()
		c.JSON(http.StatusOK, gin.H{
			"timestamp":    timestamp,
			"randomString": randomString,
		})
	})

	err := router.Run(":" + port)
	if err != nil {
		return
	}

	for {
		timestamp := getCurrentTimestamp()
		fmt.Printf("%s: %s\n", timestamp, randomString)
		time.Sleep(5 * time.Second)
	}

}

func getCurrentTimestamp() string {
	return time.Now().Format(time.RFC3339)
}
