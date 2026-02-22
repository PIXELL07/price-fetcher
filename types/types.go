package types

import "time"

type PriceFetchRequest struct {
	Ticker string `json:"ticker"`
}

type PriceResponse struct {
	Ticker    string    `json:"ticker"`
	Price     float64   `json:"price"`
	Currency  string    `json:"currency"` // e.g. "USD"
	Timestamp time.Time `json:"timestamp"`
}

type MarketStats struct {
	Ticker     string    `json:"ticker"`
	Price      float64   `json:"price"`
	MarketCap  float64   `json:"market_cap"`
	Volume24h  float64   `json:"volume_24h"`
	Change24h  float64   `json:"change_24h_pct"` // percentage e.g. -2.5
	High24h    float64   `json:"high_24h"`
	Low24h     float64   `json:"low_24h"`
	CircSupply float64   `json:"circulating_supply"`
	Timestamp  time.Time `json:"timestamp"`
}

type CoinInfo struct {
	Ticker      string   `json:"ticker"`
	Name        string   `json:"name"`
	Category    string   `json:"category"` // e.g. "Layer 1", "DeFi", "Stablecoin"
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Website     string   `json:"website"`
}

type SupportedTickersResponse struct {
	Tickers []string `json:"tickers"`
	Count   int      `json:"count"`
}

type BatchPriceRequest struct {
	Tickers []string `json:"tickers"`
}

type BatchPriceResponse struct {
	Prices    []PriceResponse `json:"prices"`
	Timestamp time.Time       `json:"timestamp"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
