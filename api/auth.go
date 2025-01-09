package api

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"github.com/emmaahmads/summafy/util"
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func (server *Server) generateJWT(username string) (string, error) {
	util.MyGinLogger("Starting JWT generation for username:", username)

	var jwtKey = []byte(server.secretKey)
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		util.MyGinLogger("Error generating JWT:", err.Error())
		return "", err
	}

	util.MyGinLogger("JWT generation successful for username:", username)
	util.MyGinLogger("JWT:", tokenString)
	return tokenString, nil
}

func (server *Server) middlewareAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username, err := server.GetUser(ctx.Request.Header)
		if err != nil {
			util.MyGinLogger(err.Error())
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no authorization header"})
			return
		}
		util.MyGinLogger("Extracted username from JWT:", username)
		user, err := server.store.GetUser(ctx, username)
		if err != nil {
			util.MyGinLogger(err.Error())
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "not able to retrive user"})
			return
		}
		util.MyGinLogger("User retrieved:", user.Username)
		ctx.Set("username", user.Username)
		ctx.Next()
	}
}

func (server *Server) GetUser(headers http.Header) (string, error) {
	var jwtKey = []byte(server.secretKey)
	util.MyGinLogger("Retrieving Authorization header")
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		util.MyGinLogger("Authorization header is missing")
		return "", errors.New("malformed authorization header")
	}

	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) != 2 {
		return "", errors.New("malformed authorization header")
	}

	tokenStr := splitAuth[1]
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "", errors.New("invalid token signature")
		}
		return "", errors.New("could not parse token")
	}

	if !token.Valid {
		return "", errors.New("invalid token")
	}

	return claims.Username, nil
}
