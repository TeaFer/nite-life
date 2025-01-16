package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

type APIServer struct {
	listenAddr string
}

type apiFunc func(*gin.Context) error

type apiError struct {
	Error string `json:error`
}

func NewAPIServer(listenAddr string) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
	}
}

func (s *APIServer) Run() {
	router := gin.Default()

	router.Any("/user", makeHandlerFunc(s.handleUser))

	log.Println("JSON API server running on port:", s.listenAddr)
	router.Run(s.listenAddr)
}

func makeHandlerFunc(f apiFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := f(c)
		if err != nil {
			c.Error(err)
			c.JSON(400, apiError{Error: err.Error()})
		}
	}
}

func (s *APIServer) handleUser(c *gin.Context) error {
	switch c.Request.Method {
	case "GET":
		return s.handleGetUser(c)
	case "POST":
		return s.handleCreateUser(c)
	case "DELETE":
		return s.handleDeleteUser(c)
	default:
		return fmt.Errorf("method not supported: %s", c.Request.Method)
	}
}

func (s *APIServer) handleGetUser(c *gin.Context) error {

}

func (s *APIServer) handleCreateUser(c *gin.Context) error {

}

func (s *APIServer) handleDeleteUser(c *gin.Context) error {

}
