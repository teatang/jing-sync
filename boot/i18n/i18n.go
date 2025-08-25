package i18n

import (
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var i18nBundle *i18n.Bundle

func GetI18nBundle() *i18n.Bundle {
	if i18nBundle == nil {
		InitI18n()
	}
	return i18nBundle

}

func InitI18n() {
	i18nBundle = i18n.NewBundle(language.Chinese)
	i18nBundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	i18nBundle.MustLoadMessageFile("locales/active.en.toml")
	i18nBundle.MustLoadMessageFile("locales/active.zh-CN.toml")
}
