package main

import "time"

type CreateAccountRequest struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	DisplayName string `json:"display_name" binding:"required"`
	FullName    string `json:"full_name" binding:"required"`
	Gender      byte   `json:"gender" binding:"required"`
	IsHost      bool   `json:"is_host" binding:"required"`
}

type CreateEventRequest struct {
	HostID          int       `json:"host_id" binding:"required"`
	Name            string    `json:"name" binding:"required"`
	Description     string    `json:"description" binding:"required"`
	Capacity        int       `json:"capacity" binding:"required"`
	StartAt         time.Time `json:"start_at" binding:"required"`
	EndAt           time.Time `json:"end_at" binding:"required"`
	LocationName    string    `json:"location_name" binding:"required"`
	LocationAddress string    `json:"location_address" binding:"required"`
	LocationCity    string    `json:"location_city" binding:"required"`
	LocationState   string    `json:"location_state" binding:"required"`
	LocationCountry string    `json:"location_country" binding:"required"`
	LocationZip     string    `json:"location_zip" binding:"required"`
}

type Account struct {
	ID          int       `json:"id"`
	Username    string    `json:"username"`
	Password    string    `json:"-"`
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

type Ticket struct {
	ID              int       `json:"id"`
	TicketTypeID    int       `json:"ticket_type_id"`
	OwnerId         int       `json:"owner_id"`
	PurchasedAt     time.Time `json:"purchased_at"`
	TicketTypeName  string    `json:"ticket_type_name"`
	TicketTypePrice float32   `json:"ticket_type_price"`
	EventID         int       `json:"event_id"`
	EventName       string    `json:"event_name"`
	EventStartAt    time.Time `json:"event_start_at"`
	EventEndAt      time.Time `json:"event_end_at"`
}

// Participant represents a user who is attending an event. It contains
// information of the user, the event, and the ticket they purchased to
// attend the event.
type Participant struct {
	EventID         int       `json:"event_id"`
	UserID          int       `json:"user_id"`
	DisplayName     string    `json:"display_name"`
	FullName        string    `json:"full_name"`
	Gender          string    `json:"gender"`
	TicketID        int       `json:"ticket_id"`
	PurchasedAt     time.Time `json:"purchased_at"`
	TicketTypeID    int       `json:"ticket_type_id"`
	TicketTypeName  string    `json:"ticket_type_name"`
	TicketTypePrice float32   `json:"ticket_type_price"`
}

func NewEvent(hostId int, name string, description string, capacity int, startAt time.Time,
	endAt time.Time, locationName string, locationAddress string, locationCity string,
	locationState string, locationCountry string, locationZip string) *Event {
	return &Event{
		HostID:          hostId,
		Name:            name,
		Description:     description,
		Capacity:        capacity,
		StartAt:         startAt,
		EndAt:           endAt,
		LocationName:    locationName,
		LocationAddress: locationAddress,
		LocationCity:    locationCity,
		LocationState:   locationState,
		LocationCountry: locationCountry,
		LocationZip:     locationZip,
		CreatedAt:       time.Now(),
	}
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
