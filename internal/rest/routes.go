package rest

import (
	"github.com/gin-gonic/gin"
)

func (server *Server) setupRouter() {
	router := gin.Default()
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	adminOnly := router.Group("/").Use(adminMiddleware(server.tokenMaker))

	// users
	router.POST("/users", server.createUser)
	adminOnly.GET("/users", server.listUsers)
	authRoutes.GET("/users/:id", server.findUserById)

	// auth
	router.POST("/login", server.login)
	router.POST("/logout", server.logout)
	router.POST("/tokens/renew_access", server.renewAccessToken)

	server.router = router
}
