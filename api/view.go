package api

import (
	"context"
	"net/http"
	"time"

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

type getDocIdRequest struct {
	DocId int `uri:"id" binding:"required"`
}

//	@BasePath	/api/v1

// ViewDocumentsUploaded godoc
//
//	@Summary		view document
//	@Schemes
//	@Description	view document
//	@Tags			document
//	@Accept			json
//	@Produce		json
//	@Success		200	{string}	Helloworld
//	@Router			/example/hello/world [get]
func (server *Server) HandlerViewDocumentsUploaded(c *gin.Context) {
	var req getDocIdRequest
	var document []DocumentView
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	doc, err := server.store.Queries.GetDocument(context.Background(), int64(req.DocId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
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
