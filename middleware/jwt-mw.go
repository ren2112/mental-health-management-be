package middleware

import (
	"github.com/gin-gonic/gin"
	"mental-health-management-be/constants"
	"mental-health-management-be/response"
	"mental-health-management-be/utils"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 1. 获取 Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.CommonResp(c, 1, "未携带Token", nil)
			c.Abort()
			return
		}

		// 2. 解析 JWT
		claims, err := utils.ParseJWT(authHeader)
		if err != nil {
			response.CommonResp(c, 1, "Token无效或已过期", nil)
			c.Abort()
			return
		}

		// 4. 写入上下文（后续接口可直接取）
		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)

		// 继续执行
		c.Next()
	}
}

func TeacherAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		roleValue, exists := c.Get("role")
		if !exists {
			response.CommonResp(c, 1, "未登录", nil)
			c.Abort()
			return
		}

		role := roleValue.(int)

		if role != constants.RoleTeacher {
			response.CommonResp(c, 1, "仅老师可操作", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}

func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		roleValue, exists := c.Get("role")
		if !exists {
			response.CommonResp(c, 1, "未登录", nil)
			c.Abort()
			return
		}

		role := roleValue.(int)

		if role != constants.RoleAdmin {
			response.CommonResp(c, 1, "仅管理员可操作", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
