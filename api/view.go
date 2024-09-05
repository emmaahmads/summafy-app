package api

import (
	"context"
	"net/http"
	"time"

	"github.com/emmaahmads/summafy/util"
	"github.com/gin-gonic/gin"
)

type DocumentView struct {
	Filename  string
	UserId    string
	Size      string
	CreatedAt time.Time
	Summary   string
	Download  string
}

func (server *Server) HandlerViewDocuments(c *gin.Context) {

	util.MyGinLogger("HEYYYYYYY")
	c.HTML(200, "view.html", gin.H{})

	c.Header("Content-Type", "text/html")
}

type getDocIdRequest struct {
	DocId int `uri:"id" binding:"required"`
}

func (server *Server) HandlerViewDocumentsUploaded(c *gin.Context) {
	var req getDocIdRequest
	var document []DocumentView
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	doc, err := server.store.Queries.GetDocument(context.Background(), int64(req.DocId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	sum, _ := server.store.Queries.GetSummary(context.Background(), doc.ID)

	download := "/download/" + doc.FileName
	doc_info := DocumentView{
		Filename:  doc.FileName,
		UserId:    doc.Username,
		Size:      "None",
		CreatedAt: doc.CreatedAt,
		Summary:   string(sum.Summary),
		Download:  download,
	}
	document = append(document, doc_info)
	c.HTML(200, "view.html", gin.H{
		"document": document,
	})

	c.Header("Content-Type", "text/html")
}
