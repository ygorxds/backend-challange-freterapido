package main

import (
    "encoding/json"
    "net/http"
    "github.com/jmoiron/sqlx"
)

func PostQuoteHandler(db *sqlx.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
            return
        }

        var input map[string]interface{}
        if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
            http.Error(w, "Erro ao decodificar o JSON", http.StatusBadRequest)
            return
        }

        // Implementar a lógica de chamada da API Frete Rápido e gravar no banco de dados

        response := map[string]interface{}{
            "carrier": []map[string]interface{}{
                {"name": "EXPRESSO FR", "service": "Rodoviário", "deadline": "3", "price": 17},
                {"name": "Correios", "service": "SEDEX", "deadline": "1", "price": 20.99},
            },
        }
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
    }
}

func GetMetricsHandler(db *sqlx.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        lastQuotes := r.URL.Query().Get("last_quotes")

        // Implementar a lógica para buscar as métricas no banco de dados

        metrics := map[string]interface{}{
            "carrier_count":     2,
            "total_price":       37.99,
            "average_price":     18.99,
            "cheapest_shipping": 17,
            "most_expensive_shipping": 20.99,
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(metrics)
    }
}
