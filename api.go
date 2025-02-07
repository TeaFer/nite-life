package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/TeaFer/nite-life/auth"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
	router.Any("/accounts/:id", makeHandlerFuncMiddleware(s.handleProtectedResource),
		makeHandlerFunc(s.handleAccountById))
	router.Any("/accounts/:id/tickets", makeHandlerFuncMiddleware(s.handleProtectedResource),
		makeHandlerFunc(s.handleGetTicketsByAccountId))

	router.Any("/events", makeHandlerFunc(s.handleEvent))
	router.Any("/events/:id", makeHandlerFunc(s.handleEventById))

	log.Println("JSON API server running on port:", s.listenAddr)
	router.Run(s.listenAddr)
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
	id, err := getID(c)
	if err != nil {
		return err
	}

	account, err := s.store.GetAccountById(id)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, account)
	return nil
}

func (s *APIServer) handleDeleteAccountById(c *gin.Context) error {
	id, err := getID(c)
	if err != nil {
		return err
	}
	if err := s.store.DeleteAccountById(id); err != nil {
		return err
	}
	c.Status(http.StatusOK)
	return nil
}

func (s *APIServer) handleGetTicketsByAccountId(c *gin.Context) error {
	id, err := getID(c)
	if err != nil {
		return err
	}

	tickets, err := s.store.GetTicketsByAccountId(id)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, tickets)
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
	id, err := getID(c)
	if err != nil {
		return err
	}

	event, err := s.store.GetEventById(id)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, event)
	return nil
}

func getID(c *gin.Context) (int, error) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid id provided: %s", idStr)
	}

	return id, nil
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

func makeHandlerFuncMiddleware(f apiFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := f(c)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusBadRequest,
				apiError{
					Error:   "Bad request. Please try again.",
					Message: err.Error(),
				})
			c.Abort()
		}
		c.Next()
	}
}

func (s *APIServer) handleProtectedResource(c *gin.Context) error {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		return fmt.Errorf("authorization token is missing")
	}

	err := godotenv.Load()
	if err != nil {
		return err
	}
	jwtKey := []byte(os.Getenv("JWT_KEY"))

	claims, err := auth.ValidateJWT(tokenString, jwtKey)
	if err != nil {
		return err
	}
	resourceId, err := getID(c)
	if err != nil {
		return err
	}

	if err := s.checkClaims(claims, resourceId); err != nil {
		return err
	}

	c.Set("claims", claims)
	return nil
}

func (s *APIServer) checkClaims(claims *auth.Claims, resourceId int) error {
	if claims.ID == 0 || claims.ID != resourceId {
		return fmt.Errorf("you are unauthorized to make this request")
	}
	return nil
}
