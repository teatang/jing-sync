package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetDefaultConfig(t *testing.T) {
	SetDefaultConfig(EnvTypeDev)
	require.NotNil(t, Cfg)
	assert.Equal(t, "jing-sync", Cfg.SiteName)
	assert.Equal(t, "jing-sync.db", Cfg.DbName)
	assert.Equal(t, 8888, Cfg.Port)
	assert.Equal(t, 60, Cfg.Timeout)
	assert.Equal(t, EnvTypeDev, Cfg.Env)
}

func TestInitConfig(t *testing.T) {
	// 备份现有配置状态
	oldCfg := Cfg
	t.Cleanup(func() {
		Cfg = oldCfg
		_ = os.RemoveAll("data")
	})

	// 准备临时测试配置目录
	err := os.MkdirAll("data", 0755)
	require.NoError(t, err)

	testConfigFile := filepath.Join("data", "config_test.json")
	configJSON := `{
		"site_name": "jing-sync-test",
		"db_name": ":memory:",
		"port": 9999,
		"timeout": 30
	}`
	err = os.WriteFile(testConfigFile, []byte(configJSON), 0644)
	require.NoError(t, err)

	Cfg = nil
	err = InitConfig(EnvTypeTest)
	require.NoError(t, err)
	require.NotNil(t, Cfg)
	assert.Equal(t, "jing-sync-test", Cfg.SiteName)
	assert.Equal(t, ":memory:", Cfg.DbName)
	assert.Equal(t, 9999, Cfg.Port)
	assert.Equal(t, 30, Cfg.Timeout)

	// 重复调用 InitConfig 应当直接返回 nil (已初始化状态)
	err = InitConfig(EnvTypeTest)
	require.NoError(t, err)
}
