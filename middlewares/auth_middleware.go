package middlewares

import (
	// "net/http"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWT声明结构
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// JWT认证中间件
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// tokenString := c.GetHeader("Authorization")
		// if tokenString == "" {
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证Token"})
		// 	c.Abort()
		// 	return
		// }

		// claims := &Claims{}
		// token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// 	return jwtSecret, nil
		// })

		// if err != nil || !token.Valid {
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "无效Token"})
		// 	c.Abort()
		// 	return
		// }

		// c.Set("claims", claims)
		c.Next()
	}
}