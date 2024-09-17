package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/emmaahmads/summafy/util"
)

// HandlerLoginPage handles GET requests to the /login endpoint
func (server *Server) HandlerLoginPage(c *gin.Context) {
	c.HTML(200, "login.html", gin.H{
		"title": "Login",
	})
}

type LoginResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	Success  bool   `json:"success"`
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
		util.MyGinLogger("Failed to get user", userInput.Username, "error:", err.Error())
		c.JSON(400, gin.H{"error": "Invalid username or password"})
		return
	}

	// Compare the provided password with the stored hashed password
	if err := util.CheckPassword(userInput.Password, user.HashedPassword); err != nil {
		util.MyGinLogger("Failed to check password", userInput.Username, "error:", err.Error())
		c.JSON(400, gin.H{"error": "Invalid username or password"})
		return
	}

	// Generate a JWT token for the user
	token, _, err := server.tokenMaker.CreateToken(user.Username, "user", time.Hour*24)
	c.Writer.Header().Set("Authorization", "Bearer "+token)
	if err != nil {
		util.MyGinLogger("Failed to generate JWT token", userInput.Username, "error:", err.Error())
		c.JSON(500, gin.H{"error": "Failed to generate JWT token"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token:    token,
		Username: user.Username,
		Success:  true,
	})

	util.MyGinLogger("user logged in", user.Username)

	// redirect to /dashboard
	c.Redirect(http.StatusFound, "/dashboard")
}
