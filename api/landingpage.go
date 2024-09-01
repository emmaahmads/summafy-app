package api

import (
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
	c.HTML(200, "index.html", gin.H{
		"title": "Emma Summafy - Home Page",
		"link":  link,
	})

	c.Header("Content-Type", "text/html")

}

func (server *Server) HandlerLandingPageTest(c *gin.Context) {
	c.HTML(200, "test.html", gin.H{
		"emma": "Emma",
	})

	c.Header("Content-Type", "text/html")

}
