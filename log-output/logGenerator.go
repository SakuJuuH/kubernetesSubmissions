//go:build generator
// +build generator

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
)

var randomString string = uuid.New().String()

func main() {

	file, err := os.OpenFile("/usr/src/app/files/logs.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	for {
		timestamp := time.Now().Format(time.RFC3339)
		logLine := fmt.Sprintf("%s: %s\n", timestamp, randomString)
		if _, err := file.WriteString(logLine); err != nil {
			fmt.Printf("Error writing to file: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(logLine)

		time.Sleep(5 * time.Second)
	}
}

func getCurrentTimestamp() string {
	return time.Now().Format(time.RFC3339)
}
