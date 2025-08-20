package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router := gin.Default()

	fmt.Printf("Server started in port %s\n", port)

	router.Run(":" + port).Error()
}
