package api

import (
	"github.com/gin-gonic/gin"
)

func (server *Server) SetupRoutes() {
	router := gin.Default()

	// users
	router.POST("/users", server.CreateUser)
	router.GET("/users", server.ListUsers)
	router.GET("/users/:id", server.FindUserById)

	// auth
	router.POST("/login", server.HandleLogin)
	router.POST("/logout", server.HandleLogout)

	server.router = router
}
