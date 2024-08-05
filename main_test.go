package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var dbTest *sqlx.DB

// TestMain configura o ambiente de teste (banco de dados).
func TestMain(m *testing.M) {
	connStr := "user=admin password=admin dbname=frete sslmode=disable"
	dbTest = sqlx.MustConnect("postgres", connStr)
	defer dbTest.Close()

	// Configurações necessárias no banco de dados para testes
	setupTestDatabase()

	code := m.Run()

	os.Exit(code)
}

// Configuração do banco de dados de teste
func setupTestDatabase() {
	// Criar uma tabela de testes separada
	createTableQuery := `
    CREATE TABLE IF NOT EXISTS quotes_test (
        id SERIAL PRIMARY KEY,
        carrier TEXT,
        service TEXT,
        deadline TEXT,
        price REAL
    );
    `
	dbTest.MustExec(createTableQuery)
}

func TestPostQuoteHandler(t *testing.T) {
	// Dados de entrada simulados para a rota /quote
	requestBody := QuoteRequest{
		Shipper: struct {
			RegisteredNumber string `json:"registered_number"`
			Token            string `json:"token"`
			PlatformCode     string `json:"platform_code"`
		}{
			RegisteredNumber: "25438296000158",
			Token:            "1d52a9b6b78cf07b08586152459a5c90",
			PlatformCode:     "5AKVkHqCn",
		},
		Recipient: struct {
			Type             int    `json:"type"`
			RegisteredNumber string `json:"registered_number"`
			StateInscription string `json:"state_inscription"`
			Country          string `json:"country"`
			Zipcode          int    `json:"zipcode"`
		}{
			Type:             0,
			RegisteredNumber: "25438296000158",
			StateInscription: "12345678",
			Country:          "BR",
			Zipcode:          29161376,
		},
		Dispatchers: []Dispatcher{
			{
				RegisteredNumber: "25438296000158",
				Zipcode:          29161376,
				TotalPrice:       100.0,
				Volumes: []Volume{
					{
						Amount:        1,
						AmountVolumes: 1,
						Category:      "7",
						Sku:           "abc-teste-123",
						Tag:           "tag1",
						Description:   "Descrição do volume",
						Height:        0.2,
						Width:         0.2,
						Length:        0.2,
						UnitaryPrice:  50.0,
						UnitaryWeight: 5.0,
						Consolidate:   false,
						Overlaid:      false,
						Rotate:        false,
					},
				},
			},
		},
		Channel:        "",
		Filter:         0,
		Limit:          0,
		Identification: "",
		Reverse:        false,
		SimulationType: []int{0},
		Returns: struct {
			Composition  bool `json:"composition"`
			Volumes      bool `json:"volumes"`
			AppliedRules bool `json:"applied_rules"`
		}{
			Composition:  false,
			Volumes:      false,
			AppliedRules: false,
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Erro ao codificar JSON: %v", err)
	}

	// Cria uma requisição POST para a rota /quote
	req, err := http.NewRequest("POST", "/quote", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatalf("Erro ao criar requisição: %v", err)
	}

	req.Header.Set("Authorization", "Bearer 1d52a9b6b78cf07b08586152459a5c90")
	req.Header.Set("X-Platform-Code", "5AKVkHqCn")
	req.Header.Set("Content-Type", "application/json")

	// Cria um ResponseRecorder para capturar a resposta
	rr := httptest.NewRecorder()
	handler := PostQuoteHandler(dbTest)
	handler.ServeHTTP(rr, req)

	// Verifica o status da resposta
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status code esperado %v, mas obteve %v", http.StatusOK, status)
	}

	// Verifica se a resposta JSON está correta
	var response map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Erro ao decodificar resposta: %v", err)
	}

	if _, ok := response["carrier"]; !ok {
		t.Errorf("Resposta esperada contém 'carrier', mas não foi encontrado")
	}

	// Insira os dados de teste na tabela de teste
	_, err = dbTest.Exec("INSERT INTO quotes_test (carrier, service, deadline, price) VALUES ($1, $2, $3, $4)",
		"EXPRESSO FR", "Rodoviário", "3", 17.00)
	if err != nil {
		t.Errorf("Erro ao inserir dados de teste na tabela quotes_test: %v", err)
	}
}

func TestGetMetricsHandler(t *testing.T) {
	// Cria uma requisição GET para a rota /metrics
	req, err := http.NewRequest("GET", "/metrics", nil)
	if err != nil {
		t.Fatalf("Erro ao criar requisição: %v", err)
	}

	// Cria um ResponseRecorder para capturar a resposta
	rr := httptest.NewRecorder()
	handler := GetMetricsHandler(dbTest)
	handler.ServeHTTP(rr, req)

	// Verifica o status da resposta
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status code esperado %v, mas obteve %v", http.StatusOK, status)
	}

	// Verifica se a resposta JSON não está vazia
	var metricsResponse map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&metricsResponse); err != nil {
		t.Fatalf("Erro ao decodificar resposta: %v", err)
	}

	if len(metricsResponse) == 0 {
		t.Errorf("Resposta esperada não deve estar vazia")
	}
}
