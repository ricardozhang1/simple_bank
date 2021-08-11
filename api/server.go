package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/token"
	"github.com/techschool/simplebank/util"
	"log"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	config		util.Config
	store 		db.Store
	tokenMaker 	token.Maker
	router 		*gin.Engine  // 初始化时，并不传入这个参数，在gin.Default()得到*gin.Engine后传入
}

// NewServer create a new HTTP server and setup router.
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker, err: %v", err)
	}

	server := &Server{
		config: config,
		store: store,
		tokenMaker: tokenMaker,
	}

	// currency注册验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("currency", validCurrency)
		if err != nil {
			log.Fatalf("failed register currency validator, err: %v", err)
		}
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter()  {
	router := gin.Default()

	authRouter := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRouter.POST("/accounts", server.createAccount)
	authRouter.GET("/accounts/:id", server.getAccount)
	authRouter.GET("/accounts", server.listAccount)
	authRouter.POST("/transfers", server.createTransfer)

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	server.router = router
}

// Start run the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// errorResponse 统一处理错误信息返回的格式，方便阅读，简化代码
func errorResponse(err error) gin.H {
	return gin.H{"message": err.Error()}
}



