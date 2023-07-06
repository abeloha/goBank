package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(CreateAccountRequestModel) (*Account, error)
	UpdateAccount(*Account) error
	DeleteAccount(int) error
	GetAccountByID(int) (*Account, error)
	GetAccounts(int) ([]*Account, error)
}

type PostgreStore struct {
	db *sql.DB
}

func (s *PostgreStore) Init() error {
	return s.CreateAccountTable()
}

func (s *PostgreStore) CreateAccountTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS accounts (
        id SERIAL PRIMARY KEY,
        first_name VARCHAR(255) NOT NULL,
		last_name VARCHAR(255) NOT NULL,
        number VARCHAR(255) NOT NULL,
        balance int NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	
	`
	_, err := s.db.Exec(query)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (s *PostgreStore) CreateAccount(accountReq CreateAccountRequestModel) (*Account, error) {

	account := NewAccount(accountReq)

	query := `
		insert into accounts (first_name, last_name, number, balance)
		values ($1, $2, $3, $4) RETURNING id, created_at;
	`

	err := s.db.QueryRow(query, account.FirstName, account.LastName, account.Number, account.Balance).Scan(&account.ID, &account.CreatedAt)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return account, nil
}

func (s *PostgreStore) UpdateAccount(account *Account) error {
	return nil
}
func (s *PostgreStore) DeleteAccount(id int) error {
	_, err := s.db.Exec("DELETE FROM accounts WHERE id = $1", id)
	return err
}
func (s *PostgreStore) GetAccountByID(id int) (*Account, error) {

	query := `
        select id, first_name, last_name, number, balance, created_at from accounts where id = $1
	`
	account := &Account{}

	err := s.db.QueryRow(query, id).Scan(&account.ID, &account.FirstName, &account.LastName, &account.Number, &account.Balance, &account.CreatedAt)

	return account, err
}

func (s *PostgreStore) GetAccounts(limit int) ([]*Account, error) {

	query := `
        select id, first_name, last_name, number, balance, created_at from accounts limit $1
	`
	rows, err := s.db.Query(query, limit)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	var accounts []*Account

	for rows.Next() {
		account, err := scanRows(rows)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		accounts = append(accounts, account)
	}

	err = rows.Err()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return accounts, nil
}

func scanRows(rows *sql.Rows) (*Account, error) {
	var account Account
	err := rows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.Balance,
		&account.CreatedAt,
	)

	return &account, err
}

func NewPostgreStore(user, dbname, password string) (*PostgreStore, error) {

	connStr := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", user, dbname, password)
	fmt.Println(connStr)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Error connecing:", err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &PostgreStore{
		db: db,
	}, nil
}
