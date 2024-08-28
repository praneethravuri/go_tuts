package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type Records struct {
	region       string
	country      string
	itemType     string
	salesChannel string
}

func (r *Records) saveRecords(rec []string) {
	r.region = rec[0]
	r.country = rec[1]
	r.itemType = rec[2]
	r.salesChannel = rec[3]
}

func processData(start, end int, rows [][]string, results chan<- Records, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := start; i < end; i++ {
		r := Records{}
		r.saveRecords(rows[i])
		results <- r
	}
}

func main() {
	start := time.Now()

	var wg sync.WaitGroup

	file, err := os.Open("../10000-records.csv")
	if err != nil {
		log.Fatal("Error while reading file: ", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		log.Fatal("Error while reading rows: ", err)
	}

	var allRecords []Records

	numOfRows := len(rows)

	results := make(chan Records, numOfRows)

	wg.Add(2)
	go processData(0, numOfRows/2, rows, results, &wg)
	go processData(numOfRows/2, numOfRows, rows, results, &wg)

	wg.Wait()
	close(results)

	for r := range results {
		allRecords = append(allRecords, r)
	}

	elapsed := time.Since(start)

	fmt.Printf("Total records: %d\n", len(allRecords))
	fmt.Println("First record:", allRecords[0])
	fmt.Println("Last record:", allRecords[len(allRecords)-1])
	fmt.Printf("Record processing time: %v\n", elapsed)

}
