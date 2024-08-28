package main

import (
	"fmt"
	"time"
)

func worker(done chan bool) {
	fmt.Println("Working...")
	time.Sleep(time.Second)
	fmt.Println("done")

	done <- true
}

func receiveMsg(ping chan<- string, msg string) {
	ping <- msg
}

func transferPing(ping1 <-chan string, ping2 chan<- string) {
	temp := <-ping1
	ping2 <- temp
}

func main() {
	// unbuffered channel
	messages := make(chan string)
	go func() {
		messages <- "ping"
	}()
	msg := <-messages
	fmt.Println(msg)

	// buffered channels
	messages2 := make(chan string, 2)
	messages2 <- "first message"
	messages2 <- "second message"

	fmt.Println(<-messages2)
	fmt.Println(<-messages2)

	// channel synchronization
	done := make(chan bool, 1)
	go worker(done)

	<-done

	// channel directions
	ping1 := make(chan string, 1)
	ping2 := make(chan string, 1)

	receiveMsg(ping1, "sent ping")
	fmt.Println(<-ping1)

	receiveMsg(ping1, "sent second ping")
	transferPing(ping1, ping2)
	fmt.Println(<-ping2)

	// select statements
	s1 := make(chan string)
	s2 := make(chan string)

	go func() {
		time.Sleep(3 * time.Second)
		s1 <- "message 1"
		fmt.Println("Send message to s1")
	}()

	go func() {
		time.Sleep(1 * time.Second)
		s2 <- "message 2"
		fmt.Println("Send message to s2")
	}()

	for i := range 2 {
		i++
		select {
		case msg1 := <-s1:
			fmt.Println("Received: ", msg1)
		case msg2 := <-s2:
			fmt.Println("Received: ", msg2)
		}
	}

}
