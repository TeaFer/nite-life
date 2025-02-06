package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type APIServer struct {
	listenAddr string
	store      Storage
}

type apiFunc func(*gin.Context) error

type apiError struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	router := gin.Default()

	router.Any("/accounts", makeHandlerFunc(s.handleAccount))
	router.Any("/accounts/:id", makeHandlerFunc(s.handleAccountById))
	router.Any("/events", makeHandlerFunc(s.handleEvent))
	router.Any("/events/:id", makeHandlerFunc(s.handleEventById))

	log.Println("JSON API server running on port:", s.listenAddr)
	router.Run(s.listenAddr)
}

func makeHandlerFunc(f apiFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := f(c)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusBadRequest,
				apiError{
					Error:   "Bad request. Please try again.",
					Message: err.Error(),
				})
		}
	}
}

func (s *APIServer) handleAccount(c *gin.Context) error {
	switch c.Request.Method {
	case "GET":
		return s.handleGetAccount(c)
	case "POST":
		return s.handleCreateAccount(c)
	default:
		return fmt.Errorf("method not supported: %s", c.Request.Method)
	}
}

func (s *APIServer) handleAccountById(c *gin.Context) error {
	switch c.Request.Method {
	case "GET":
		return s.handleGetAccountById(c)
	case "DELETE":
		return s.handleDeleteAccountById(c)
	default:
		return fmt.Errorf("method not supported: %s", c.Request.Method)
	}
}

func (s *APIServer) handleEvent(c *gin.Context) error {
	switch c.Request.Method {
	case "GET":
		return s.handleGetEvent(c)
	case "POST":
		return s.handleCreateEvent(c)
	}
	return nil
}

func (s *APIServer) handleEventById(c *gin.Context) error {
	switch c.Request.Method {
	case "GET":
		return s.handleGetEventById(c)
	}
	return nil
}

func (s *APIServer) handleGetAccount(c *gin.Context) error {
	accounts, err := s.store.GetAccount()
	if err != nil {
		return err
	}
	c.JSON(200, accounts)
	return nil
}

func (s *APIServer) handleCreateAccount(c *gin.Context) error {
	createAccountReq := new(CreateAccountRequest)
	c.BindJSON(createAccountReq)
	Account := NewAccount(
		createAccountReq.Username,
		createAccountReq.Password,
		createAccountReq.DisplayName,
		createAccountReq.FullName,
		createAccountReq.Gender, createAccountReq.IsHost)
	err := s.store.CreateAccount(Account)
	if err != nil {
		return err
	}

	c.JSON(http.StatusCreated, Account)
	return nil
}

func (s *APIServer) handleGetAccountById(c *gin.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("invalid id provided: %s", idStr)
	}

	account, err := s.store.GetAccountById(id)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, account)
	return nil
}

func (s *APIServer) handleDeleteAccountById(c *gin.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("invalid id provided: %s", idStr)
	}
	if err := s.store.DeleteAccountById(id); err != nil {
		return err
	}
	c.Status(200)
	return nil
}

func (s *APIServer) handleGetEvent(c *gin.Context) error {
	events, err := s.store.GetEvent()
	if err != nil {
		return err
	}
	c.JSON(200, events)
	return nil
}

func (s *APIServer) handleCreateEvent(c *gin.Context) error {
	createEventReq := new(CreateEventRequest)
	c.BindJSON(createEventReq)
	Event := NewEvent(
		createEventReq.HostID,
		createEventReq.Name,
		createEventReq.Description,
		createEventReq.Capacity,
		createEventReq.StartAt,
		createEventReq.EndAt,
		createEventReq.LocationName,
		createEventReq.LocationAddress,
		createEventReq.LocationCity,
		createEventReq.LocationState,
		createEventReq.LocationCountry,
		createEventReq.LocationZip)
	err := s.store.CreateEvent(Event)
	if err != nil {
		return err
	}

	c.JSON(http.StatusCreated, Event)
	return nil
}

func (s *APIServer) handleGetEventById(c *gin.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("invalid id provided: %s", idStr)
	}

	event, err := s.store.GetEventById(id)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, event)
	return nil
}
