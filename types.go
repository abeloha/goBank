package main

import (
	"fmt"
	"math/rand"
)

type Account struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Number    string `json:"number"`
	Balance   int64  `json:"balance"`
}

func NewAccount(FirstName, LastName string) *Account {
	return &Account{
		ID:        rand.Intn(10000),
		FirstName: FirstName,
		LastName:  LastName,
		Number:    fmt.Sprint(rand.Intn(100000000)),
	}
}
