package main

import (
    "encoding/csv"
    "fmt"
    "log"
    "os"
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

func main() {

	start := time.Now()

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

    for _, rec := range rows {
        r := Records{}
        r.saveRecords(rec)
        allRecords = append(allRecords, r)
    }

	elapsed := time.Since(start)
	
    fmt.Printf("Total records: %d\n", len(allRecords))
    fmt.Println("First record:", allRecords[0])
    fmt.Println("Last record:", allRecords[len(allRecords)-1])
	fmt.Printf("Record processing time: %v\n", elapsed)
}