package api

import (
	"net/http"

	db "github.com/emmaahmads/summafy/db/sqlc"
	"github.com/emmaahmads/summafy/util"
	"github.com/gin-gonic/gin"
)

//	@BasePath	/api/v1

// PingExample godoc
//
//	@Summary	signup
//	@Description	signup
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		signupInput	true	"signup input"
//	@Success		200		{object}	signupOutput
//	@Router			/signup [post]
func (server *Server) HandlerSignup(c *gin.Context) {
	var userInput struct {
		Username        string `json:"username" binding:"required"`
		Email           string `json:"email" binding:"required,email"`
		Password        string `json:"password" binding:"required,min=4"`
		ConfirmPassword string `json:"confirm_password" binding:"required,min=4"`
	}

	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		util.MyGinLogger(err.Error())
		return
	}

	if userInput.Password != userInput.ConfirmPassword {
		c.JSON(402, gin.H{"error": "Passwords do not match"})
		util.MyGinLogger("Passwords do not match")
		return
	}

	hashedPassword, err := util.HashPassword(userInput.Password)
	if err != nil {
		c.JSON(501, gin.H{"error": "Failed to hash password"})
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
		c.JSON(500, gin.H{"error": err})
		return
	}
	token, err := server.generateJWT(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(201, gin.H{"status": true, "user": user, "token": token})
}
