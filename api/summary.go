package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type createSummaryRequest struct {
	DocID            int64  `json:"doc_id" binding:"required"`
	SummaryGenerated string `json:"summary_generated" binding:"required"`
}

func (server *Server) HandlerCreateSummary(c *gin.Context) {
	var req createSummaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// summary, err := server.store.CreateSummary(context.Background(), db.CreateSummaryParams{
	// 	DocID:   req.DocID,
	// 	Param1:  true,
	// 	Param2:  true,
	// 	Summary: []byte(req.SummaryGenerated),
	// })

	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, err.Error())
	// 	return
	// }

	c.JSON(http.StatusOK, "createsummary endpoint")
}

func (server *Server) HandlerCreateSummarySubscription(c *gin.Context) {
	c.JSON(http.StatusOK, "getsummary endpoint")
}