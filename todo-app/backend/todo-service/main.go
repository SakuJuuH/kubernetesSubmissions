package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Todo struct {
	ID   int    `json:"id" db:"id"`
	Task string `json:"task" db:"task"`
	Done bool   `json:"done" db:"done"`
}

var (
	port             = os.Getenv("PORT")
	allowedOrigins   = os.Getenv("ALLOWED_ORIGINS")
	dbHost           = os.Getenv("POSTGRES_HOST")
	dbUser           = os.Getenv("POSTGRES_USER")
	dbPass           = os.Getenv("POSTGRES_PASSWORD")
	dbName           = os.Getenv("POSTGRES_DB")
	db               *sqlx.DB
	randomArticleURL = os.Getenv("RANDOM_ARTICLE_URL")
)

func getTodos() ([]Todo, error) {
	todos := make([]Todo, 0)
	err := db.Select(&todos, "SELECT id, task, done FROM todos")
	return todos, err
}

func addTodo(task string) (Todo, error) {
	var todo Todo
	err := db.Get(&todo, "INSERT INTO todos (task) VALUES ($1) RETURNING id, task, done", task)
	if err != nil {
		return Todo{}, err
	}
	return todo, nil
}

func main() {
	if port == "" {
		fmt.Println("$PORT must be set")
		os.Exit(1)
	}

	if allowedOrigins == "" {
		allowedOrigins = "*"
	}

	initDB()
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}(db)

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
		todos, err := getTodos()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, todos)
	})

	router.POST("/api/todos/random", func(c *gin.Context) {
		if randomArticleURL == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "RANDOM_ARTICLE_URL is not set"})
			return
		}

		resp, err := http.Get(randomArticleURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		articleURL := resp.Request.URL.String()
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Printf("Error closing response body: %v", err)
			}
		}(resp.Body)

		task := fmt.Sprintf("Read: %s", articleURL)

		createdTodo, err := addTodo(task)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"New todo created": createdTodo,
		})
	})

	router.POST("/api/todos", func(c *gin.Context) {
		var requestTodo struct {
			Task string `json:"task" binding:"required"`
		}

		if err := c.BindJSON(&requestTodo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		newTodo, err := addTodo(requestTodo.Task)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, newTodo)
	})

	fmt.Printf("Server starting on port %s\n", port)
	err := router.Run(":" + port)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		os.Exit(1)
	}
}

func initDB() {
	if dbHost == "" || dbName == "" || dbUser == "" || dbPass == "" {
		fmt.Println("Database configuration environment variables must be set")
		os.Exit(1)
	}

	connStr := fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable", dbHost, dbUser, dbPass, dbName)

	var err error
	db, err = sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS todos (
			id SERIAL PRIMARY KEY,
			task TEXT NOT NULL,
			done BOOLEAN NOT NULL DEFAULT FALSE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}
}
