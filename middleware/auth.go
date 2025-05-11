package middleware

import (
	"finance-backend/controllers"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Récupérer le token du header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token d'authentification manquant"})
			c.Abort()
			return
		}

		// Vérifier le format "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Format de token invalide"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Parser et vérifier le token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return controllers.JwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token invalide"})
			c.Abort()
			return
		}

		// Extraire les claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token invalide"})
			c.Abort()
			return
		}

		// Récupérer l'ID de l'utilisateur
		userID, ok := claims["user_id"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token invalide"})
			c.Abort()
			return
		}

		// Stocker l'ID de l'utilisateur dans le contexte
		c.Set("user_id", uint(userID))
		c.Next()
	}
}
