package main

import (
	"fmt"
	"time"
)

func f(from string){
	for i:= range 3{
		fmt.Printf("%s : %d\n", from, i)
	}
}

func main(){
	f("direct")

	go f("goroutine")

	go func(msg string){
		fmt.Println(msg)
	}("going")
	
	// use can subsitute the wait with a channel - channel sync
	time.Sleep(time.Second)

	fmt.Println("\ndone")

}