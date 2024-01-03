package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	connections "github.com/AkbarFikri/signconnect_backend/Connections"
	models "github.com/AkbarFikri/signconnect_backend/Models"

)

func AuthJWTToken(c *gin.Context) {
	// GET the authorization cookies
	tokenString, err := c.Cookie("X-Authorization")

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	// validate the JWT Token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to load cookies",
		})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// check exp
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		// Find user with the token
		var user models.User
		connections.DB.First(&user, claims["sub"])

		if user.ID == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		// add the user to req
		c.Set("user", user)

		// next
		c.Next()
	} else {
		fmt.Println(err)
	}
}

func AuthAPIKey(c *gin.Context) {
	// get API Key Header
	apiKey := c.GetHeader("A-Authorization")
	fmt.Print(apiKey)

	if apiKey != os.Getenv("API_KEY") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Your not allowed to access this website",
		})
	}

	c.Next()
}
