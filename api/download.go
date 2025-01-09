package api

import (
	"net/http"

	"github.com/emmaahmads/summafy/util"
	"github.com/gin-gonic/gin"
)

type getDocFilenameRequest struct {
	Filename string `uri:"filename" binding:"required"`
}

//	@BasePath	/api/v1

// PingExample godoc
//
//	@Summary	ping example hello
//	@Schemes
//	@Description	do ping
//	@Tags			example
//	@Accept			json
//	@Produce		json
//	@Success		200	{string}	Helloworld
//	@Router			/example/helloworld [get]
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
