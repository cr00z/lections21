package main

import (
	"context"
	"hw-async/domain"
	"hw-async/generator"
	"math"
	"os"
	"os/signal"
	"syscall"
	"time"
	"sync"

	log "github.com/sirupsen/logrus"
)

var tickers = []string{"AAPL", "SBER", "NVDA", "TSLA"}

func proceedCandle1m(
	wg *sync.WaitGroup,
	logger *log.Logger,
	in <-chan domain.Price) <-chan domain.Candle {
	
	out := make(chan domain.Candle)
	go func() {
		defer wg.Done()
		candles := make(map[string]domain.Candle)

		for price := range in {
			logger.Infof("in price: %+v", price)
			period, _ := domain.PeriodTS(domain.CandlePeriod1m, price.TS)
			candle, inMap := candles[price.Ticker]

			if candle.TS != period {
				if inMap {
					logger.Infof("out 1m candle: %+v", candle)
					out <- candle
				}

				candles[price.Ticker] = domain.Candle{
					Ticker: price.Ticker,
					Period: domain.CandlePeriod1m,
					Open:   price.Value,
					High:   price.Value,
					Low:    price.Value,
					Close:  price.Value,
					TS:     period,
				}
			} else {
				if price.Value > candle.High {
					candle.High = price.Value
				} else if price.Value < candle.Low {
					candle.Low = price.Value
				}
				candle.Close = price.Value
				candles[price.Ticker] = candle
			}
		}
		close(out)
		logger.Info("candle 1m gorutine done")
	}()
	return out
}

func convertPriceToCandle(wg *sync.WaitGroup,
	logger *log.Logger,
	in <-chan domain.Price) <-chan domain.Candle {

	out := make(chan domain.Candle)
	go func() {
		defer wg.Done()

		for price := range in {
			out <- domain.Candle{
				Ticker: price.Ticker,
				Open:   price.Value,
				High:   price.Value,
				Low:    price.Value,
				Close:  price.Value,
				TS:     price.TS,
			}
		}

		close(out)
		logger.Info("convert price to candle gorutine done")
	}()
	return out
}

func proceedCandle(
	wg *sync.WaitGroup,
	logger *log.Logger,
	candlePeriod domain.CandlePeriod,
	in <-chan domain.Candle) <-chan domain.Candle {

	out := make(chan domain.Candle)

	go func() {
		defer wg.Done()
		candles := make(map[string]domain.Candle)

		for currentCandle := range in {
			currentPeriod, _ := domain.PeriodTS(candlePeriod, currentCandle.TS)
			storedCandle, inMap := candles[currentCandle.Ticker]

			if storedCandle.TS != currentPeriod {
				if inMap {
					logger.Infof("out %v candle: %+v", candlePeriod, storedCandle)
					out <- storedCandle
				}
				currentCandle.Period = candlePeriod
				currentCandle.TS = currentPeriod
				candles[currentCandle.Ticker] = currentCandle
			} else {
				storedCandle.High = math.Max(storedCandle.High, currentCandle.High)
				storedCandle.Low = math.Min(storedCandle.Low, currentCandle.Low)
				storedCandle.Close = currentCandle.Close
				candles[currentCandle.Ticker] = storedCandle
			}
		}

		close(out)
		logger.Infof("candle %v gorutine done", candlePeriod)
	}()
	return out
}

func main() {
	logger := log.New()
	ctx, cancel := context.WithCancel(context.Background())

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	var wg sync.WaitGroup

	pg := generator.NewPricesGenerator(generator.Config{
		Factor:  10,
		Delay:   time.Millisecond * 500,
		Tickers: tickers,
	})

	logger.Info("start prices generator...")
	prices := pg.Prices(ctx)

	wg.Add(4)
	candles := convertPriceToCandle(&wg, logger, prices)
	candle1m := proceedCandle(&wg, logger, domain.CandlePeriod1m, candles)
	candle2m := proceedCandle(&wg, logger, domain.CandlePeriod2m, candle1m)
	candle10m := proceedCandle(&wg, logger, domain.CandlePeriod10m, candle2m)
	for {
		select {
		case <-candle10m:
		case <-termChan:
			logger.Info("shutdown signal received")
			cancel()
			wg.Wait()
			logger.Info("all prices processed")
			return
		}
	}
}
