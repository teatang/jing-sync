package i18n

import (
	"testing"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestI18nBundleInitialization(t *testing.T) {
	bundle := GetI18nBundle()
	require.NotNil(t, bundle)

	// 测试中文翻译
	zhLocalizer := i18n.NewLocalizer(bundle, "zh-CN")
	zhMsg, err := zhLocalizer.Localize(&i18n.LocalizeConfig{
		MessageID: "error_msg.token_check_failed",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, zhMsg)

	// 测试英文翻译
	enLocalizer := i18n.NewLocalizer(bundle, "en")
	enMsg, err := enLocalizer.Localize(&i18n.LocalizeConfig{
		MessageID: "error_msg.token_check_failed",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, enMsg)

	// 验证中英文内容不同
	assert.NotEqual(t, zhMsg, enMsg)
}
