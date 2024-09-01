package api

import (
	db "github.com/emmaahmads/summafy/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  db.Store
	router *gin.Engine
	// sdk for aws s3
	// openai
}

func NewServer(store db.Store) *Server {
	server := &Server{
		store: store,
	}
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

	r.GET("/", server.HandlerLandingPage)
	r.GET("/test", server.HandlerLandingPageTest)

	server.router = r
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
