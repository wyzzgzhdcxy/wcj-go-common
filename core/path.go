package core

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

// ExecPath 获取执行命令所在路径
func ExecPath() string {
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Printf("The ExecPath failed: %s\n\n", err.Error())
	}
	return path
}

func GetAppDir() string {
	path, _ := os.Executable()
	dir := filepath.Dir(path)
	macAppPath := "/proxy.app/Contents/MacOS"
	if strings.HasSuffix(dir, macAppPath) {
		dir = strings.ReplaceAll(dir, macAppPath, "")
	}
	return dir
}
