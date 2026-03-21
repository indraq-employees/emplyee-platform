package middleware

import (
	"net/http"
	"strings"

	"employee-platform/backend/internal/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			utils.Error(c, http.StatusUnauthorized, "missing or invalid authorization header")
			c.Abort()
			return
		}

		claims, err := utils.ParseJWT(parts[1], jwtSecret)
		if err != nil {
			utils.Error(c, http.StatusUnauthorized, "invalid token")
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}