package main

import "time"

type CreateAccountRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	Gender   byte   `json:"gender"`
	IsHost   bool   `json:"is_host"`
}

type Account struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	FullName  string    `json:"full_name"`
	Gender    byte      `json:"gender"`
	IsHost    bool      `json:"is_host"`
	CreatedAt time.Time `json:"created_at"`
}

func NewAccount(username string, password string,
	fullName string, gender byte, isHost bool) *Account {
	return &Account{
		Username:  username,
		Password:  password,
		FullName:  fullName,
		Gender:    gender,
		IsHost:    isHost,
		CreatedAt: time.Now(),
	}
}
