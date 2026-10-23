package filter

import (
	"path"
	"path/filepath"
	"strings"
)

// IsFileSizeAllowed 检查文件大小是否在允许区间内
// minSize: 允许的最小字节数（含），为 nil 表示不限制
// maxSize: 允许的最大字节数（含），为 nil 表示不限制
func IsFileSizeAllowed(size int64, minSize *int64, maxSize *int64) bool {
	if minSize != nil && size < *minSize {
		return false
	}
	if maxSize != nil && size > *maxSize {
		return false
	}
	return true
}

// NormalizeVirtualPath 规范化虚拟路径
// 统一使用正斜杠，保证以 / 开头，去除结尾的多余斜杠（根目录 / 除外）
func NormalizeVirtualPath(p string) string {
	cleaned := strings.ReplaceAll(p, "\\", "/")
	if !strings.HasPrefix(cleaned, "/") {
		cleaned = "/" + cleaned
	}
	cleaned = path.Clean(cleaned)
	return cleaned
}

// VirtualPathsOverlap 检查两个虚拟路径是否存在重叠（相同或相互嵌套）
// 用于在同步前做安全检查，防止源路径与目标路径相互包含导致死循环
func VirtualPathsOverlap(firstPath, secondPath string) bool {
	first := NormalizeVirtualPath(firstPath)
	second := NormalizeVirtualPath(secondPath)

	if first == second {
		return true
	}

	// 根路径与其他任何路径都重叠
	if first == "/" || second == "/" {
		return true
	}

	// 检查第一条路径是否是第二条路径的祖先
	if strings.HasPrefix(second, strings.TrimSuffix(first, "/")+"/") {
		return true
	}

	// 检查第二条路径是否是第一条路径的祖先
	if strings.HasPrefix(first, strings.TrimSuffix(second, "/")+"/") {
		return true
	}

	return false
}

// PathFilter 路径通配符与排除规则过滤器
type PathFilter struct {
	patterns []string
}

// NewPathFilter 根据换行分隔的规则文本创建过滤器
func NewPathFilter(rules string) *PathFilter {
	var patterns []string
	lines := strings.Split(rules, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		patterns = append(patterns, trimmed)
	}
	return &PathFilter{patterns: patterns}
}

// IsExcluded 检查指定相对或绝对文件路径是否匹配排除规则
func (pf *PathFilter) IsExcluded(filePath string) bool {
	cleanPath := strings.TrimPrefix(strings.ReplaceAll(filePath, "\\", "/"), "/")
	baseName := path.Base(cleanPath)

	for _, pattern := range pf.patterns {
		pattern = strings.TrimPrefix(pattern, "/")

		// 目录排除，如 node_modules/
		if strings.HasSuffix(pattern, "/") {
			dirPattern := strings.TrimSuffix(pattern, "/")
			if cleanPath == dirPattern || strings.HasPrefix(cleanPath, dirPattern+"/") {
				return true
			}
			// 支持多级中的目录名匹配，例如 a/b/node_modules/c
			parts := strings.Split(cleanPath, "/")
			for _, part := range parts {
				if part == dirPattern {
					return true
				}
			}
		}

		// ** 递归通配符，如 .git/** 或 **/tmp/*
		if strings.Contains(pattern, "**") {
			prefix := strings.Split(pattern, "/**")[0]
			if cleanPath == prefix || strings.HasPrefix(cleanPath, prefix+"/") {
				return true
			}
		}

		// 单级通配符直接匹配文件名，如 *.tmp, *.log
		if matched, _ := filepath.Match(pattern, baseName); matched {
			return true
		}

		// 全路径匹配，如 temp/*
		if matched, _ := filepath.Match(pattern, cleanPath); matched {
			return true
		}
	}

	return false
}
