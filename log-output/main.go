package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PingPongResp struct {
	Pongs int64 `json:"pongs"`
}

var (
	port        = os.Getenv("PORT")
	pingPongURL = os.Getenv("PING_PONG_URL")
	message     = os.Getenv("MESSAGE")
)

var randomString = uuid.New().String()

func main() {
	if port == "" {
		port = "3000"
	}

	if pingPongURL == "" {
		pingPongURL = "http://localhost:3001/api"
	}

	if message == "" {
		message = "env variable not set"
	}

	fileContent, err := os.ReadFile("/usr/src/app/files/information.txt")
	var fileContentStr string
	if err != nil {
		fileContentStr = "file not found"
	} else {
		fileContentStr = string(fileContent)
	}

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":     "Welcome to the Log Output Service!",
			"status_code": http.StatusOK,
		})
	})

	router.GET("/log", func(c *gin.Context) {
		count, err := getPingCount()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		timestamp := time.Now().Format(time.RFC3339)

		c.JSON(http.StatusOK, gin.H{
			"file content": fileContentStr,
			"env variable": fmt.Sprintf("MESSAGE=%s", message),
			timestamp:      randomString,
			"Ping / Pongs": count,
		})
	})

	err = router.Run(":" + port)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		os.Exit(1)
	}
}

func getPingCount() (int64, error) {
	u, err := url.Parse(pingPongURL)
	if err != nil {
		return 0, fmt.Errorf("unable to parse PING_PONG_URL: %w", err)
	}
	u.Path = "/pings"

	resp, err := http.Get(u.String())
	if err != nil {
		return 0, fmt.Errorf("unable to get ping count: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("unable to read response body: %w", err)
	}

	var pingPongResp PingPongResp

	if err := json.Unmarshal(body, &pingPongResp); err != nil {
		return 0, fmt.Errorf("unable to unmarshal response: %w", err)
	}

	return pingPongResp.Pongs, nil
}
