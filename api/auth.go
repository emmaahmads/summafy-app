package api

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtKey = []byte(os.Getenv("SECRET_KEY"))

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func generateJWT(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func (server *Server) middlewareAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username, err := GetApi(ctx.Request.Header)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no authorization header"})
			return
		}
		user, err := server.store.GetUser(ctx, username)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "not able to retrive user"})
			return
		}
		ctx.Set("user", user)
		ctx.Next()
	}
}

func GetApi(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("malformed authorization header")

	}
	splitAuth := strings.Split(authHeader, " ")

	return splitAuth[1], nil
}
