package api

import (
	"github.com/gin-gonic/gin"
	"github.com/kavinddd/mangtoon_be/internal/db"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := Server{store: store}
	server.SetupRoutes()
	return &server
}

func (server *Server) Run(address string) error {
	return server.router.Run(address)
}
