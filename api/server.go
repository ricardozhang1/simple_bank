package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/techschool/simplebank/db/sqlc"
	"log"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	store db.Store
	router *gin.Engine
}

// NewServer create a new HTTP server and setup router.
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	// 注册验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("currency", validCurrency)
		if err != nil {
			log.Fatalf("failed register currency validator, err: %v", err)
		}
	}

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)
	router.POST("/transfers", server.createTransfer)

	server.router = router
	return server
}

// Start run the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"message": err.Error()}
}



