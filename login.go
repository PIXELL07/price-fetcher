package main

import (
	"context"
	"time"

	"github.com/PIXELL07/price-fetcher/types"
	"github.com/sirupsen/logrus"
)

type loggingService struct {
	next PriceFetcher
}

func NewLoggingService(next PriceFetcher) PriceFetcher {
	return &loggingService{next: next}
}

func (s *loggingService) FetchPrice(ctx context.Context, ticker string) (price float64, err error) {
	defer func(begin time.Time) {
		logrus.WithFields(logrus.Fields{
			"requestID": ctx.Value("requestID"),
			"method":    "FetchPrice",
			"ticker":    ticker,
			"took":      time.Since(begin),
			"err":       err,
			"price":     price,
		}).Info("fetchPrice")
	}(time.Now())
	return s.next.FetchPrice(ctx, ticker)
}

func (s *loggingService) FetchMarketStats(ctx context.Context, ticker string) (stats *types.MarketStats, err error) {
	defer func(begin time.Time) {
		logrus.WithFields(logrus.Fields{
			"requestID": ctx.Value("requestID"),
			"method":    "FetchMarketStats",
			"ticker":    ticker,
			"took":      time.Since(begin),
			"err":       err,
		}).Info("fetchMarketStats")
	}(time.Now())
	return s.next.FetchMarketStats(ctx, ticker)
}

func (s *loggingService) FetchCoinInfo(ctx context.Context, ticker string) (info *types.CoinInfo, err error) {
	defer func(begin time.Time) {
		logrus.WithFields(logrus.Fields{
			"requestID": ctx.Value("requestID"),
			"method":    "FetchCoinInfo",
			"ticker":    ticker,
			"took":      time.Since(begin),
			"err":       err,
		}).Info("fetchCoinInfo")
	}(time.Now())
	return s.next.FetchCoinInfo(ctx, ticker)
}

func (s *loggingService) FetchSupportedTickers(ctx context.Context) (resp *types.SupportedTickersResponse, err error) {
	defer func(begin time.Time) {
		logrus.WithFields(logrus.Fields{
			"requestID": ctx.Value("requestID"),
			"method":    "FetchSupportedTickers",
			"took":      time.Since(begin),
			"err":       err,
		}).Info("fetchSupportedTickers")
	}(time.Now())
	return s.next.FetchSupportedTickers(ctx)
}

func (s *loggingService) BatchFetchPrice(ctx context.Context, tickers []string) (resp *types.BatchPriceResponse, err error) {
	defer func(begin time.Time) {
		logrus.WithFields(logrus.Fields{
			"requestID": ctx.Value("requestID"),
			"method":    "BatchFetchPrice",
			"tickers":   tickers,
			"took":      time.Since(begin),
			"err":       err,
		}).Info("batchFetchPrice")
	}(time.Now())
	return s.next.BatchFetchPrice(ctx, tickers)
}
