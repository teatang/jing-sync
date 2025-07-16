package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type BaseController struct{}

// Success 成功响应
func (bc *BaseController) Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    0,
		"data":    data,
	})
}

// Error 错误响应
func (bc *BaseController) Error(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, gin.H{
		"success": false,
		"code":    code,
		"message": msg,
	})
}
