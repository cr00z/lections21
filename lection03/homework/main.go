package main

import (
	"context"
	"fmt"
	"hw-async/domain"
	"hw-async/generator"
	"math"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

var tickers = []string{"AAPL", "SBER", "NVDA", "TSLA"}

func convertPriceToCandle(wg *sync.WaitGroup, logger *log.Logger,
	in <-chan domain.Price) <-chan domain.Candle {
	out := make(chan domain.Candle)

	go func() {
		defer wg.Done()

		for price := range in {
			logger.Infof("in price: %+v", price)

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

func proceedCandle(wg *sync.WaitGroup, logger *log.Logger, file *os.File,
	currentPeriod domain.CandlePeriod, in <-chan domain.Candle) <-chan domain.Candle {
	out := make(chan domain.Candle)

	go func() {
		defer wg.Done()

		candles := make(map[string]domain.Candle)

		for currentCandle := range in {
			currentTS, _ := domain.PeriodTS(currentPeriod, currentCandle.TS)
			storedCandle, inMap := candles[currentCandle.Ticker]

			if storedCandle.TS != currentTS {
				if inMap {
					logger.Infof("out %v candle: %+v", currentPeriod, storedCandle)

					if file != nil {
						fmt.Fprintf(file, "%s,%s,%f,%f,%f,%f\n",
							storedCandle.Ticker,
							storedCandle.TS,
							storedCandle.Open,
							storedCandle.High,
							storedCandle.Low,
							storedCandle.Close)
					}

					out <- storedCandle
				}
				currentCandle.Period = currentPeriod
				currentCandle.TS = currentTS
				candles[currentCandle.Ticker] = currentCandle
			} else {
				storedCandle.High = math.Max(storedCandle.High, currentCandle.High)
				storedCandle.Low = math.Min(storedCandle.Low, currentCandle.Low)
				storedCandle.Close = currentCandle.Close
				candles[currentCandle.Ticker] = storedCandle
			}
		}

		close(out)
		logger.Infof("candle %v gorutine done", currentPeriod)
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
	for _, candlePeriod := range []domain.CandlePeriod{
		domain.CandlePeriod1m,
		domain.CandlePeriod2m,
		domain.CandlePeriod10m,
	} {
		file, err := os.Create("candles_" + string(candlePeriod) + ".csv")
		if err != nil {
			logger.Error(err)
			termChan <- syscall.SIGINT
		} else {
			defer file.Close()
		}
		candles = proceedCandle(&wg, logger, file, candlePeriod, candles)
	}

	for {
		select {
		case <-candles:
			// get 10m candle, do nothing
		case <-termChan:
			logger.Info("shutdown signal received")
			cancel()
			wg.Wait()
			logger.Info("all prices processed")
			return
		}
	}
}
