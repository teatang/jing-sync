package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ResponseSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    0,
		"data":    data,
	})
}

func ResponseError(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, gin.H{
		"success": false,
		"code":    code,
		"message": message,
	})
}
