package api

import (
	"github.com/Ritik1101-ux/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves Http request for our banking services

type Server struct {
	store  db.Store
	router *gin.Engine
}

//NewServer creates a new HTTP server and setup routing

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}	

	router.POST("/accounts", server.createAccount)
	router.GET("/account/:id", server.getAccount)
	router.GET("/accounts", server.ListAccounts)

	router.POST("/transfers", server.createTransfer)

	router.POST("/user", server.createUser)
	server.router = router

	return server

}

//Start runs the HTTP server on a Specific Address

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
