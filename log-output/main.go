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

var PingPongURL = os.Getenv("PING_PONG_URL")
var randomString = uuid.New().String()

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	if PingPongURL == "" {
		PingPongURL = "http://localhost:3001/"
	}

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		count, err := getPingCount()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		timestamp := time.Now().Format(time.RFC3339)

		c.JSON(http.StatusOK, gin.H{
			timestamp:      randomString,
			"Ping / Pongs": count,
		})
	})

	err := router.Run(":" + port)
	if err != nil {
		return
	}
}

func getPingCount() (int64, error) {
	u, err := url.Parse(PingPongURL)
	if err != nil {
		return 0, fmt.Errorf("unable to parse PING_PONG_URL: %w", err)
	}
	u.Path = "pings"

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
