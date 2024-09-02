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
	r.Use(gin.Logger())
	r.LoadHTMLGlob("templates/*")
	r.StaticFile("/style.css", "templates/style.css")
	r.GET("/", server.HandlerLandingPage)
	r.GET("/upload", server.HandlerUploadPage)
	r.GET("/view", server.HandlerViewPage)
	r.POST("/upload", server.HandlerUploadDoc)

	server.router = r
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
