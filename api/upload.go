package api

import (
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	db "github.com/emmaahmads/summafy/db/sqlc"
	"github.com/emmaahmads/summafy/util"
	"github.com/gin-gonic/gin"
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

	uploaded_file, err := server.UploadFileToS3("./upload/", local_file)

	// remove local copy
	os.Remove(dst)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//start transaction of New Document
	_, err = server.store.NewDocumentTx(c, db.NewDocumentParams{
		Username:   "emma",
		IsPrivate:  false,
		HasSummary: true,
		FileName:   uploaded_file,
		Param1:     false,
		Param2:     false,
		Summary:    "no summary no summary no summary thank you",
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": uploaded_file})

}

func (server *Server) UploadFileToS3(fileDir string, file *os.File) (string, error) {
	sess, err := session.NewSession(server.aws)

	if err != nil {
		return "", err
	}

	filename := strings.Split(file.Name(), "/")[2]
	util.MyGinLogger(filename)
	svc := s3.New(sess)
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(server.s3_bucket),
		Key:    aws.String(filename),
		Body:   file,
		// todo add tags for username
	})

	if err != nil {
		return "", err
	}

	return filename, nil
}
