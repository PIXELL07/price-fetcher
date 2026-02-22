syntax = "proto3";

option go_package = "github.com/PIXELL07/price-fetcher/proto";

import "google/protobuf/timestamp.proto";


service PriceFetcher {
    // Fetch the current price for a single ticker.
    rpc FetchPrice(PriceRequest) returns (PriceResponse);

    // Fetch extended market stats for a single ticker.
    rpc FetchMarketStats(PriceRequest) returns (MarketStatsResponse);

    // Fetch static metadata / coin info for a single ticker.
    rpc FetchCoinInfo(PriceRequest) returns (CoinInfoResponse);

    // List all supported tickers.
    rpc FetchSupportedTickers(Empty) returns (SupportedTickersResponse);

    // Fetch prices for multiple tickers in one call.
    rpc BatchFetchPrice(BatchPriceRequest) returns (BatchPriceResponse);
}


// Empty is used for RPCs that require no input.
message Empty {}


message PriceRequest {
    string ticker = 1;
}

message BatchPriceRequest {
    repeated string tickers = 1;
}


message PriceResponse {
    string                     ticker    = 1;
    float                      price     = 2;
    string                     currency  = 3; // e.g. "USD"
    google.protobuf.Timestamp  timestamp = 4;
}

message MarketStatsResponse {
    string                     ticker             = 1;
    float                      price              = 2;
    float                      market_cap         = 3;
    float                      volume_24h         = 4;
    float                      change_24h_pct     = 5;
    float                      high_24h           = 6;
    float                      low_24h            = 7;
    float                      circulating_supply = 8;
    google.protobuf.Timestamp  timestamp          = 9;
}

message CoinInfoResponse {
    string          ticker      = 1;
    string          name        = 2;
    string          category    = 3; // e.g. "Layer 1", "DeFi", "Stablecoin"
    string          description = 4;
    repeated string tags        = 5;
    string          website     = 6;
}

message SupportedTickersResponse {
    repeated string tickers = 1;
    int32           count   = 2;
}

message BatchPriceResponse {
    repeated PriceResponse prices    = 1;
    google.protobuf.Timestamp timestamp = 2;
}


message ErrorResponse {
    int32  code    = 1;
    string message = 2;
}
