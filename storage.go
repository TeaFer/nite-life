package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	GetAccount() ([]*Account, error)
	DeleteAccountById(int) error
	GetAccountById(int) (*Account, error)
	GetEvent() ([]*Event, error)
	CreateEvent(*Event) error
	GetEventById(int) (*Event, error)
	GetTicketsByAccountId(int) ([]*Ticket, error)
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
	query := `SELECT id, username, display_name, full_name, gender, is_host, created_at FROM account`

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
	rows, err := s.db.Query(`SELECT id, username, display_name, full_name, gender, 
	is_host, created_at FROM account WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("Account %d not found", id)
}

func (s *PostgresStore) GetTicketsByAccountId(id int) ([]*Ticket, error) {
	query := `SELECT ticket.id, ticket.ticket_type_id, ticket.owner_id, ticket.purchased_at,
	ticket_type.name, ticket_type.price, event.id, event.name, event.start_at, event.end_at 
	FROM ticket 
	join ticket_type on ticket.ticket_type_id = ticket_type.id 
	join event on ticket_type.event_id = event.id WHERE ticket.owner_id = $1`
	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}

	tickets := []*Ticket{}
	for rows.Next() {
		ticket, err := scanIntoTicket(rows)
		if err != nil {
			return nil, err
		}

		tickets = append(tickets, ticket)
	}
	return tickets, nil
}

func (s *PostgresStore) GetEvent() ([]*Event, error) {
	query := `SELECT id, host_id, name, description, capacity, 
	start_at, end_at, location_name, location_address, location_city, 
	location_state, location_zip, created_at
	FROM event`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	events := []*Event{}
	for rows.Next() {
		event, err := scanIntoEvent(rows)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

func (s *PostgresStore) CreateEvent(event *Event) error {
	query := `INSERT INTO event
	(host_id, name, description, capacity, start_at, end_at, 
	location_name, location_address, location_city, location_state,
	location_country, location_zip, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

	_, err := s.db.Query(
		query,
		event.HostID,
		event.Name,
		event.Description,
		event.Capacity,
		event.StartAt,
		event.EndAt,
		event.LocationName,
		event.LocationAddress,
		event.LocationCity,
		event.LocationState,
		event.LocationCountry,
		event.LocationZip,
		event.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStore) GetEventById(id int) (*Event, error) {
	query := `SELECT id, host_id, name, description, capacity, 
	start_at, end_at, location_name, location_address, location_city, 
	location_state, location_zip, created_at
	FROM event`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanIntoEvent(rows)
	}

	return nil, fmt.Errorf("Event %d not found", id)
}

func scanIntoEvent(rows *sql.Rows) (*Event, error) {
	event := new(Event)
	err := rows.Scan(
		&event.ID,
		&event.HostID,
		&event.Name,
		&event.Description,
		&event.Capacity,
		&event.StartAt,
		&event.EndAt,
		&event.LocationName,
		&event.LocationAddress,
		&event.LocationCity,
		&event.LocationState,
		&event.LocationCountry,
		&event.LocationZip,
		&event.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return event, nil
}

func scanIntoTicket(rows *sql.Rows) (*Ticket, error) {
	ticket := new(Ticket)
	err := rows.Scan(
		&ticket.ID,
		&ticket.TicketTypeID,
		&ticket.OwnerId,
		&ticket.PurchasedAt,
		&ticket.TicketTypeName,
		&ticket.TicketTypePrice,
		&ticket.EventID,
		&ticket.EventName,
		&ticket.EventStartAt,
		&ticket.EventEndAt,
	)
	if err != nil {
		return nil, err
	}

	return ticket, nil
}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := new(Account)
	err := rows.Scan(
		&account.ID,
		&account.Username,
		&account.Password,
		&account.DisplayName,
		&account.FullName,
		&account.Gender,
		&account.IsHost,
		&account.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (s PostgresStore) addFilters(query string, filters string) (string, []interface{}, error) {
	if len(filters) == 0 {
		return query, nil, nil
	}

	// parse filter query parameter string to a map with fields as keys and filter values as values
	params, err := url.ParseQuery(filters)
	if err != nil {
		return query, nil, fmt.Errorf("failed to parse filter: %v", err)
	}

	// construct the where clause of the query based on the filters
	var conditions []string
	var values []interface{}
	i := 1
	for field, val := range params {
		conditions = append(conditions, fmt.Sprintf("%s = $%d", field, i))
		values = append(values, val[0])
		i++
	}

	// append the where clause and conditions to the base query
	query += " WHERE "
	query += strings.Join(conditions, " AND ")

	return query, values, nil
}

func (s PostgresStore) addSorts(query string, sorts string) (string, error) {
	if len(sorts) == 0 {
		return query, nil
	}

}
