package api

import (
	"time"

	db "github.com/emmaahmads/summafy/db/sqlc"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	r := gin.Default()
	r.Use(gin.Logger())
	// Set Access-Control-Allow-Origin header
	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8081"}, /* TODO - Use env var in production for container bridged networking */
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.POST("/signup", server.HandlerSignup)
	r.POST("/login", server.HandlerLogin)
	r.GET("/logout", func(c *gin.Context) {
		// placeholder
	})

	api := r.Group("/api/v1")
	api.Use(server.middlewareAuth())
	{
		api.GET("/dashboard", server.HandlerDashboard)
		api.GET("/viewdoc/:id", server.HandlerViewDocumentsUploaded)
		api.POST("/upload", server.HandlerUploadDoc)
		api.GET("/download/:filename", server.HandlerDownloadDoc)
		api.POST("/notification", server.HandlerNotification)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	server.router = r
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
