package api

import (
	db "github.com/emmaahmads/summafy/db/sqlc"
	"github.com/emmaahmads/summafy/util"
	"github.com/gin-gonic/gin"
)

// HandlerSignupPage handles GET requests to the /signup endpoint
func (server *Server) HandlerSignupPage(c *gin.Context) {
	c.HTML(200, "signup.html", gin.H{
		"title": "Sign Up",
	})
}

// HandlerSignup handles POST requests to the /signup endpoint
func (server *Server) HandlerSignup(c *gin.Context) {
	var userInput struct {
		Username        string `json:"username" binding:"required"`
		Email           string `json:"email" binding:"required,email"`
		Password        string `json:"password" binding:"required,min=4"`
		ConfirmPassword string `json:"confirm-password" binding:"required,min=4"`
	}

	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if userInput.Password != userInput.ConfirmPassword {
		c.JSON(400, gin.H{"error": "Passwords do not match"})
		return
	}

	hashedPassword, err := util.HashPassword(userInput.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to hash password"})
		return
	}

	arg := db.CreateUserParams{
		Username:       userInput.Username,
		HashedPassword: hashedPassword,
		FullName:       util.RandomString(int(util.RandomInt(5, 10))),
		Email:          userInput.Email,
	}

	user, err := server.store.CreateUser(c, arg)
	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_username_key\"" {
			c.JSON(400, gin.H{"error": "Username already exists"})
			return
		}
		c.JSON(500, gin.H{"error": "Failed to create user"})
		return
	}

	c.Set("username", userInput.Username)

	c.JSON(201, gin.H{"status": true, "user": user})
}
