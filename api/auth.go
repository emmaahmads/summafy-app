package api

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func (server *Server) generateJWT(username string) (string, error) {
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
		return "", err
	}

	return tokenString, nil
}

func (server *Server) middlewareAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, err := server.GetClaimsFromCookie(ctx)
		if err != nil || claims.ExpiresAt < time.Now().Unix() {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		user, err := server.store.GetUser(ctx, claims.Username)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "not able to retrieve user"})
			return
		}
		ctx.Set("username", user.Username)
		ctx.Next()
	}
}

// Extract JWT from HttpOnly cookie
func (server *Server) GetClaimsFromCookie(ctx *gin.Context) (*Claims, error) {
	var jwtKey = []byte(server.secretKey)
	claims := &Claims{}
	tokenStr, err := ctx.Cookie("session_token")
	if err != nil {
		return nil, errors.New("missing session cookie")
	}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, errors.New("invalid token signature")
		}
		return nil, errors.New("could not parse token")
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

// Deprecated: for backward compatibility with header-based auth
func (server *Server) GetClaimsFromJWT(headers http.Header) (*Claims, error) {
	var jwtKey = []byte(server.secretKey)
	claims := &Claims{}
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return nil, errors.New("malformed authorization header")
	}

	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) != 2 {
		return nil, errors.New("malformed authorization header")
	}

	tokenStr := splitAuth[1]

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, errors.New("invalid token signature")
		}
		return nil, errors.New("could not parse token")
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
