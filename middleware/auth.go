package middleware

import (
	"go-gin-api/models"
	"go-gin-api/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing or invalid."})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := utils.ParseToken(tokenStr, false)
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token."})
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Set("userID", claims["user_id"])
		c.Set("email", claims["email"])
		c.Next()
	}
}

func PermissionMiddleware(db *gorm.DB, requiredPermission string) gin.HandlerFunc {
	return func (c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		var user models.User
		if err := db.Preload("Role.Permissions").First(&user, userID).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}

		hasPermission := false
		for _, permission := range user.Role.Permissions {
			if permission.Name == requiredPermission {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			return
		}

		c.Set("userRole", user.Role)
		c.Set("userPermissions", user.Role.Permissions)

		c.Next()
	}
}