package main

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/PIXELL07/price-fetcher/types"
)

type PriceFetcher interface {
	FetchPrice(context.Context, string) (float64, error)
	FetchMarketStats(context.Context, string) (*types.MarketStats, error)
	FetchCoinInfo(context.Context, string) (*types.CoinInfo, error)
	FetchSupportedTickers(context.Context) (*types.SupportedTickersResponse, error)
	BatchFetchPrice(context.Context, []string) (*types.BatchPriceResponse, error)
}

type priceFetcher struct{}

func (s *priceFetcher) FetchPrice(ctx context.Context, ticker string) (float64, error) {
	time.Sleep(100 * time.Millisecond)
	price, ok := priceMocks[ticker]
	if !ok {
		return 0, fmt.Errorf("ticker (%s) is not supported", ticker)
	}
	return price, nil
}

func (s *priceFetcher) FetchMarketStats(ctx context.Context, ticker string) (*types.MarketStats, error) {
	time.Sleep(100 * time.Millisecond)
	price, ok := priceMocks[ticker]
	if !ok {
		return nil, fmt.Errorf("ticker (%s) is not supported", ticker)
	}
	// Mock derived stats — in production these would come from a real data source.
	return &types.MarketStats{
		Ticker:     ticker,
		Price:      price,
		MarketCap:  price * 1_000_000,
		Volume24h:  price * 50_000,
		Change24h:  -1.5,
		High24h:    price * 1.03,
		Low24h:     price * 0.97,
		CircSupply: 1_000_000,
		Timestamp:  time.Now().UTC(),
	}, nil
}

func (s *priceFetcher) FetchCoinInfo(ctx context.Context, ticker string) (*types.CoinInfo, error) {
	time.Sleep(50 * time.Millisecond)
	info, ok := coinInfoMocks[ticker]
	if !ok {
		return nil, fmt.Errorf("ticker (%s) is not supported", ticker)
	}
	return &info, nil
}

func (s *priceFetcher) FetchSupportedTickers(ctx context.Context) (*types.SupportedTickersResponse, error) {
	tickers := make([]string, 0, len(priceMocks))
	for t := range priceMocks {
		tickers = append(tickers, t)
	}
	sort.Strings(tickers)
	return &types.SupportedTickersResponse{
		Tickers: tickers,
		Count:   len(tickers),
	}, nil
}

func (s *priceFetcher) BatchFetchPrice(ctx context.Context, tickers []string) (*types.BatchPriceResponse, error) {
	time.Sleep(100 * time.Millisecond)
	now := time.Now().UTC()
	prices := make([]types.PriceResponse, 0, len(tickers))
	for _, ticker := range tickers {
		price, ok := priceMocks[ticker]
		if !ok {
			return nil, fmt.Errorf("ticker (%s) is not supported", ticker)
		}
		prices = append(prices, types.PriceResponse{
			Ticker:    ticker,
			Price:     price,
			Currency:  "USD",
			Timestamp: now,
		})
	}
	return &types.BatchPriceResponse{Prices: prices, Timestamp: now}, nil
}

var priceMocks = map[string]float64{
	// Layer 1
	"BTC":   20_000,
	"ETH":   2_000,
	"SOL":   150,
	"BNB":   300,
	"ADA":   0.45,
	"AVAX":  35,
	"DOT":   7.5,
	"MATIC": 0.85,
	"TRX":   0.08,
	"TON":   2.10,
	"NEAR":  3.50,
	"ICP":   5.20,
	"FTM":   0.40,
	"ALGO":  0.17,
	//  DeFi
	"ARB":  1.20,
	"OP":   1.80,
	"LINK": 14.50,
	"UNI":  6.30,
	"AAVE": 95,
	"CRV":  0.55,
	"MKR":  1_200,
	"SNX":  2.80,
	"LDO":  2.10,
	"GRT":  0.14,
	// Stablecoins
	"USDT": 1.00,
	"USDC": 1.00,
	"DAI":  1.00,
	"FRAX": 1.00,
	"TUSD": 1.00,
	// Exchange tokens
	"OKB": 45,
	"CRO": 0.09,
	"KCS": 7.50,
	// Meme
	"XRP":   0.55,
	"LTC":   85,
	"DOGE":  0.08,
	"SHIB":  0.000009,
	"XLM":   0.12,
	"ATOM":  10.50,
	"ETC":   18,
	"XMR":   155,
	"PEPE":  0.0000015,
	"FLOKI": 0.00003,
}

var coinInfoMocks = map[string]types.CoinInfo{
	"BTC":  {Ticker: "BTC", Name: "Bitcoin", Category: "Layer 1", Description: "The original decentralised cryptocurrency.", Tags: []string{"pow", "store-of-value"}, Website: "https://bitcoin.org"},
	"ETH":  {Ticker: "ETH", Name: "Ethereum", Category: "Layer 1", Description: "Programmable blockchain with smart contracts.", Tags: []string{"smart-contracts", "pos"}, Website: "https://ethereum.org"},
	"SOL":  {Ticker: "SOL", Name: "Solana", Category: "Layer 1", Description: "High-throughput blockchain using Proof of History.", Tags: []string{"pos", "high-throughput"}, Website: "https://solana.com"},
	"BNB":  {Ticker: "BNB", Name: "BNB", Category: "Layer 1", Description: "Native token of the BNB Chain ecosystem.", Tags: []string{"exchange-token", "bsc"}, Website: "https://bnbchain.org"},
	"USDT": {Ticker: "USDT", Name: "Tether", Category: "Stablecoin", Description: "USD-pegged stablecoin issued by Tether.", Tags: []string{"stablecoin", "usd-pegged"}, Website: "https://tether.to"},
	"USDC": {Ticker: "USDC", Name: "USD Coin", Category: "Stablecoin", Description: "Regulated USD-pegged stablecoin by Circle.", Tags: []string{"stablecoin", "regulated"}, Website: "https://circle.com/usdc"},
	"LINK": {Ticker: "LINK", Name: "Chainlink", Category: "DeFi", Description: "Decentralised oracle network for smart contracts.", Tags: []string{"oracle", "defi"}, Website: "https://chain.link"},
	"UNI":  {Ticker: "UNI", Name: "Uniswap", Category: "DeFi", Description: "Governance token of the Uniswap DEX protocol.", Tags: []string{"dex", "amm", "defi"}, Website: "https://uniswap.org"},
	"DOGE": {Ticker: "DOGE", Name: "Dogecoin", Category: "Meme", Description: "The original meme coin, based on the Shiba Inu meme.", Tags: []string{"meme", "pow"}, Website: "https://dogecoin.com"},
	"ARB":  {Ticker: "ARB", Name: "Arbitrum", Category: "Layer 2", Description: "Optimistic rollup scaling solution for Ethereum.", Tags: []string{"l2", "rollup"}, Website: "https://arbitrum.io"},
	"OP":   {Ticker: "OP", Name: "Optimism", Category: "Layer 2", Description: "Optimistic rollup L2 built on Ethereum.", Tags: []string{"l2", "rollup"}, Website: "https://optimism.io"},
	"AAVE": {Ticker: "AAVE", Name: "Aave", Category: "DeFi", Description: "Decentralised lending and borrowing protocol.", Tags: []string{"lending", "defi"}, Website: "https://aave.com"},
}
