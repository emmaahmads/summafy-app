package api

import (
	"net/http"
	"time"

	db "github.com/emmaahmads/summafy/db/sqlc"
	"github.com/gin-contrib/sessions"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

type Server struct {
	secretKey string
	store     db.Store
	router    *gin.Engine
	s3bucket  string
}

func NewServer(store db.Store, s3bucket string, secretkey string) *Server {
	server := &Server{
		store:     store,
		s3bucket:  s3bucket,
		secretKey: secretkey,
	}
	mycookie := cookie.NewStore([]byte("mysecretkey"))
	r := gin.Default()
	r.Use(gin.Logger())
	// Set Access-Control-Allow-Origin header
	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8081"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

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

	api := r.Group("/api/v1")
	api.Use(server.middlewareAuth())
	{
		api.GET("/dashboard", server.HandlerDashboard)
		api.GET("/upload", server.HandlerUploadPage)
		api.GET("/view", server.HandlerViewDocuments)
		api.GET("/viewdoc/:id", server.HandlerViewDocumentsUploaded)
		api.POST("/upload", server.HandlerUploadDoc)
		api.GET("/download/:filename", server.HandlerDownloadDoc)
		api.POST("/notification", server.HandlerNotification)
	}

	server.router = r
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
