package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aeof/douyin/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type Claims struct {
	jwt.StandardClaims
}

var secretKey = []byte(config.C.JWT.SecretKey)

func getUserID(c *gin.Context) int64 {
	// token's user ID
	idVal, exists := c.Get("id")
	if !exists {
		log.Panicf("id %v not found\n", idVal)
	}
	var id int64
	fmt.Sscanf(idVal.(string), "%d", &id)

	return id
}

func AuthMiddleware(c *gin.Context) {
	// Check the token parameter
	token := c.Query("token")
	if token == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, Status{
			StatusCode:    http.StatusUnauthorized,
			StatusMessage: "用户认证错误",
		})
		return
	}

	// Check the token is valid and store user ID to the context
	claims := parseToken(token)
	if claims == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, Status{
			StatusCode:    http.StatusUnauthorized,
			StatusMessage: "用户认证错误",
		})
		return
	}
	fmt.Printf("%#v\n", claims)
	c.Set("id", claims.Id)
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
