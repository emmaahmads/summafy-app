package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/emmaahmads/summafy/util"
)

//	@BasePath	/api/v1

// Login godoc
//
//	@Summary		login
//	@Description	login
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{string}	Helloworld
//	@Router			/example/helloworld [get]
func (server *Server) HandlerLogin(c *gin.Context) {
	var userInput struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&userInput); err != nil {
		util.MyGinLogger(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	// Verify the user credentials
	user, err := server.store.GetUser(c, userInput.Username)
	if err != nil {
		util.MyGinLogger(err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}
	util.MyGinLogger(userInput.Password, user.HashedPassword)
	// Compare the provided password with the stored hashed password
	if err := util.CheckPassword(userInput.Password, user.HashedPassword); err != nil {
		util.MyGinLogger(err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	token, err := server.generateJWT(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	// Set JWT as HttpOnly cookie
	c.SetCookie("session_token", token, 86400, "/", "", false, true)
	c.JSON(200, gin.H{"success": true, "username": user.Username})
}
