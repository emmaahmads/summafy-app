package api

import (
	"net/http"

	db "github.com/emmaahmads/summafy/db/sqlc"
	"github.com/gin-contrib/sessions"

	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store    db.Store
	router   *gin.Engine
	s3bucket string
}

func NewServer(store db.Store, s3bucket string) *Server {
	server := &Server{
		store:    store,
		s3bucket: s3bucket,
	}
	mycookie := cookie.NewStore([]byte("mysecretkey"))
	r := gin.Default()
	r.Use(gin.Logger())
	// Set Access-Control-Allow-Origin header
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Next()
	})

	r.Use(sessions.Sessions("mysession", mycookie))
	r.LoadHTMLGlob("templates/*")
	r.StaticFile("/style.css", "templates/style.css")
	r.GET("/signup", server.HandlerSignupPage)
	r.POST("/signup", server.HandlerSignup)
	r.GET("/login", server.HandlerLoginPage)
	r.GET("/", server.HandlerLoginPage)
	r.POST("/login", server.HandlerLogin)
	r.GET("/logout", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Delete("username")
		session.Save()
		c.Redirect(http.StatusFound, "/login")
	})

	// need token
	r.GET("/dashboard", server.HandlerLandingPage)
	r.GET("/upload", server.HandlerUploadPage)
	r.GET("/view", server.HandlerViewDocuments)
	r.GET("/viewdoc/:id", server.HandlerViewDocumentsUploaded)
	r.POST("/upload", server.HandlerUploadDoc)
	r.GET("/download/:filename", server.HandlerDownloadDoc)
	r.POST("/notification", server.HandlerNotification)
	r.POST("/createsummary", server.HandlerCreateSummary)
	r.GET("/createsummary", server.HandlerCreateSummarySubscription)
	r.GET("/display-notified-objects", server.HandlerDisplayNotifiedObjects)

	server.router = r
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
