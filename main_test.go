package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

func TestPostQuoteHandler(t *testing.T) {
	// Configurar a variável de ambiente
	os.Setenv("DATABASE_URL", "postgres://usuario:usuario@url:porta/bdName?sslmode=disable") //por a url na hora de testar do dotenv

	// Abrir conexão com o banco de dados
	var err error
	connStr := os.Getenv("DATABASE_URL")
	db, err = sqlx.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	// Criar tabela se não existir
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS quotes (
		id SERIAL PRIMARY KEY,
		carrier TEXT,
		service TEXT,
		deadline TEXT,
		price REAL
	);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		t.Fatalf("Error creating table: %v", err)
	}

	// criando dados fakes mokkados (adicionar dados reais, removi para evitar expor no projeto público)
	quoteRequest := QuoteRequest{
		Shipper: struct {
			RegisteredNumber string "json:\"registered_number\""
			Token            string "json:\"token\""
			PlatformCode     string "json:\"platform_code\""
		}{
			RegisteredNumber: "123456789",
			Token:            "token123",
			PlatformCode:     "platform123",
		},
		Recipient: struct {
			Type             int    "json:\"type\""
			RegisteredNumber string "json:\"registered_number\""
			StateInscription string "json:\"state_inscription\""
			Country          string "json:\"country\""
			Zipcode          int    "json:\"zipcode\""
		}{
			Type:             1,
			RegisteredNumber: "987654321",
			StateInscription: "insc123",
			Country:          "BR",
			Zipcode:          12345678,
		},
		Dispatchers: []Dispatcher{
			{
				RegisteredNumber: "123456789",
				Zipcode:          12345678,
				TotalPrice:       100.0,
				Volumes: []Volume{
					{
						Amount:       1,
						Sku:          "sku123",
						Description:  "desc123",
						Height:       10.0,
						Width:        10.0,
						Length:       10.0,
						UnitaryPrice: 50.0,
						UnitaryWeight: 1.0,
					},
				},
			},
		},
		Channel:        "channel123",
		Filter:         1,
		Limit:          1,
		Identification: "id123",
		Reverse:        false,
		SimulationType: []int{1},
		Returns: struct {
			Composition  bool "json:\"composition\""
			Volumes      bool "json:\"volumes\""
			AppliedRules bool "json:\"applied_rules\""
		}{
			Composition:  true,
			Volumes:      true,
			AppliedRules: true,
		},
	}

	reqBody, err := json.Marshal(quoteRequest)
	if err != nil {
		t.Fatalf("Error marshaling request body: %v", err)
	}

	req, err := http.NewRequest("POST", "/quote", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer token123")
	req.Header.Set("X-Platform-Code", "platform123")

	rr := httptest.NewRecorder()
	handler := PostQuoteHandler(db) // Pasando o DB como argumento

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var response QuoteResponse
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Error decoding response body: %v", err)
	}

	//verificações adicioneis
	if len(response.Dispatchers) == 0 {
		t.Errorf("Expected dispatchers in response but got none")
	}
}
