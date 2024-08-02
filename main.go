package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"crypto/tls"

	_ "github.com/lib/pq"
)

// Estruturas para a API de Quote
type Volume struct {
	Amount             int     `json:"amount"`
	AmountVolumes      int     `json:"amount_volumes"`
	Category           string  `json:"category"`
	Sku                string  `json:"sku"`
	Tag                string  `json:"tag"`
	Description        string  `json:"description"`
	Height             float64 `json:"height"`
	Width              float64 `json:"width"`
	Length             float64 `json:"length"`
	UnitaryPrice       float64 `json:"unitary_price"`
	UnitaryWeight      float64 `json:"unitary_weight"`
	Consolidate        bool    `json:"consolidate"`
	Overlaid           bool    `json:"overlaid"`
	Rotate             bool    `json:"rotate"`
}

type Dispatcher struct {
	RegisteredNumber  string  `json:"registered_number"`
	Zipcode           int     `json:"zipcode"`
	TotalPrice        float64 `json:"total_price"`
	Volumes           []Volume `json:"volumes"`
}

type QuoteRequest struct {
	Shipper struct {
		RegisteredNumber string `json:"registered_number"`
		Token            string `json:"token"`
		PlatformCode     string `json:"platform_code"`
	} `json:"shipper"`
	Recipient struct {
		Type              int    `json:"type"`
		RegisteredNumber  string `json:"registered_number"`
		StateInscription  string `json:"state_inscription"`
		Country           string `json:"country"`
		Zipcode           int    `json:"zipcode"`
	} `json:"recipient"`
	Dispatchers       []Dispatcher `json:"dispatchers"`
	Channel           string       `json:"channel"`
	Filter            int          `json:"filter"`
	Limit             int          `json:"limit"`
	Identification    string       `json:"identification"`
	Reverse           bool         `json:"reverse"`
	SimulationType    []int        `json:"simulation_type"`
	Returns           struct {
		Composition      bool `json:"composition"`
		Volumes          bool `json:"volumes"`
		AppliedRules     bool `json:"applied_rules"`
	} `json:"returns"`
}

type Carrier struct {
	Carrier        string  `json:"carrier"`
	Service        string  `json:"service"`
	Deadline       string  `json:"deadline"`
	Price          float64 `json:"price"`
}

type QuoteResponse struct {
	Dispatchers []struct {
		ID                  string `json:"id"`
		Offers              []struct {
			Carrier struct {
				CompanyName        string `json:"company_name"`
				Logo               string `json:"logo"`
				Name               string `json:"name"`
				Reference          int    `json:"reference"`
				RegisteredNumber   string `json:"registered_number"`
				StateInscription   string `json:"state_inscription"`
			} `json:"carrier"`
			CarrierOriginalDeliveryTime struct {
				Days           int    `json:"days"`
				EstimatedDate  string `json:"estimated_date"`
			} `json:"carrier_original_delivery_time"`
			CostPrice                   float64 `json:"cost_price"`
			DeliveryTime                struct {
				Days           int    `json:"days"`
				EstimatedDate  string `json:"estimated_date"`
			} `json:"delivery_time"`
			Expiration                  string `json:"expiration"`
			FinalPrice                  float64 `json:"final_price"`
			HomeDelivery                bool   `json:"home_delivery"`
			Identifier                  string `json:"identifier"`
			Modal                       string `json:"modal"`
			Offer                       int    `json:"offer"`
			OriginalDeliveryTime        struct {
				Days           int    `json:"days"`
				EstimatedDate  string `json:"estimated_date"`
			} `json:"original_delivery_time"`
			Service                     string `json:"service"`
			ServiceCode                 string `json:"service_code"`
			SimulationType              int    `json:"simulation_type"`
			Weights                     struct {
				Real float64 `json:"real"`
			} `json:"weights"`
		} `json:"offers"`
		RegisteredNumberDispatcher string  `json:"registered_number_dispatcher"`
		RegisteredNumberShipper    string  `json:"registered_number_shipper"`
		RequestID                  string  `json:"request_id"`
		TotalPrice                 float64 `json:"total_price"`
		ZipcodeOrigin              int     `json:"zipcode_origin"`
	} `json:"dispatchers"`
}

var db *sql.DB

func main() {
	var err error

	// Abrir conexão com o banco de dados
	connStr := os.Getenv("DATABASE_URL")
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
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
	);
	`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		fmt.Println("Error creating table:", err)
		return
	}

	// Definir rotas
	http.HandleFunc("/quote", PostQuoteHandler())
	http.HandleFunc("/metrics", GetMetricsHandler())

	fmt.Println("Servidor rodando na porta 8080")
	http.ListenAndServe(":8080", nil)
}

func PostQuoteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}

		// Decodificar dados da requisição
		var requestData QuoteRequest
		if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
			http.Error(w, "Erro ao decodificar o JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		authHeader := r.Header.Get("Authorization")
		platformCodeHeader := r.Header.Get("X-Platform-Code")

		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusBadRequest)
			return
		}
		if platformCodeHeader == "" {
			http.Error(w, "X-Platform-Code header is required", http.StatusBadRequest)
			return
		}

		url := "https://sp.freterapido.com/api/v3/quote/simulate"

		// Criar cliente HTTP com configuração TLS personalizada
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // Ignorar validação do certificado (não para produção)
			},
		}
		client := &http.Client{Transport: tr}

		reqBody, err := json.Marshal(requestData)
		if err != nil {
			http.Error(w, "Erro ao codificar o JSON para a requisição: "+err.Error(), http.StatusInternalServerError)
			return
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
		if err != nil {
			http.Error(w, "Erro ao criar a requisição: "+err.Error(), http.StatusInternalServerError)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", authHeader)
		req.Header.Set("X-Platform-Code", platformCodeHeader)

		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "Erro ao fazer a solicitação: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := ioutil.ReadAll(resp.Body)
			http.Error(w, fmt.Sprintf("Status da resposta não OK: %d, resposta: %s", resp.StatusCode, body), http.StatusInternalServerError)
			return
		}

		var apiResponse QuoteResponse
		if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
			http.Error(w, "Erro ao decodificar a resposta da API: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Criar estrutura para armazenar a resposta no formato desejado
		var carriers []Carrier
		for _, dispatcher := range apiResponse.Dispatchers {
			for _, offer := range dispatcher.Offers {
				carriers = append(carriers, Carrier{
					Carrier:     offer.Carrier.Name,
					Service:  offer.Service,
					Deadline: strconv.Itoa(offer.DeliveryTime.Days),
					Price:    offer.FinalPrice,
				})

				// Salvar cotações no banco de dados
				_, err = db.Exec("INSERT INTO quotes (carrier, service, deadline, price) VALUES ($1, $2, $3, $4)",
					offer.Carrier.Name, offer.Service, strconv.Itoa(offer.DeliveryTime.Days), offer.FinalPrice)
				if err != nil {
					http.Error(w, "Erro ao salvar dados no banco de dados: "+err.Error(), http.StatusInternalServerError)
					return
				}
			}
		}

		response := map[string]interface{}{
			"carrier": carriers,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func GetMetricsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}

		rows, err := db.Query("SELECT carrier, service, deadline, price FROM quotes")
		if err != nil {
			http.Error(w, "Erro ao consultar o banco de dados: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var quotes []map[string]interface{}
		for rows.Next() {
			var carrier, service, deadline string
			var price float64
			if err := rows.Scan(&carrier, &service, &deadline, &price); err != nil {
				http.Error(w, "Erro ao ler dados do banco de dados: "+err.Error(), http.StatusInternalServerError)
				return
			}
			quote := map[string]interface{}{
				"carrier":    carrier,
				"service":  service,
				"deadline": deadline,
				"price":    price,
			}
			quotes = append(quotes, quote)
		}

		if err := rows.Err(); err != nil {
			http.Error(w, "Erro ao iterar sobre resultados: "+err.Error(), http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"quotes": quotes,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
