package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func GetI18nMsg(msgId string, c *gin.Context) string {
	localizer := c.MustGet("localizer").(*i18n.Localizer)
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: msgId,
	})

	if err != nil {
		return ""
	}
	return msg
}
