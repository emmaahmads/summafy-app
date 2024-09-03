package api

import (
	"net/http"
	"os"
	"strings"

	"github.com/emmaahmads/summafy/util"
	"github.com/gin-gonic/gin"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

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
	os.Remove(dst)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": message})

}

func (server *Server) UploadFileToS3(fileDir string, file *os.File) (string, error) {
	sess, err := session.NewSession(server.aws)

	if err != nil {
		return "", err
	}

	filename := strings.Split(file.Name(), "/")[2]
	util.MyLogger(filename)
	svc := s3.New(sess)
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(server.s3_bucket),
		Key:    aws.String(filename),
		Body:   file,
	})

	if err != nil {
		return "", err
	}

	return filename, nil
}
