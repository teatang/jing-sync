package i18n

import (
	"jing-sync/boot/logger"
	"jing-sync/public"

	"fmt"
	"io/fs"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pelletier/go-toml/v2"
	"golang.org/x/text/language"
)

var i18nBundle *i18n.Bundle

func GetI18nBundle() *i18n.Bundle {
	if i18nBundle == nil {
		I18nInit()
	}
	return i18nBundle
}

func I18nInit() {
	i18nBundle = i18n.NewBundle(language.Chinese)

	// 告诉 go-i18n 如何解析 "toml" 格式的文件
	i18nBundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	folderPath := "i18n"
	files, err := fs.ReadDir(public.I18nFiles, folderPath)
	if err != nil {
		logger.GetLogger().Fatalf("i18n: Failed to read embedded locales directory: %v", err)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filePath := fmt.Sprintf("%s/%s", folderPath, file.Name())
		fileContent, err := public.I18nFiles.ReadFile(filePath)
		if err != nil {
			logger.GetLogger().Panicf("i18n: Failed to read embedded file %s: %v", filePath, err)
		}
		// go-i18n 会根据文件名的后缀 (如 .toml) 来选择已注册的 unmarshaler
		_, err = i18nBundle.ParseMessageFileBytes(fileContent, file.Name())
		if err != nil {
			// 打印更详细的错误信息，包括文件名和原始错误
			logger.GetLogger().Panicf("i18n: Failed to parse message file %s: %v", file.Name(), err)
		}
		logger.GetLogger().Infof("i18n: Successfully loaded messages from embedded file: %s\n", file.Name())
	}
}
