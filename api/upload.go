package api

import (
	"context"
	"fmt"
	"net/http"
	"os"

	db "github.com/emmaahmads/summafy/db/sqlc"
	"github.com/emmaahmads/summafy/util"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (server *Server) HandlerUploadPage(c *gin.Context) {
	c.HTML(200, "uploadform.html", gin.H{
		"title": "Upload",
	})

	c.Header("Content-Type", "text/html")
}

func (server *Server) HandlerUploadDoc(c *gin.Context) {

	username_str, ok := sessions.Default(c).Get("username").(string)
	util.MyGinLogger(username_str)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No username provided"})
		return
	}
	user, err := server.store.Queries.GetUser(context.Background(), username_str)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
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

	/* TODO for unit test */
	/* // get summary from open AI
	summary, err := server.SummarizeTextFile("./upload/test.txt")

	if err != nil {
		summary = "Summary could not be generated"
	}
	util.MyGinLogger(summary)

	// remove local copy */
	os.Remove(dst)

	//start transaction of New Document
	doc, err := server.store.NewDocumentTx(c, db.NewDocumentParams{
		Username:   user.Username,
		IsPrivate:  false,
		HasSummary: true,
		FileName:   uploaded_file.filename,
		Param1:     false,
		Param2:     false,
		Summary:    string(uploaded_file.summary),
	})

	c.Header("Content-Type", "text/html")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Redirect(http.StatusFound, fmt.Sprintf("/viewdoc/%d", doc.Document.ID))
}
