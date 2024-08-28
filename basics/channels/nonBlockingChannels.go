package main

import "fmt"

func main(){
	c1 := make(chan string, 1)

	select{
	case msg := <- c1:
		fmt.Println("Received: ", msg)
	default:
		fmt.Println("Channel is empty or did not receive message")
	}

	select{
	case c1 <- "hi":
		fmt.Println("Received: ", <-c1)
	default:
		fmt.Println("Did not receive anything")
	}
}