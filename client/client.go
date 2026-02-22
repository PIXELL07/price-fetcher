package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/PIXELL07/price-fetcher/types"
)

type Client struct {
	endpoint   string
	httpClient *http.Client
}

func New(endpoint string) *Client {
	return &Client{
		endpoint:   endpoint,
		httpClient: http.DefaultClient,
	}
}

// GET /price?ticker=<ticker>
func (c *Client) FetchPrice(ctx context.Context, ticker string) (*types.PriceResponse, error) {
	url := fmt.Sprintf("%s/price?ticker=%s", c.endpoint, ticker)
	priceResp := new(types.PriceResponse)
	if err := c.get(ctx, url, priceResp); err != nil {
		return nil, err
	}
	return priceResp, nil
}

// GET /market?ticker=<ticker>
func (c *Client) FetchMarketStats(ctx context.Context, ticker string) (*types.MarketStats, error) {
	url := fmt.Sprintf("%s/market?ticker=%s", c.endpoint, ticker)
	stats := new(types.MarketStats)
	if err := c.get(ctx, url, stats); err != nil {
		return nil, err
	}
	return stats, nil
}

// GET /info?ticker=<ticker>
func (c *Client) FetchCoinInfo(ctx context.Context, ticker string) (*types.CoinInfo, error) {
	url := fmt.Sprintf("%s/info?ticker=%s", c.endpoint, ticker)
	info := new(types.CoinInfo)
	if err := c.get(ctx, url, info); err != nil {
		return nil, err
	}
	return info, nil
}

// GET /tickers
func (c *Client) FetchSupportedTickers(ctx context.Context) (*types.SupportedTickersResponse, error) {
	url := fmt.Sprintf("%s/tickers", c.endpoint)
	resp := new(types.SupportedTickersResponse)
	if err := c.get(ctx, url, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// POST /batch   body: {"tickers":[...]}
func (c *Client) BatchFetchPrice(ctx context.Context, tickers []string) (*types.BatchPriceResponse, error) {
	url := fmt.Sprintf("%s/batch", c.endpoint)
	body, err := json.Marshal(types.BatchPriceRequest{Tickers: tickers})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	httpResp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("service responded with non-OK status: %d", httpResp.StatusCode)
	}

	batchResp := new(types.BatchPriceResponse)
	if err := json.NewDecoder(httpResp.Body).Decode(batchResp); err != nil {
		return nil, err
	}
	return batchResp, nil
}

func (c *Client) get(ctx context.Context, url string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("service responded with non-OK status: %d", resp.StatusCode)
	}
	return json.NewDecoder(resp.Body).Decode(out)
}
