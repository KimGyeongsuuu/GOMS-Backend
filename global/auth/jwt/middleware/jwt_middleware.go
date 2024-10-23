package middleware

import (
	"GOMS-BACKEND-GO/repository"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func AuthorizeRoleJWT(secret []byte, requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token not provided"})
			c.Abort()
			return
		}

		tokenString := strings.TrimSpace(strings.Replace(authHeader, "Bearer", "", 1))

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		authority, ok := claims["authority"].(string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "No authority found in token"})
			c.Abort()
			return
		}

		if authority != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: Insufficient role"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func AccountMiddleware(accountRepo *repository.AccountRepository, secretKey []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized1"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized2"})
			c.Abort()
			return
		}
		accountID, ok := claims["accountID"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized3"})
			c.Abort()
			return
		}

		accountIDUint64 := uint64(accountID)

		account, err := accountRepo.FindByAccountID(c.Request.Context(), accountIDUint64)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized4"})
			c.Abort()
			return
		}

		c.Set("account", account)
		c.Next()
	}
}
