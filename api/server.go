package api

import (
	"github.com/Ritik1101-ux/simplebank/db/sqlc"
	"github.com/Ritik1101-ux/simplebank/token"
	"github.com/Ritik1101-ux/simplebank/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves Http request for our banking services

type Server struct {
	config     utils.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

//NewServer creates a new HTTP server and setup routing

func NewServer(config utils.Config, store db.Store) (*Server, error) {

	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmtericKey)

	if err != nil {
		return nil, err
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()

	return server, nil

}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/user", server.createUser)
	router.POST("/users/login", server.loginUser)

	authRoutes:=router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/account/:id", server.getAccount)
	authRoutes.GET("/accounts", server.ListAccounts)
	authRoutes.POST("/transfers", server.createTransfer)

	server.router = router
}

//Start runs the HTTP server on a Specific Address

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
