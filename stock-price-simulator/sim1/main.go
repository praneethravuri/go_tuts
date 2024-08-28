package main

import (
	"fmt"
	"math/rand/v2"
	"strings"
	"sync"
	"time"
)

type Stock struct {
	symbol           string
	currentPrice     float32
	historicalPrices []float32
	mu               sync.Mutex
}

func priceGenerator(stock *Stock, wg *sync.WaitGroup) {
	defer wg.Done()

	stock.mu.Lock()
	defer stock.mu.Unlock()

	changeValue := rand.IntN(2)
	delta := rand.Float32() * 10
	if changeValue == 0 {
		stock.currentPrice -= delta
	} else {
		stock.currentPrice += delta
	}
	stock.historicalPrices = append(stock.historicalPrices, stock.currentPrice)

	fmt.Printf("| %-5s | %8.2f | [%s]\n",
		stock.symbol,
		stock.currentPrice,
		formatHistoricalPrices(stock.historicalPrices))
}

func formatHistoricalPrices(prices []float32) string {
	strPrices := make([]string, len(prices))
	for i, price := range prices {
		strPrices[i] = fmt.Sprintf("%.2f", price)
	}
	return strings.Join(strPrices, ", ")
}

func main() {
	stockSymbols := []string{"AAPL", "NVDA", "MSFT", "GOOG", "AMZN"}
	stocks := make([]*Stock, 0, len(stockSymbols))

	for _, symbol := range stockSymbols {
		stock := &Stock{
			symbol:           symbol,
			currentPrice:     rand.Float32() * 100,
			historicalPrices: make([]float32, 0),
		}
		stocks = append(stocks, stock)
	}

	ticker := time.NewTicker(1 * time.Second)
	done := make(chan bool)
	var wg sync.WaitGroup

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				wg.Add(len(stocks))
				for _, stock := range stocks {
					go priceGenerator(stock, &wg)
				}
				fmt.Println()
			}
		}
	}()

	time.Sleep(10 * time.Second)

	ticker.Stop()
	done <- true
	wg.Wait()

	fmt.Println("Simulation completed")
}
