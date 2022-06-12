package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/YJ9938/DouYin/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type Claims struct {
	jwt.StandardClaims
}

var secretKey = []byte(config.C.JWT.SecretKey)

// Query form's auth token rather than the query string
func FormAuthMiddleware(c *gin.Context) {
	auth(c, c.PostForm("token"))
}

// Auth middleware used to handle authentication for every route that needs auth
func QueryAuthMiddleware(c *gin.Context) {
	// Check the token parameter
	auth(c, c.Query("token"))
}

func auth(c *gin.Context, token string) {
	if token == "" {
		c.AbortWithStatusJSON(http.StatusOK, Response{
			StatusCode: http.StatusBadRequest,
			StatusMsg:  "鉴权参数错误",
		})
		return
	}

	// Check the token is valid and store user ID to the context
	claims := parseToken(token)
	if claims == nil {
		c.AbortWithStatusJSON(http.StatusOK, Response{
			StatusCode: http.StatusBadRequest,
			StatusMsg:  "用户鉴权错误",
		})
		return
	}
	log.Printf("Token user ID: %s\n", claims.Id)

	// If auth success, we pass an 'id' to gin's context
	id, _ := strconv.ParseInt(claims.Id, 10, 64)
	c.Set("id", id)
}

// signJWT signs a JWT and returns it.
func signJWT(userID int64) (string, error) {
	expireTime := time.Now().Add(time.Duration(config.C.JWT.ExpireMinutes) * time.Minute)
	claims := Claims{
		StandardClaims: jwt.StandardClaims{
			Id:        fmt.Sprintf("%d", userID),
			ExpiresAt: expireTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	return token.SignedString(secretKey)
}

// parseToken verify a JWT string and returns its claims.
func parseToken(tokenString string) *Claims {
	var claims Claims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		return nil
	}
	return token.Claims.(*Claims)
}
