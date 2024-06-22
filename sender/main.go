package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	url := "http://localhost:5000/query"

	for range 10000000 {
		format := fmt.Sprintf(`{"key":"%d"}`, rand.Intn(1000000000000))
		jsonData := []byte(format)
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Fatalf("Failed to send POST request: %v", err)
		}

		defer resp.Body.Close()

		time.Sleep(200 * time.Millisecond)
	}
}
