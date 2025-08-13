package middlewares

import (
	"jing-sync/utils"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWT认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			utils.ResponseError(c, http.StatusUnauthorized, "未提供认证Token")
			c.Abort()
			return
		}

		claims := &utils.Claims{}
		jwtSecret, _ := utils.GetSecretKey()
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			utils.ResponseError(c, http.StatusUnauthorized, "无效Token")
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}