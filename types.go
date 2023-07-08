package main

import (
	"fmt"
	"math/rand"
	"time"
)

type CreateAccountRequestModel struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type TransferRequestModel struct {
	AccountNumber string `json:"account_number"`
	Amount        int    `json:"amount"`
	Remarks       string `json:"remarks"`
}

type Account struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Number    string    `json:"number"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

func NewAccount(data CreateAccountRequestModel) *Account {
	return &Account{
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Number:    fmt.Sprint(rand.Intn(100000000)),
		CreatedAt: time.Now().UTC(),
	}
}
