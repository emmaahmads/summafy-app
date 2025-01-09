package api

import (
	"net/http"

	"github.com/emmaahmads/summafy/util"
	"github.com/gin-gonic/gin"
)

type getDocFilenameRequest struct {
	Filename string `uri:"filename" binding:"required"`
}

// DownloadDocument godoc
//
//	@Summary		download document
//	@Description	download document by filename
//	@Tags			document
//	@Accept			json
//	@Produce		json
//	@Param			filename	path		string	true	"filename"
//	@Success		200		{string}	string
//	@Router			/documents/{filename} [get]
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
