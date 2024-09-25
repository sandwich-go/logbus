package debug

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sandwich-go/boost/xpanic"
)

const modFileName = "go.mod"

var moduleFilePath string

// MustInitAppModuleFilePath 初始化app 的mod file路径
func MustInitAppModuleFilePath() {
	md, err := findModPath()
	xpanic.WhenError(err)
	if !strings.HasSuffix(md, string(os.PathSeparator)) {
		md = md + string(os.PathSeparator)
	}
	moduleFilePath = md
}

func findModPath() (string, error) {
	currentPath, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("could not get current path: %w", err)
	}
	for {
		exists := true
		prospectiveFile := filepath.Join(currentPath, modFileName)
		if _, err := os.Stat(prospectiveFile); os.IsNotExist(err) {
			exists = false
		}
		if !exists {
			lastPath := currentPath
			currentPath = filepath.Dir(currentPath)
			if lastPath == currentPath {
				return "", errors.New("could not find go.mod file")
			}
			continue
		}
		return currentPath, nil
	}
}
