package public

import "embed"

//go:embed all:web
var WebFiles embed.FS

//go:embed all:i18n
var I18nFiles embed.FS
