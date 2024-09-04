package api

import (
	"github.com/gin-gonic/gin"
)

func (server *Server) HandlerViewPage(c *gin.Context) {
	c.HTML(200, "view.html", gin.H{
		"emma": "Emma",
	})

	c.Header("Content-Type", "text/html")

}
