package logger

import (
	"jing-sync/boot/config"

	"github.com/gin-gonic/gin"
)

type WebErrorObj struct {
	Type      string `json:"type"`
	Method    string `json:"method"`
	Path      string `json:"path"`
	GetParams string `json:"get_params"`
	ErrMsg    string `json:"err_msg"`
}

func WebError(c *gin.Context, err_type config.ErrType) {
	e := WebErrorObj{
		Type:      config.ErrTypeName[err_type],
		Method:    c.Request.Method,
		Path:      c.Request.URL.Path,
		GetParams: c.Request.URL.RawQuery,
		ErrMsg:    config.ErrTypeMsg[err_type],
	}
	GetLogger().Error(e)
}
