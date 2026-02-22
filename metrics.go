package main

import (
	"context"

	"github.com/PIXELL07/price-fetcher/types"
	"github.com/sirupsen/logrus"
)

type metricService struct {
	next PriceFetcher
}

func NewMetricService(next PriceFetcher) PriceFetcher {
	return &metricService{next: next}
}

func (s *metricService) FetchPrice(ctx context.Context, ticker string) (price float64, err error) {
	defer func() {
		logrus.WithFields(logrus.Fields{"method": "FetchPrice", "ticker": ticker, "price": price, "err": err}).Debug("metrics")
	}()
	return s.next.FetchPrice(ctx, ticker)
}

func (s *metricService) FetchMarketStats(ctx context.Context, ticker string) (stats *types.MarketStats, err error) {
	defer func() {
		logrus.WithFields(logrus.Fields{"method": "FetchMarketStats", "ticker": ticker, "err": err}).Debug("metrics")
	}()
	return s.next.FetchMarketStats(ctx, ticker)
}

func (s *metricService) FetchCoinInfo(ctx context.Context, ticker string) (info *types.CoinInfo, err error) {
	defer func() {
		logrus.WithFields(logrus.Fields{"method": "FetchCoinInfo", "ticker": ticker, "err": err}).Debug("metrics")
	}()
	return s.next.FetchCoinInfo(ctx, ticker)
}

func (s *metricService) FetchSupportedTickers(ctx context.Context) (resp *types.SupportedTickersResponse, err error) {
	defer func() {
		logrus.WithFields(logrus.Fields{"method": "FetchSupportedTickers", "err": err}).Debug("metrics")
	}()
	return s.next.FetchSupportedTickers(ctx)
}

func (s *metricService) BatchFetchPrice(ctx context.Context, tickers []string) (resp *types.BatchPriceResponse, err error) {
	defer func() {
		logrus.WithFields(logrus.Fields{"method": "BatchFetchPrice", "tickers": tickers, "err": err}).Debug("metrics")
	}()
	return s.next.BatchFetchPrice(ctx, tickers)
}
