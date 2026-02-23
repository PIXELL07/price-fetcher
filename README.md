# price-fetcher

A lightweight crypto price fetcher microservice written in Go. It exposes a JSON REST API for querying mock cryptocurrency prices, market stats, and coin metadata across 40+ supported tickers.

## Features

- Single price lookup by ticker
- Market stats (market cap, 24h volume, high/low, circulating supply)
- Coin info (name, category, description, tags, website)
- List all supported tickers
- Batch price lookup for multiple tickers in one request
- Structured request logging via [logrus](https://github.com/sirupsen/logrus)
- Metrics layer using the middleware/decorator pattern
- Type-safe Go client package for programmatic usage
- Docker support with a minimal `scratch`-based image

## Project Structure

```
price-fetcher/
├── main.go           # Entry point, wires up service layers and HTTP server
├── api.go            # HTTP handlers and JSON API server
├── service.go        # Core PriceFetcher interface and mock implementation
├── login.go          # Logging middleware (decorator pattern)
├── metrics.go        # Metrics middleware (decorator pattern)
├── types/
│   └── types.go      # Shared request/response types
├── client/
│   └── client.go     # Go HTTP client for consuming the API
├── proto/
│   └── service.pro   # Protobuf service definition
├── DockerFile
└── MakeFile
```

## Requirements

- Go 1.22+
- Docker (optional, for containerised runs)
- `golangci-lint` (optional, for linting)
- `protoc` + `protoc-gen-go` (optional, for proto generation)

## Getting Started

### Run locally

```bash
# Build and start the server on :3000
make run

# Or run without rebuilding (faster during development)
make dev

# Custom listen address
go run . -listenaddr=:8080
```

### Run with Docker

```bash
make docker-build
make docker-run
```

The container exposes port `3000` and runs as a fully static binary on a `scratch` base image.

## API Endpoints

All responses are JSON. Errors return `{"code": <int>, "message": "<string>"}`.

### `GET /price?ticker=BTC`

Returns the current price for a single ticker.

```json
{
  "ticker": "BTC",
  "price": 20000,
  "currency": "USD",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

### `GET /market?ticker=ETH`

Returns extended market statistics for a ticker.

```json
{
  "ticker": "ETH",
  "price": 2000,
  "market_cap": 2000000000,
  "volume_24h": 100000000,
  "change_24h_pct": -1.5,
  "high_24h": 2060,
  "low_24h": 1940,
  "circulating_supply": 1000000,
  "timestamp": "2024-01-01T00:00:00Z"
}
```

### `GET /info?ticker=SOL`

Returns descriptive information about a coin.

```json
{
  "ticker": "SOL",
  "name": "Solana",
  "category": "Layer 1",
  "description": "High-throughput blockchain using Proof of History.",
  "tags": ["pos", "high-throughput"],
  "website": "https://solana.com"
}
```

### `GET /tickers`

Returns all supported tickers sorted alphabetically.

```json
{
  "tickers": ["AAVE", "ADA", "ALGO", "..."],
  "count": 40
}
```

### `POST /batch`

Fetches prices for multiple tickers in a single request.

**Request body:**
```json
{ "tickers": ["BTC", "ETH", "SOL"] }
```

**Response:**
```json
{
  "prices": [
    { "ticker": "BTC", "price": 20000, "currency": "USD", "timestamp": "..." },
    { "ticker": "ETH", "price": 2000,  "currency": "USD", "timestamp": "..." },
    { "ticker": "SOL", "price": 150,   "currency": "USD", "timestamp": "..." }
  ],
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## Supported Tickers

| Category       | Tickers                                                      |
|----------------|--------------------------------------------------------------|
| Layer 1        | BTC, ETH, SOL, BNB, ADA, AVAX, DOT, MATIC, TRX, TON, NEAR, ICP, FTM, ALGO |
| DeFi           | ARB, OP, LINK, UNI, AAVE, CRV, MKR, SNX, LDO, GRT            |
| Stablecoins    | USDT, USDC, DAI, FRAX, TUSD                                  |
| Exchange Tokens| OKB, CRO, KCS                                                |
| Other / Meme   | XRP, LTC, DOGE, SHIB, XLM, ATOM, ETC, XMR, PEPE, FLOKI       |

## Go Client

A typed client package is included for integrating the service into other Go applications.

```go
import "github.com/PIXELL07/price-fetcher/client"

c := client.New("http://localhost:3000")

// Fetch a single price
price, err := c.FetchPrice(ctx, "BTC")

// Fetch market stats
stats, err := c.FetchMarketStats(ctx, "ETH")

// Fetch coin info
info, err := c.FetchCoinInfo(ctx, "SOL")

// List all tickers
tickers, err := c.FetchSupportedTickers(ctx)

// Batch fetch
batch, err := c.BatchFetchPrice(ctx, []string{"BTC", "ETH", "SOL"})
```

## Makefile Reference

| Command            | Description                              |
|--------------------|------------------------------------------|
| `make build`       | Compile the binary to `bin/price-fetcher`|
| `make run`         | Build and start the server               |
| `make dev`         | Run with `go run` (no build step)        |
| `make test`        | Run all tests with race detector         |
| `make lint`        | Lint with `golangci-lint`                |
| `make proto`       | Generate Go code from proto definition   |
| `make docker-build`| Build the Docker image                   |
| `make docker-run`  | Run the Docker container on port 3000    |
| `make clean`       | Remove build artifacts                   |

## Architecture

The service is structured using the **decorator (middleware) pattern**. The core `priceFetcher` struct implements the `PriceFetcher` interface. It is then wrapped by:

1. `MetricService` — logs debug-level metrics for every method call
2. `LoggingService` — logs structured info-level entries including request ID, method, ticker, duration, and any error

This keeps each concern isolated and makes it straightforward to add further middleware (e.g. caching, rate limiting) without touching business logic.
