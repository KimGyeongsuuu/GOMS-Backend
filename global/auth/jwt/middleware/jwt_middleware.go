package middleware

import (
	"GOMS-BACKEND-GO/repository"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func AccountMiddleware(accountRepo *repository.MongoAccountRepository, secretKey []byte) gin.HandlerFunc {
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
			c.JSON(http.StatusUnauthorized, gin.H{"error": "not valid access token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		accountID, ok := claims["sub"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "not found account id in claims"})
			c.Abort()
			return
		}

		accountIDObjectID, err := primitive.ObjectIDFromHex(accountID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid account id format"})
			c.Abort()
			return
		}

		account, err := accountRepo.FindByAccountID(c.Request.Context(), accountIDObjectID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized4"})
			c.Abort()
			return
		}

		c.Set("account", account)
		c.Next()
	}
}
