package api

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	db "github.com/emmaahmads/summafy/db/sqlc"
	"github.com/emmaahmads/summafy/util"
	"github.com/gin-gonic/gin"
)

func (server *Server) HandlerUploadPage(c *gin.Context) {
	c.HTML(200, "uploadform.html", gin.H{
		"title": "Upload",
	})

	c.Header("Content-Type", "text/html")
}

func (server *Server) HandlerUploadDoc(c *gin.Context) {
	username_str := c.GetString("username")

	user, err := server.store.Queries.GetUser(context.Background(), username_str)
	if err != nil {
		util.MyGinLogger("Error fetching user:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	contentType := c.Request.Header.Get("Content-Type")
	if contentType != "" {
		contentType = strings.Split(contentType, ";")[0]
	}
	util.MyGinLogger("Content-Type:", contentType)

	if contentType == "multipart/form-data" {
		if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
			util.MyGinLogger("Error parsing form body:", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
			return
		}
	} else {
		util.MyGinLogger("Unsupported Content-Type")
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "Unsupported Content-Type"})
		return
	}

	file, header, err := c.Request.FormFile("files")
	if err != nil {
		util.MyGinLogger("Error retrieving file from request:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	// Create a new file in the local filesystem
	filePath := filepath.Join("uploads", header.Filename)
	dst, err := os.Create(filePath)
	if err != nil {
		util.MyGinLogger("Error creating file:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer dst.Close()

	// Copy the uploaded file to the destination file
	if _, err := io.Copy(dst, file); err != nil {
		util.MyGinLogger("Error saving file:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	local_file, err := os.Open(filePath)
	if err != nil {
		util.MyGinLogger("Error opening local file:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	util.MyGinLogger("Uploading file to S3")
	uploaded_file, err := server.UploadFileToS3("./upload/", local_file)
	if err != nil {
		util.MyGinLogger("Error uploading file to S3:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	os.Remove(filePath)

	util.MyGinLogger("Starting transaction for new document")
	doc, err := server.store.NewDocumentTx(c, db.NewDocumentParams{
		Username:   user.Username,
		IsPrivate:  false,
		HasSummary: true,
		FileName:   uploaded_file.filename,
		Param1:     false,
		Param2:     false,
		Summary:    string("No summary"),
	})

	if err != nil {
		util.MyGinLogger("Error creating new document transaction:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"document": doc.Document.ID})
}
