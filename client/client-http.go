package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type ServerResponse struct {
	USDBRL struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	var data ServerResponse
	if err := json.Unmarshal(body, &data.USDBRL); err != nil {
		log.Fatalf("Error decoding JSON: %v", err)
	}

	file, err := os.Create("cotacao.txt")
	if err != nil {
		log.Fatalf("Error writing to file: %v", err)
	}
	defer file.Close()

	file.WriteString("Dólar: " + data.USDBRL.Bid)

	log.Println("Arquivo criado com sucesso!")
}
