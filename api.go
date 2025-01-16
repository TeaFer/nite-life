package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

type APIServer struct {
	listenAddr string
}

func NewAPIServer(listenAddr string) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
	}
}

func (s *APIServer) Run() {
	router := gin.Default()

	log.Println("JSON API server running on port:", s.listenAddr)
	router.Run(s.listenAddr)
}
