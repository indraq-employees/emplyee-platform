package middleware

import (
	"net/http"

	"employee-platform/backend/internal/utils"
	"github.com/gin-gonic/gin"
)

func RequireRoles(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		for _, allowed := range allowedRoles {
			if role == allowed {
				c.Next()
				return
			}
		}

		utils.Error(c, http.StatusForbidden, "forbidden")
		c.Abort()
	}
}