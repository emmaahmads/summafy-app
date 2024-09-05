package api

import (
	"net/http"

	"github.com/emmaahmads/summafy/util"
	"github.com/gin-gonic/gin"
)

type getDocFilenameRequest struct {
	Filename string `uri:"filename" binding:"required"`
}

func (server *Server) HandlerDownloadDoc(c *gin.Context) {
	var req getDocFilenameRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	util.MyGinLogger(req.Filename)
	_, err := server.DownloadFileFromS3(req.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.String(http.StatusOK, req.Filename)

}
