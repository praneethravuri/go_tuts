package main

import (
	"fmt"
	"time"
)

func main(){
	ticker := time.NewTicker(300 * time.Millisecond)
	done := make(chan bool)
	numberOfTics := 0

	go func(){
		for{
			select{
			case <-done:
				return
			case t := <-ticker.C:
				numberOfTics += 1
				fmt.Println("Tick at: ", t)
			}
		}
	}()

	time.Sleep(2000 * time.Millisecond)
	ticker.Stop()
	done <-true
	fmt.Println("Ticker stopped")
	fmt.Println(numberOfTics)
}