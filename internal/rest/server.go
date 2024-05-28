package rest

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kavinddd/mangtoon_be/internal/db"
	"github.com/kavinddd/mangtoon_be/internal/token"
	"github.com/kavinddd/mangtoon_be/pkg/util"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewJwtMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}
	server.setupRouter()
	return &server, nil
}

func (server *Server) Run(address string) error {
	return server.router.Run(address)
}
