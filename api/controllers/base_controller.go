package controllers

import (
	"github.com/gin-gonic/gin"
	"jing-sync/internal/utils"
)

type BaseController struct{}

// Success 成功响应
func (bc *BaseController) Success(c *gin.Context, data interface{}) {
	utils.ResponseSuccess(c, data)
}

// Error 错误响应
func (bc *BaseController) Error(c *gin.Context, code int, msg string) {
	utils.ResponseError(c, code, msg)
}
