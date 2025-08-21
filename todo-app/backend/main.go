package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Todo App API",
			"status":  http.StatusOK,
		})
	})

	fmt.Printf("Server started in port %s\n", port)

	router.Run(":" + port).Error()
}
