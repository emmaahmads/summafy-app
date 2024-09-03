package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Link struct {
	Href string
	Rel  string
}

func (server *Server) HandlerLandingPage(c *gin.Context) {
	link := Link{
		Href: "https://example.com/link",
		Rel:  "related",
	}
	c.HTML(200, "dashboard.html", gin.H{
		"title": "Emma Summafy - Home Page",
		"link":  link,
	})

	c.Header("Content-Type", "text/html")

}

func (server *Server) HandlerUploadPage(c *gin.Context) {
	c.HTML(200, "uploadform.html", gin.H{
		"emma": "Emma",
	})

	c.Header("Content-Type", "text/html")

	file, err := c.FormFile("document")
	if err != nil {
		return
	}

	gin.Logger()
	fmt.Println(file.Filename)

}

func (server *Server) HandlerViewPage(c *gin.Context) {
	c.HTML(200, "view.html", gin.H{
		"emma": "Emma",
	})

	c.Header("Content-Type", "text/html")

}

func (server *Server) HandlerUploadDoc(c *gin.Context) {
	file, err := c.FormFile("document")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dst := "./upload/" + file.Filename
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	local_file, err := os.Open(dst)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	message, err := server.UploadFileToS3("./upload/", local_file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": message})

}

func (server *Server) UploadFileToS3(fileDir string, file *os.File) (string, error) {
	sess, _ := session.NewSession(server.aws)

	svc := s3.New(sess)
	_, err := svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(server.s3_bucket),
		Key:    aws.String(file.Name()),
		Body:   file,
	})

	if err != nil {
		return "", err
	}

	return "https://emmaahmadsproject1.s3.amazonaws.com/" + file.Name(), nil
}
