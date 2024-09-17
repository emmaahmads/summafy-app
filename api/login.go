package api

import (
	"github.com/gin-gonic/gin"

	"github.com/emmaahmads/summafy/util"
)

// HandlerLoginPage handles GET requests to the /login endpoint
func (server *Server) HandlerLoginPage(c *gin.Context) {
	c.HTML(200, "login.html", gin.H{
		"title": "Login",
	})
}

// HandlerLogin handles POST requests to the /login endpoint
func (server *Server) HandlerLogin(c *gin.Context) {
	var userInput struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Verify the user credentials
	user, err := server.store.GetUser(c, userInput.Username)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid username or password"})
		return
	}

	// Compare the provided password with the stored hashed password
	if err := util.CheckPassword(userInput.Password, user.HashedPassword); err != nil {
		c.JSON(400, gin.H{"error": "Invalid username or password"})
		return
	}

	/* // Generate a JWT token for the user
	token, err := util.GenerateJWTToken(user.Username)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate JWT token"})
		return
	}
	*/
	c.JSON(200, gin.H{"token": "N/A", "username": user.Username})
}
