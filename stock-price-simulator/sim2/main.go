package main

import (
	"cmp"
	"fmt"
	"math/rand/v2"
	"slices"
	"sync"
	"time"
)

type Stock struct {
	name             string
	symbol           string
	price            float32
	historicalPrices []float32
	mu               sync.Mutex
}

type Market struct {
	stocks []*Stock
	mu     sync.Mutex
}

func (m *Market) changePrice(stock *Stock) {
	stock.mu.Lock()
	defer stock.mu.Unlock()
	stock.price += (rand.Float32() - 0.5) * 10
	stock.historicalPrices = append(stock.historicalPrices, stock.price)
}

func (m *Market) DisplayPrices() {
	m.mu.Lock()
	defer m.mu.Unlock()
	fmt.Println("Symbol | Name | Current Price | Historical Prices")
	fmt.Println("------------------------------------------------")

	compareSymbols := func(a, b *Stock) int {
		return cmp.Compare(a.symbol, b.symbol)
	}

	slices.SortFunc(m.stocks, compareSymbols)
	for _, stock := range m.stocks {
		stock.mu.Lock()
		fmt.Printf("%-6s | %-20s | %-13.2f | %v\n",
			stock.symbol,
			stock.name,
			stock.price,
			stock.historicalPrices)
		stock.mu.Unlock()
	}
	fmt.Print("\n\n")
}

func (m *Market) UpdateAllPrices() {
	var wg sync.WaitGroup
	for _, stock := range m.stocks {
		wg.Add(1)
		go func(s *Stock) {
			defer wg.Done()
			m.changePrice(s)
		}(stock)
	}
	wg.Wait()
}

func main() {
	m := Market{
		stocks: []*Stock{
			{name: "Apple", symbol: "AAPL", price: 150.25, historicalPrices: []float32{148.56, 150.25, 149.80, 151.60, 147.92}},
			{name: "Microsoft", symbol: "MSFT", price: 290.17, historicalPrices: []float32{288.30, 289.67, 290.17, 291.32, 287.72}},
			{name: "Amazon", symbol: "AMZN", price: 102.11, historicalPrices: []float32{100.79, 101.56, 102.11, 103.29, 99.92}},
			{name: "Alphabet", symbol: "GOOGL", price: 105.22, historicalPrices: []float32{104.61, 105.97, 105.22, 106.06, 103.73}},
			{name: "Meta", symbol: "META", price: 195.32, historicalPrices: []float32{193.62, 194.02, 195.32, 196.64, 192.53}},
			{name: "Tesla", symbol: "TSLA", price: 180.14, historicalPrices: []float32{178.97, 181.67, 180.14, 183.25, 177.49}},
			{name: "NVIDIA", symbol: "NVDA", price: 277.77, historicalPrices: []float32{275.79, 278.01, 277.77, 280.61, 274.86}},
			{name: "JPMorgan Chase", symbol: "JPM", price: 138.73, historicalPrices: []float32{137.22, 138.05, 138.73, 139.83, 136.67}},
			{name: "Johnson & Johnson", symbol: "JNJ", price: 158.14, historicalPrices: []float32{157.56, 158.32, 158.14, 159.07, 156.89}},
			{name: "Visa", symbol: "V", price: 232.48, historicalPrices: []float32{230.93, 231.89, 232.48, 233.76, 229.85}},
		},
	}

	ticker := time.NewTicker(200 * time.Millisecond)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				m.UpdateAllPrices()
				m.DisplayPrices()
			}
		}
	}()

	time.Sleep(1 * time.Second)
	ticker.Stop()
	done <- true

}
