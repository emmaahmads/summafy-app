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
//	@Success		201		{object}	signupOutput
//	@Router			/signup [post]
func (server *Server) HandlerSignup(c *gin.Context) {
	var userInput struct {
		Username        string `json:"username" binding:"required"`
		Email           string `json:"email" binding:"required,email"`
		Password        string `json:"password" binding:"required,min=4"`
		ConfirmPassword string `json:"confirm_password" binding:"required,min=4"`
	}

	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		util.MyGinLogger(err.Error())
		return
	}

	if userInput.Password != userInput.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords do not match"})
		util.MyGinLogger("Passwords do not match")
		return
	}

	hashedPassword, err := util.HashPassword(userInput.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	token, err := server.generateJWT(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	// Set JWT as HttpOnly cookie
	c.SetCookie("session_token", token, 86400, "/", "", !util.IsDevelopment, true)
	c.JSON(http.StatusCreated, gin.H{"status": true, "user": user})
}
