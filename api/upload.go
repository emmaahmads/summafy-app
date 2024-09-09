package api

import (
	"fmt"
	"net/http"
	"os"

	db "github.com/emmaahmads/summafy/db/sqlc"
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
		Username:   "emma",
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
