package debug

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sandwich-go/boost/xpanic"
	"golang.org/x/mod/modfile"
)

const modFileName = "go.mod"

var moduleFilePathWithSeparator string
var packagePathWithSeparator string

// MustInitAppModuleFilePath 初始化app 的mod file路径
func MustInitAppModuleFilePath() {
	md, pp, err := findModPath()
	xpanic.WhenError(err)
	if !strings.HasSuffix(md, string(os.PathSeparator)) {
		md = md + string(os.PathSeparator)
	}
	if pp != "" && !strings.HasSuffix(pp, string(os.PathSeparator)) {
		pp = pp + string(os.PathSeparator)
	}
	moduleFilePathWithSeparator = md
	packagePathWithSeparator = pp
}

// getModFile gets the module file from the root of the repo. It returns an error if the module cannot be found.
func getModFile(modFile string) (*modfile.File, error) {
	//read the file
	//nolint: gosec
	modContents, err := os.ReadFile(modFile)
	if err != nil {
		return nil, fmt.Errorf("could not read modfile: %w", err)
	}

	parsedFile, err := modfile.Parse(modFile, modContents, nil)
	if err != nil {
		return nil, fmt.Errorf("could not parse mod file:%s :%w", modFile, err)
	}

	return parsedFile, nil
}
func findModPath() (string, string, error) {
	currentPath, err := os.Getwd()
	if err != nil {
		return "", "", fmt.Errorf("could not get current path: %w", err)
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
				return "", "", errors.New("could not find go.mod file")
			}
			continue
		}
		packagePath := ""
		modFile, err := getModFile(prospectiveFile)
		if modFile != nil {
			packagePath = modFile.Module.Mod.String()
		}
		return currentPath, packagePath, err
	}
}
