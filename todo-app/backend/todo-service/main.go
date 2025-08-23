package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Todo struct {
	ID   string `json:"id"`
	Task string `json:"task"`
	Done bool   `json:"done"`
}

var (
	port           = os.Getenv("PORT")
	allowedOrigins = os.Getenv("ALLOWED_ORIGINS")
)

var todos = []Todo{}

func main() {
	if port == "" {
		fmt.Println("$PORT must be set")
		os.Exit(1)
	}

	if allowedOrigins == "" {
		allowedOrigins = "*"
	}

	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", allowedOrigins)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	router.GET("/api/todos", func(c *gin.Context) {
		c.JSON(http.StatusOK, todos)
	})

	router.POST("/api/todos", func(c *gin.Context) {
		var newTodo Todo

		if err := c.BindJSON(&newTodo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		newTodo.ID = uuid.New().String()
		newTodo.Done = false

		fmt.Printf("%+v\n", newTodo)
		todos = append(todos, newTodo)
		c.JSON(http.StatusCreated, newTodo)
	})

	fmt.Printf("Server starting on port %s\n", port)
	err := router.Run(":" + port)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		os.Exit(1)
	}
}
