package main

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

var randomString = uuid.New().String()
var timestamp string

func main() {
	for {
		timestamp = time.Now().Format(time.RFC3339)
		fmt.Printf("%s: %s\n", timestamp, randomString)
		time.Sleep(5 * time.Second)
	}
}
