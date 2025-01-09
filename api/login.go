package api

import (
	"net/http"

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
		util.MyGinLogger(err.Error())
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	// Verify the user credentials
	user, err := server.store.GetUser(c, userInput.Username)
	if err != nil {
		util.MyGinLogger(err.Error())
		c.JSON(402, gin.H{"error": "Invalid username or password"})
		return
	}
	util.MyGinLogger(userInput.Password, user.HashedPassword)
	// Compare the provided password with the stored hashed password
	if err := util.CheckPassword(userInput.Password, user.HashedPassword); err != nil {
		util.MyGinLogger(err.Error())
		c.JSON(403, gin.H{"error": "Invalid username or password"})
		return
	}

	token, err := server.generateJWT(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(200, gin.H{"success": true, "token": token, "username": user.Username})
}
