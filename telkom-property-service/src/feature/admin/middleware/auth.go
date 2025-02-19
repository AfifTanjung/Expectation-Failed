package middleware

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)


func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: No token provided"})
			c.Abort()
			return
		}


		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Invalid token format"})
			c.Abort()
			return
		}

		
		tokenString := tokenParts[1]
		payload, err := decodeJWTPayload(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Invalid token payload"})
			c.Abort()
			return
		}

		
		
		
		// var admin model.Admin
		// if err := db.Where("token = ?", tokenString).First(&admin).Error; err != nil {
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized invalid token"})
		// 	c.Abort()
		// 	return 
		// }

		
		
		c.Set("user", payload)
		
		c.Next()
	
	}
}

func decodeJWTPayload(token string) (map[string]interface{}, error) {
	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		return nil, http.ErrNoCookie
	}


	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return nil, err
	}

	return payload, nil
}


