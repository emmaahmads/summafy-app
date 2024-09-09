package api

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	db "github.com/emmaahmads/summafy/db/sqlc"
	"github.com/emmaahmads/summafy/util"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store     db.Store
	router    *gin.Engine
	aws       *aws.Config
	s3_bucket string
	apiKey    string
}

func NewServer(store db.Store, aws_conf *awsConfig, apiKey string) *Server {
	util.MyGinLogger(aws_conf.s3_bucket, aws_conf.region, aws_conf.creds[0], aws_conf.creds[1], aws_conf.creds[2])
	server := &Server{
		store: store,
		aws: &aws.Config{
			Region:      aws.String(aws_conf.region),
			Credentials: credentials.NewStaticCredentials(aws_conf.creds[0], aws_conf.creds[1], aws_conf.creds[2]),
		},
		s3_bucket: aws_conf.s3_bucket,
		apiKey:    apiKey,
	}

	r := gin.Default()
	r.Use(gin.Logger())
	r.LoadHTMLGlob("templates/*")
	r.StaticFile("/style.css", "templates/style.css")
	r.GET("/", server.HandlerLandingPage)
	r.GET("/dashboard", server.HandlerLandingPage)
	r.GET("/upload", server.HandlerUploadPage)
	r.GET("/view", server.HandlerViewDocuments)
	r.GET("/viewdoc/:id", server.HandlerViewDocumentsUploaded)
	r.POST("/upload", server.HandlerUploadDoc)
	r.GET("/download/:filename", server.HandlerDownloadDoc)
	r.POST("/notification", server.HandlerNotification)
	r.GET("/display-notified-objects", server.HandlerDisplayNotifiedObjects)

	server.router = r
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
