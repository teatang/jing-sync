package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatBytes(t *testing.T) {
	assert.Equal(t, "0 B", FormatBytes(0))
	assert.Equal(t, "500 B", FormatBytes(500))
	assert.Equal(t, "1.00 KB", FormatBytes(1024))
	assert.Equal(t, "1.50 MB", FormatBytes(1572864))
	assert.Equal(t, "2.00 GB", FormatBytes(2147483648))
	assert.Equal(t, "1.00 TB", FormatBytes(1099511627776))
}

func TestSecureRandString(t *testing.T) {
	str1, err := SecureRandString(16)
	require.NoError(t, err)
	assert.Equal(t, 16, len(str1))

	str2, err := SecureRandString(16)
	require.NoError(t, err)
	assert.NotEqual(t, str1, str2, "两次随机生成的字符串不应该相同")
}

func TestSHA256(t *testing.T) {
	hash := SHA256("hello world")
	assert.Equal(t, "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9", hash)
}

func TestFileUtil(t *testing.T) {
	tmpDir := t.TempDir()
	testDir := tmpDir + "/test_sub_dir"

	assert.False(t, FileExists(testDir))
	err := EnsureDir(testDir)
	require.NoError(t, err)
	assert.True(t, FileExists(testDir))

	assert.Equal(t, "123", UintToString(123))
}

func TestReadOrSet(t *testing.T) {
	tmpFile := t.TempDir() + "/secret.test"
	defer os.Remove(tmpFile)

	val1, err := ReadOrSet(tmpFile, "default_val", false)
	require.NoError(t, err)
	assert.Equal(t, "default_val", val1)

	// 不覆盖读取
	val2, err := ReadOrSet(tmpFile, "another_val", false)
	require.NoError(t, err)
	assert.Equal(t, "default_val", val2)

	// 强制覆盖
	val3, err := ReadOrSet(tmpFile, "forced_val", true)
	require.NoError(t, err)
	assert.Equal(t, "forced_val", val3)
}
