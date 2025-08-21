//go:build reader
// +build reader

package main

import (
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		content, err := readLogFileLines()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		count, err := readPingCount()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"Ping / Pongs:": count,
			"logs":          content,
		})
	})

	err := router.Run(":" + port)
	if err != nil {
		return
	}
}

func readPingCount() (int64, error) {
	file, err := os.Open("/usr/src/app/files/ping_count.txt")
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return 0, err
	}

	count, err := strconv.ParseInt(string(content), 10, 64)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func readLogFileLines() ([]string, error) {
	file, err := os.Open("/usr/src/app/files/logs.txt")
	if err != nil {
		if os.IsNotExist(err) {
			return []string{"Waiting for logs to be generated..."}, nil
		}
		return nil, err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	var nonEmptyLines []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			nonEmptyLines = append(nonEmptyLines, line)
		}
	}

	return nonEmptyLines, nil
}
