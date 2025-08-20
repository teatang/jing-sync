package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type PageInfo struct {
	Page  int   `json:"page"`
	Size  int   `json:"size"`
	Total int64 `json:"total"`
}

type PageList[T any] struct {
	List       []T      `json:"list"`
	Pagination PageInfo `json:"pagination"`
}

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
		"msg": message,
	})
}
