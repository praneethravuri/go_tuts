package main

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v7"
	"math/rand"
	"sync"
	"time"
)

type Article struct {
	Topic   string
	Content string
}

type User struct {
	Name               string
	Email              string
	SubscribedArticles []Article
	mu                 sync.Mutex
}

func generateArticle() Article {
	return Article{
		Topic:   gofakeit.Quote(),
		Content: gofakeit.Paragraph(2, 5, 10, "."),
	}
}

func generateUser() User {
	return User{
		Name:               gofakeit.Name(),
		Email:              gofakeit.Email(),
		SubscribedArticles: []Article{},
	}
}

func main() {
	users := generateRandomUsers(10)

	var wg sync.WaitGroup

	articleChan := make(chan Article, 100)
	done := make(chan bool)

	wg.Add(2)
	go generateArticlesRandomly(articleChan, done, &wg)
	go SubscribeUser(articleChan, users, &wg)

	wg.Wait()

	close(articleChan)

	for _, user := range users {
		fmt.Printf("User: %s | Subscribed Articles:\n", user.Name)
		for _, article := range user.SubscribedArticles {
			fmt.Printf(" - %s\n", article.Topic)
		}
		fmt.Println()
	}
}

func generateRandomUsers(count int) []User {
	users := make([]User, count)
	for i := 0; i < count; i++ {
		users[i] = generateUser()
	}
	return users
}

func generateArticlesRandomly(articleChan chan<- Article, done chan<- bool, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 10; i++ {
		article := generateArticle()

		select {
		case articleChan <- article:
			fmt.Println("New article generated")
		default:
			fmt.Println("Channel full, discarding article")
		}
		randomDuration := time.Duration(rand.Intn(5)+1) * time.Second
		time.Sleep(randomDuration)
	}
	done <- true
}

func SubscribeUser(articleChan <-chan Article, users []User, wg *sync.WaitGroup) {
	defer wg.Done()
	for article := range articleChan {
		randomUserIdx := rand.Intn(len(users))
		user := &users[randomUserIdx]

		user.mu.Lock()
		user.SubscribedArticles = append(user.SubscribedArticles, article)
		user.mu.Unlock()

		fmt.Printf("User %s subscribed to article: %s\n", user.Name, article.Topic)
	}
}
