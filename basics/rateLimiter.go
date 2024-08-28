package main

import (
	"fmt"
	"time"
)

func main(){
	requests := make(chan int, 5)
	for i:=1; i <=5; i++{
		requests <- i
	}

	close(requests)

	limiter := time.NewTicker(200 * time.Millisecond)

	for req:= range requests{
		<-limiter.C
		fmt.Println("Request: ", req, time.Now())
	}


}