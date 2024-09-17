package api

import "github.com/gin-gonic/gin"

func (server *Server) HandlerTest(c *gin.Context) {
	c.HTML(200, "test.html", gin.H{
		"title": "TEST TEST",
	})
}
