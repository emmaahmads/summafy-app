package api

import (
	"fmt"
	"net/http"

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

func (server *Server) HandlerUploadDoc(c *gin.Context) {
	file, err := c.FormFile("document")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dst := "./upload/" + file.Filename
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully"})

}
