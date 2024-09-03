package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Link struct {
	Href string
	Rel  string
}

func (server *Server) HandlerLandingPage(c *gin.Context) {
	link := Link{
		Href: "https://example.com/link",
		Rel:  "related",
	}
	c.HTML(200, "dashboard.html", gin.H{
		"title": "Emma Summafy - Home Page",
		"link":  link,
	})

	c.Header("Content-Type", "text/html")

}

func (server *Server) HandlerUploadPage(c *gin.Context) {
	c.HTML(200, "uploadform.html", gin.H{
		"emma": "Emma",
	})

	c.Header("Content-Type", "text/html")

	file, err := c.FormFile("document")
	if err != nil {
		return
	}

	gin.Logger()
	fmt.Println(file.Filename)

}

func (server *Server) HandlerViewPage(c *gin.Context) {
	c.HTML(200, "view.html", gin.H{
		"emma": "Emma",
	})

	c.Header("Content-Type", "text/html")

}
