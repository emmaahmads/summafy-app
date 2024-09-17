package api

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	db "github.com/emmaahmads/summafy/db/sqlc"
	"github.com/emmaahmads/summafy/token"
	"github.com/emmaahmads/summafy/util"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

type Server struct {
	store      db.Store
	router     *gin.Engine
	aws        *aws.Config
	tokenMaker token.Maker
	s3_bucket  string
	apiKey     string
}

func NewServer(store db.Store, aws_conf *awsConfig, apiKey string, secret_key string) *Server {
	tokenMaker, err := token.NewJWTMaker(secret_key)
	if err != nil {
		log.Fatal("cannot create token maker:", err)
	}
	server := &Server{
		store: store,
		aws: &aws.Config{
			Region:      aws.String(aws_conf.region),
			Credentials: credentials.NewStaticCredentials(aws_conf.creds[0], aws_conf.creds[1], aws_conf.creds[2]),
		},
		tokenMaker: tokenMaker,
		s3_bucket:  aws_conf.s3_bucket,
		apiKey:     apiKey,
	}

	server.setupRouter()
	return server
}

func (server *Server) setupRouter() {
	r := gin.Default()
	r.Use(gin.Logger())
	r.LoadHTMLGlob("templates/*")
	r.StaticFile("/style.css", "templates/style.css")
	r.GET("/signup", server.HandlerSignupPage)
	r.POST("/signup", server.HandlerSignup)
	r.GET("/login", server.HandlerLoginPage)
	r.POST("/login", server.HandlerLogin)
	r.GET("/test", server.HandlerTest)
	util.MyGinLogger("Setting up auth group")
	r_auth := r.Group("/").Use(authMiddleware(server.tokenMaker))
	r_auth.GET("/", server.HandlerLandingPage)
	r_auth.GET("/dashboard", server.HandlerLandingPage)
	r.GET("/upload", server.HandlerUploadPage)
	r.GET("/view", server.HandlerViewDocuments)
	r.GET("/viewdoc/:id", server.HandlerViewDocumentsUploaded)
	r.POST("/upload", server.HandlerUploadDoc)
	r.GET("/download/:filename", server.HandlerDownloadDoc)
	r.POST("/notification", server.HandlerNotification)
	r.GET("/display-notified-objects", server.HandlerDisplayNotifiedObjects)
	util.MyGinLogger("Finished setting up routes")

	server.router = r
}
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		util.MyGinLogger("Received headers:", fmt.Sprintf("%+v", ctx.Request.Header))
		header := ctx.GetHeader(authorizationHeaderKey)
		/* if header == "" {
			util.MyGinLogger("Missing Authorization Header")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		} */

		// a header has two parts: "bearer <token>"
		util.MyGinLogger("Authorization header:", ctx.GetHeader(authorizationHeaderKey))
		fields := strings.Fields(header)
		if len(fields) < 2 {
			util.MyGinLogger("Authorization Header is not in the correct format")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if strings.ToLower(fields[0]) != authorizationTypeBearer {
			util.MyGinLogger("Authorization Header is not in the correct format")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			util.MyGinLogger("VerifyToken error:", string(err.Error()))
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}
