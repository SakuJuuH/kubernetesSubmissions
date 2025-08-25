package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

type Todo struct {
	ID   int    `json:"id" db:"id"`
	Task string `json:"task" db:"task"`
	Done bool   `json:"done" db:"done"`
}

func main() {
	var port = os.Getenv("PORT")

	if port == "" {
		log.Fatal().Msg("$PORT must be set")
	}

	db := initDB()
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close database connection")
		}
	}(db)

	repo := NewTodoRepository(db)
	controller := NewTodosController(repo)

	router := gin.Default()

	router.Use(CorsMiddleware)

	router.GET("/", controller.welcome)

	router.GET("/api/todos", controller.getTodos)

	router.POST("/api/todos", controller.createTodo)

	router.POST("/api/todos/random", controller.createRandomTodo)

	log.Info().Str("port", port).Msg("Server starting")
	err := router.Run(":" + port)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}

func initDB() *sqlx.DB {
	var (
		dbHost = os.Getenv("POSTGRES_HOST")
		dbUser = os.Getenv("POSTGRES_USER")
		dbPass = os.Getenv("POSTGRES_PASSWORD")
		dbName = os.Getenv("POSTGRES_DB")
	)

	if dbHost == "" || dbName == "" || dbUser == "" || dbPass == "" {
		log.Fatal().Msg("All database environment variables not set")
	}

	connStr := fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable", dbHost, dbUser, dbPass, dbName)

	var err error
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database:")
	}

	if err = db.Ping(); err != nil {
		log.Fatal().Err(err).Msg("Failed to ping database:")
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
		log.Fatal().Err(err).Msg("Failed to create table:")
	}

	return db
}
