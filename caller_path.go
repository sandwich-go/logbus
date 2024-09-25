package logbus

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

var pathMap sync.Map

func getCallerPath(fullPath string) string {
	if v, ok := pathMap.Load(fullPath); ok {
		return v.(string)
	}

	// 获取从module到当前文件的相对路径
	relPath, err := filepath.Rel(modulePath, fullPath)
	if err != nil {
		return ""
	}

	pathMap.Store(fullPath, relPath)

	return relPath
}

// findModuleRoot 从给定的目录开始向上查找，直到找到包含go.mod文件的目录
func findModuleRoot(dir string) (string, error) {
	for {
		// 检查当前目录是否包含go.mod文件
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return dir, nil
		}
		// 获取当前目录的父目录
		parentDir := filepath.Dir(dir)
		// 如果父目录与当前目录相同，则我们已经到达了文件系统的顶部
		if parentDir == dir {
			break
		}
		dir = parentDir
	}
	// 如果没有找到go.mod文件，返回错误
	return "", fmt.Errorf("go.mod file not found")
}
