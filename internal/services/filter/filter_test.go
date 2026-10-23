package filter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsFileSizeAllowed(t *testing.T) {
	// 没有任何限制
	assert.True(t, IsFileSizeAllowed(500, nil, nil))

	// 只有最小限制
	min := int64(100)
	assert.True(t, IsFileSizeAllowed(100, &min, nil))
	assert.True(t, IsFileSizeAllowed(101, &min, nil))
	assert.False(t, IsFileSizeAllowed(99, &min, nil))

	// 只有最大限制
	max := int64(500)
	assert.True(t, IsFileSizeAllowed(500, nil, &max))
	assert.True(t, IsFileSizeAllowed(499, nil, &max))
	assert.False(t, IsFileSizeAllowed(501, nil, &max))

	// 同时有最小和最大限制
	assert.True(t, IsFileSizeAllowed(250, &min, &max))
	assert.False(t, IsFileSizeAllowed(50, &min, &max))
	assert.False(t, IsFileSizeAllowed(600, &min, &max))
}

func TestNormalizeVirtualPath(t *testing.T) {
	assert.Equal(t, "/a/b/c", NormalizeVirtualPath("a/b/c"))
	assert.Equal(t, "/a/b/c", NormalizeVirtualPath("/a/b/c/"))
	assert.Equal(t, "/a/b/c", NormalizeVirtualPath("\\a\\b\\c"))
	assert.Equal(t, "/a/b", NormalizeVirtualPath("/a/./b"))
	assert.Equal(t, "/a/b", NormalizeVirtualPath("/a/c/../b"))
	assert.Equal(t, "/", NormalizeVirtualPath(""))
}

func TestVirtualPathsOverlap(t *testing.T) {
	// 完全相同
	assert.True(t, VirtualPathsOverlap("/data/movies", "/data/movies"))
	assert.True(t, VirtualPathsOverlap("data/movies", "/data/movies/"))

	// 相互嵌套
	assert.True(t, VirtualPathsOverlap("/data/movies", "/data/movies/2026"))
	assert.True(t, VirtualPathsOverlap("/data/movies/2026", "/data/movies"))
	assert.True(t, VirtualPathsOverlap("/", "/data"))

	// 无关路径
	assert.False(t, VirtualPathsOverlap("/data/movies", "/data/music"))
	assert.False(t, VirtualPathsOverlap("/data/a", "/data/ab"))
}

func TestPathFilter(t *testing.T) {
	rules := `
# 忽略临时文件
*.tmp
*.log

# 忽略 node_modules
node_modules/

# 忽略特定文件夹内容
temp/*
.git/**
`
	pf := NewPathFilter(rules)

	// 应该被排除
	assert.True(t, pf.IsExcluded("test.tmp"))
	assert.True(t, pf.IsExcluded("dir/test.tmp"))
	assert.True(t, pf.IsExcluded("app.log"))
	assert.True(t, pf.IsExcluded("node_modules/express/index.js"))
	assert.True(t, pf.IsExcluded("node_modules"))
	assert.True(t, pf.IsExcluded("temp/cache.dat"))
	assert.True(t, pf.IsExcluded(".git/objects/abc"))

	// 不应该被排除
	assert.False(t, pf.IsExcluded("test.txt"))
	assert.False(t, pf.IsExcluded("dir/app.go"))
	assert.False(t, pf.IsExcluded("my_node_modules/index.js"))
	assert.False(t, pf.IsExcluded("template/index.html"))
}
