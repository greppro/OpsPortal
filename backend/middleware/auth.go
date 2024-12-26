package middleware

import (
	"ops-portal/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "未授权访问"})
			return
		}

		// 移除 Bearer 前缀
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 解析token
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "token无效或已过期"})
			return
		}

		// 将用户信息存储在上下文中
		c.Set("username", claims.Username)
		c.Next()
	}
}
