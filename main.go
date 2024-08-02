package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Shipper struct {
	RegisteredNumber string `json:"registered_number"`
	Token            string `json:"token"`
	PlatformCode     string `json:"platform_code"`
}

type Recipient struct {
	Type              int    `json:"type"`
	RegisteredNumber  string `json:"registered_number"`
	StateInscription  string `json:"state_inscription"`
	Country           string `json:"country"`
	Zipcode           int    `json:"zipcode"`
}

type Volume struct {
	Amount           int     `json:"amount"`
	AmountVolumes    int     `json:"amount_volumes"`
	Category         string  `json:"category"`
	Sku              string  `json:"sku"`
	Tag              string  `json:"tag"`
	Description      string  `json:"description"`
	Height           float64 `json:"height"`
	Width            float64 `json:"width"`
	Length           float64 `json:"length"`
	UnitaryPrice     float64 `json:"unitary_price"`
	UnitaryWeight    float64 `json:"unitary_weight"`
	Consolidate      bool    `json:"consolidate"`
	Overlaid         bool    `json:"overlaid"`
	Rotate           bool    `json:"rotate"`
}

type Dispatcher struct {
	RegisteredNumber string   `json:"registered_number"`
	Zipcode          int      `json:"zipcode"`
	TotalPrice       float64  `json:"total_price"`
	Volumes          []Volume `json:"volumes"`
}

type Returns struct {
	Composition   bool `json:"composition"`
	Volumes       bool `json:"volumes"`
	AppliedRules  bool `json:"applied_rules"`
}

type QuoteRequest struct {
	Shipper          Shipper      `json:"shipper"`
	Recipient        Recipient    `json:"recipient"`
	Dispatchers      []Dispatcher `json:"dispatchers"`
	Channel          string       `json:"channel"`
	Filter           int          `json:"filter"`
	Limit            int          `json:"limit"`
	Identification   string       `json:"identification"`
	Reverse          bool         `json:"reverse"`
	SimulationType   []int        `json:"simulation_type"`
	Returns          Returns      `json:"returns"`
}

func main() {
	http.HandleFunc("/quote", PostQuoteHandler())
	fmt.Println("Servidor rodando na porta 8080")
	http.ListenAndServe(":8080", nil)
}

func PostQuoteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}

		var requestData QuoteRequest
		if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
			http.Error(w, "Erro ao decodificar o JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Verifique se o zipcode não está vazio
		if requestData.Recipient.Zipcode == 0 {
			http.Error(w, "Recipient.Zipcode cannot be empty", http.StatusBadRequest)
			return
		}

		// Obtendo valores dos cabeçalhos da requisição original
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

		// Chamada à API externa
		url := "https://sp.freterapido.com/api/v3/quote/simulate"
		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}

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

		var quotes interface{}
		if err := json.NewDecoder(resp.Body).Decode(&quotes); err != nil {
			http.Error(w, "Erro ao decodificar a resposta da API: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(quotes)
	}
}
