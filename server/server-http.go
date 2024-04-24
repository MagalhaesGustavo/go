package main

import (
	"context"
	"database/sql"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"net/http"
	"time"
)

type RequestResponse struct {
	Usdbrl struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}

func main() {
	http.HandleFunc("/cotacao", cotacaoDolar)
	log.Println("Servidor iniciado em http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("Erro ao servidor:", err)
		return
	}
}

func cotacaoDolar(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	req1, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		log.Println("Erro ao criar requisição:", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	req, err := http.DefaultClient.Do(req1)
	if err != nil {
		log.Println("Erro ao fazer solicitação HTTP:", err)
		http.Error(w, "Erro ao fazer solicitação HTTP", http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Println("Erro ao ler resposta HTTP:", err)
		http.Error(w, "Erro ao ler resposta HTTP", http.StatusInternalServerError)
		return
	}

	log.Println("Resposta recebida:", string(body))

	var data RequestResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println("Erro ao decodificar JSON:", err)
		http.Error(w, "Erro ao decodificar JSON", http.StatusInternalServerError)
		return
	}

	result, err := json.Marshal(data.Usdbrl)
	if err != nil {
		log.Println("Erro ao codificar resposta JSON:", err)
		http.Error(w, "Erro ao codificar resposta JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if _, err := w.Write(result); err != nil {
		log.Println("Erro ao escrever resposta:", err)
	}

	db, err := sql.Open("sqlite3", "cotacao.db")
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS cotacao (
                        id INTEGER PRIMARY KEY,
                        bid TEXT)`)
	if err != nil {
		log.Println(err)
		return
	}

	ctxDB, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err = db.ExecContext(ctxDB,"INSERT INTO cotacao (bid) VALUES (?)", data.Usdbrl.Bid)
	if err != nil {
		log.Println(err)
		return
	}
}
