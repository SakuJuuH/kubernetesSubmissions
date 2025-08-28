package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var (
	port   = os.Getenv("PORT")
	dbName = os.Getenv("POSTGRES_DB")
	dbHost = os.Getenv("POSTGRES_HOST")
	dbUser = os.Getenv("POSTGRES_USER")
	dbPass = os.Getenv("POSTGRES_PASSWORD")
	db     *sql.DB
)

const (
	maxRetries    = 10
	retryInterval = 5 * time.Second
)

func initDB() {
	if dbHost == "" {
		fmt.Println("DB_HOST environment variable not set")
		os.Exit(1)
	}

	connStr := fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable", dbHost, dbUser, dbPass, dbName)

	var err error
	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("postgres", connStr)
		if err == nil && db.Ping() == nil {
			break
		}
		log.Printf("Failed to connect or ping database (attempt %d/%d): %v\n", i+1, maxRetries, err)
		time.Sleep(retryInterval)
	}
	if err != nil || db.Ping() != nil {
		log.Fatal("Failed to connect to database after retries:", err)
	}

	// Create table if not exists
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS ping_counter (
			id SERIAL PRIMARY KEY,
			count BIGINT NOT NULL DEFAULT 0
		)
	`)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}

	// Initialize counter if not exists
	_, err = db.Exec(`
		INSERT INTO ping_counter (count) 
		SELECT 0 
		WHERE NOT EXISTS (SELECT 1 FROM ping_counter)
	`)
	if err != nil {
		log.Fatal("Failed to initialize counter:", err)
	}
}

func incrementCounter() (int64, error) {
	var count int64
	err := db.QueryRow(`
		UPDATE ping_counter 
		SET count = count + 1 
		WHERE id = (SELECT MIN(id) FROM ping_counter)
		RETURNING count
	`).Scan(&count)
	return count, err
}

func getCount() (int64, error) {
	var count int64
	err := db.QueryRow(`SELECT count FROM ping_counter WHERE id = (SELECT MIN(id) FROM ping_counter)`).Scan(&count)
	return count, err
}

func main() {
	if port == "" {
		port = "3001"
	}

	initDB()
	defer db.Close()

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":     "Welcome to the Ping Pong Service!",
			"status_code": http.StatusOK,
		})
	})

	router.GET("/ping", func(c *gin.Context) {
		count, err := incrementCounter()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to increment count"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"pong": count,
		})
	})

	router.GET("/pings", func(c *gin.Context) {
		count, err := getCount()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get count"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"pongs": count,
		})
	})

	router.GET("/db-health", func(c *gin.Context) {
		if err := db.Ping(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database is not reachable"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Database is reachable"})
	})

	err := router.Run(":" + port)

	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		os.Exit(1)
	}
}
