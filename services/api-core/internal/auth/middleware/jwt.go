package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var secret = []byte("supersecret")

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "authorization header missing",
			})
			return
		}

		// ✅ Proper Bearer format check
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "invalid authorization format",
			})
			return
		}

		tokenString := parts[1]

		// ✅ Parse with validation
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {

			// 🔒 Check signing method
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}

			return secret, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "invalid or expired token",
			})
			return
		}

		// ✅ Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "invalid token claims",
			})
			return
		}

		// ✅ Get user_id safely
		userID, ok := claims["user_id"].(string)
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "invalid user_id in token",
			})
			return
		}

		// ✅ Store in context
		c.Set("user_id", userID)

		c.Next()
	}
}