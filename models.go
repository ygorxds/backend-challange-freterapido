package main

type QuoteRequest struct {
	Recipient struct {
		Address struct {
			Zipcode string `json:"zipcode"`
		} `json:"address"`
	} `json:"recipient"`
	Volumes []struct {
		Category      int     `json:"category"`
		Amount        int     `json:"amount"`
		UnitaryWeight float64 `json:"unitary_weight"`
		Price         float64 `json:"price"`
		Sku           string  `json:"sku"`
		Height        float64 `json:"height"`
		Width         float64 `json:"width"`
		Length        float64 `json:"length"`
	} `json:"volumes"`
}

type QuoteResponse struct {
	Carrier []struct {
		Name     string  `json:"name"`
		Service  string  `json:"service"`
		Deadline int     `json:"deadline"`
		Price    float64 `json:"price"`
	} `json:"carrier"`
}
