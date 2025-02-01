package main

import "time"

type CreateAccountRequest struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name"`
	FullName    string `json:"full_name"`
	Gender      byte   `json:"gender"`
	IsHost      bool   `json:"is_host"`
}

type Account struct {
	ID          int       `json:"id"`
	Username    string    `json:"username"`
	Password    string    `json:"password"`
	DisplayName string    `json:"display_name"`
	FullName    string    `json:"full_name"`
	Gender      byte      `json:"gender"`
	IsHost      bool      `json:"is_host"`
	CreatedAt   time.Time `json:"created_at"`
}

type Event struct {
	ID              int       `json:"id"`
	HostID          int       `json:"host_id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Capacity        int       `json:"capacity"`
	StartAt         time.Time `json:"start_at"`
	EndAt           time.Time `json:"end_at"`
	LocationName    string    `json:"location_name"`
	LocationAddress string    `json:"location_address"`
	LocationCity    string    `json:"location_city"`
	LocationState   string    `json:"location_state"`
	LocationCountry string    `json:"location_country"`
	LocationZip     string    `json:"location_zip"`
	CreatedAt       time.Time `json:"created_at"`
}

type TicketType struct {
	ID                int       `json:"id"`
	EventID           int       `json:"event_id"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	Price             float64   `json:"price"`
	TotalQuantity     int       `json:"total_quantity"`
	AvailableQuantity int       `json:"available_quantity"`
	CreatedAt         time.Time `json:"created_at"`
}

func NewAccount(username string, password string, displayName string,
	fullName string, gender byte, isHost bool) *Account {
	return &Account{
		Username:    username,
		Password:    password,
		DisplayName: displayName,
		FullName:    fullName,
		Gender:      gender,
		IsHost:      isHost,
		CreatedAt:   time.Now(),
	}
}
