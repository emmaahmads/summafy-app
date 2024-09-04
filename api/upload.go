package api

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	db "github.com/emmaahmads/summafy/db/sqlc"
	"github.com/emmaahmads/summafy/util"
	"github.com/gin-gonic/gin"
)

func (server *Server) HandlerUploadPage(c *gin.Context) {
	c.HTML(200, "uploadform.html", gin.H{
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

	uploaded_file, err := server.UploadFileToS3("./upload/", local_file)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// get summary from open AI
	summary, err := server.SummarizeTextFile("./upload/test.txt")

	if err != nil {
		summary = "Summary could not be generated"
	}
	util.MyGinLogger(summary)
	// remove local copy
	os.Remove(dst)

	//start transaction of New Document
	doc, err := server.store.NewDocumentTx(c, db.NewDocumentParams{
		Username:   "emma",
		IsPrivate:  false,
		HasSummary: true,
		FileName:   uploaded_file,
		Param1:     false,
		Param2:     false,
		Summary:    summary,
	})

	c.Header("Content-Type", "text/html")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Redirect(http.StatusFound, "/view?doc_id="+strconv.Itoa(int(doc.Document.ID)))
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
