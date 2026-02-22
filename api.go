package main

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"

	"github.com/PIXELL07/price-fetcher/types"
)

type JSONAPIServer struct {
	listenAddr string
	svc        PriceFetcher
}

func NewJSONAPIServer(listenAddr string, svc PriceFetcher) *JSONAPIServer {
	return &JSONAPIServer{
		listenAddr: listenAddr,
		svc:        svc,
	}
}

type APIFunc func(context.Context, http.ResponseWriter, *http.Request) error

func (s *JSONAPIServer) Run() {
	http.HandleFunc("/price", makeHTTPHandlerFunc(s.handleFetchPrice))
	http.HandleFunc("/market", makeHTTPHandlerFunc(s.handleFetchMarketStats))
	http.HandleFunc("/info", makeHTTPHandlerFunc(s.handleFetchCoinInfo))
	http.HandleFunc("/tickers", makeHTTPHandlerFunc(s.handleFetchSupportedTickers))
	http.HandleFunc("/batch", makeHTTPHandlerFunc(s.handleBatchFetchPrice))
	http.ListenAndServe(s.listenAddr, nil)
}

func makeHTTPHandlerFunc(apiFn APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "requestID", rand.Intn(100_000_000))
		if err := apiFn(ctx, w, r); err != nil {
			writeJSON(w, http.StatusBadRequest, &types.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			})
		}
	}
}

// GET /price?ticker=BTC
func (s *JSONAPIServer) handleFetchPrice(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ticker := r.URL.Query().Get("ticker")
	price, err := s.svc.FetchPrice(ctx, ticker)
	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, &types.PriceResponse{
		Ticker:   ticker,
		Price:    price,
		Currency: "USD",
	})
}

// GET /market?ticker=BTC
func (s *JSONAPIServer) handleFetchMarketStats(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ticker := r.URL.Query().Get("ticker")
	stats, err := s.svc.FetchMarketStats(ctx, ticker)
	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, stats)
}

// GET /info?ticker=BTC
func (s *JSONAPIServer) handleFetchCoinInfo(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ticker := r.URL.Query().Get("ticker")
	info, err := s.svc.FetchCoinInfo(ctx, ticker)
	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, info)
}

// GET /tickers
func (s *JSONAPIServer) handleFetchSupportedTickers(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	resp, err := s.svc.FetchSupportedTickers(ctx)
	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, resp)
}

// POST /batch   body: {"tickers":["BTC","ETH"]}
func (s *JSONAPIServer) handleBatchFetchPrice(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return writeJSON(w, http.StatusMethodNotAllowed, &types.ErrorResponse{
			Code:    http.StatusMethodNotAllowed,
			Message: "method not allowed — use POST",
		})
	}
	var req types.BatchPriceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}
	resp, err := s.svc.BatchFetchPrice(ctx, req.Tickers)
	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, resp)
}

func writeJSON(w http.ResponseWriter, s int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(s)
	return json.NewEncoder(w).Encode(v)
}
