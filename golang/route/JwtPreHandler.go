package route

import (
	"example/entity"
	"example/utility"
	"net/http"

	"github.com/gin-gonic/gin"
)

/**
 * @author Flyinsky
 * @email w2084151024@gmail.com
 * @date 2024/12/24 11:34
 */
func CorsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 对于预检请求，直接返回状态码 204
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func preHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 前置处理：比如检查Token，记录日志等
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(200, entity.ErrorResponse[string](401, "未提供令牌"))
			c.Abort()
			return
		}
		claims, err := utility.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusOK, entity.ErrorResponse[string](401, "token不合法或已过期，请尝试重新登陆获取。"))
			c.Abort()
			return
		}
		c.Set("userId", claims.UserID)

		// 允许请求继续进行
		c.Next()
	}
}
