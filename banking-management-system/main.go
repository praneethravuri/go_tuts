package main

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"math/rand/v2"
	"sync"
	"time"
)

var accountsMu sync.Mutex

type TransactionHistory struct {
	TransactionDate time.Time
	TransactionType string
	Amount          float64
}

type Account struct {
	AccountId string
	Name      string
	Email     string
	Balance   float64
	History   []TransactionHistory
}

type AccountManager interface {
	Deposit(amount float64)
	Withdraw(amount float64) error
	GetBalance() float64
	AddTransactionHistory(transactionType string, amount float64)
}

func (a *Account) Deposit(amount float64) {
	a.Balance += amount
	a.AddTransactionHistory("Deposit", amount)
}

func (a *Account) Withdraw(amount float64) error {
	if a.Balance < amount {
		return fmt.Errorf("insufficient funds")
	}
	a.Balance -= amount
	a.AddTransactionHistory("Withdrawal", amount)
	return nil
}

func (a *Account) GetBalance() float64 {
	return a.Balance
}

func (a *Account) AddTransactionHistory(transactionType string, amount float64) {
	a.History = append(a.History, TransactionHistory{
		TransactionDate: time.Now(),
		TransactionType: transactionType,
		Amount:          amount,
	})
}

func GenerateAccounts(accounts *[]Account, wg *sync.WaitGroup) {
	defer wg.Done()
	newAccount := Account{
		AccountId: uuid.New().String(),
		Name:      gofakeit.Name(),
		Email:     gofakeit.Email(),
		Balance:   0.0,
		History:   make([]TransactionHistory, 0),
	}
	accountsMu.Lock()
	randomDelay := time.Duration(rand.IntN(10) + 1)
	time.Sleep(randomDelay * time.Second)
	*accounts = append(*accounts, newAccount)
	accountsMu.Unlock()
}

func PerformTransaction(accounts *[]Account, wg *sync.WaitGroup, transactionType string) {
	defer wg.Done()
	accountsMu.Lock()
	numberOfAccounts := len(*accounts)
	if numberOfAccounts == 0 {
		accountsMu.Unlock()
		return
	}
	randomAccountIdx := rand.IntN(numberOfAccounts)
	user := &(*accounts)[randomAccountIdx]
	accountsMu.Unlock()

	var manager AccountManager = user
	randomAmount := rand.Float64() * 100

	time.Sleep(time.Duration(rand.IntN(3)) * time.Second)

	accountsMu.Lock()
	defer accountsMu.Unlock()

	switch transactionType {
	case "Deposit":
		manager.Deposit(randomAmount)
		fmt.Printf("Deposited $%.2f to account ID: %s | New Balance: $%.2f\n", randomAmount, user.AccountId, manager.GetBalance())
	case "Withdrawal":
		err := manager.Withdraw(randomAmount)
		if err != nil {
			fmt.Printf("Failed to withdraw $%.2f from account ID: %s | Error: %s\n", randomAmount, user.AccountId, err)
		} else {
			fmt.Printf("Withdrew $%.2f from account ID: %s | New Balance: $%.2f\n", randomAmount, user.AccountId, manager.GetBalance())
		}
	}
}

func main() {
	var wg sync.WaitGroup
	accounts := []Account{}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go GenerateAccounts(&accounts, &wg)
	}
	wg.Wait()

	for _, account := range accounts {
		fmt.Printf("Account ID: %s | Name: %s | Email: %s\n", account.AccountId, account.Name, account.Email)
	}

	for i := 0; i < 20; i++ {
		wg.Add(1)
		if rand.IntN(2) == 0 {
			go PerformTransaction(&accounts, &wg, "Deposit")
		} else {
			go PerformTransaction(&accounts, &wg, "Withdrawal")
		}
	}
	wg.Wait()

	for _, account := range accounts {
		fmt.Printf("Account ID: %s | Final Balance: $%.2f\n", account.AccountId, account.Balance)
	}
}
