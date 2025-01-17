package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	GetAccount() ([]*Account, error)
	DeleteAccountById(int) error
	GetAccountById(int) (*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=nitelife password=nighty sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	query, err := os.ReadFile("db/DDL.sql")
	if err != nil {
		panic(err)
	}
	if _, err := s.db.Exec(string(query)); err != nil {
		panic(err)
	}
	return nil
}

func (s *PostgresStore) CreateAccount(account *Account) error {
	query := `INSERT INTO account
	(Accountname, password, full_name, gender, is_host, created_at)
	VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := s.db.Query(
		query,
		account.Username,
		account.Password,
		account.FullName,
		account.Gender,
		account.IsHost, account.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStore) GetAccount() ([]*Account, error) {
	query := `SELECT * FROM account`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	accounts := []*Account{}
	for rows.Next() {
		account, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (s *PostgresStore) DeleteAccountById(id int) error {
	_, err := s.db.Exec("DELETE FROM account WHERE id = $1")
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStore) GetAccountById(id int) (*Account, error) {
	rows, err := s.db.Query("SELECT * FROM account WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("Account %d not found", id)
}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := new(Account)
	err := rows.Scan(
		&account.ID,
		&account.Username,
		&account.Password,
		&account.FullName,
		&account.Gender,
		&account.IsHost,
		&account.CreatedAt)

	if err != nil {
		return nil, err
	}

	return account, nil
}
