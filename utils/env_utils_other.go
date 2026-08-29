//go:build !windows

package utils

import (
	"fmt"
	"runtime"
)

func SetEnvViaPowerShell(key, value string, systemWide bool) error {
	return fmt.Errorf("SetEnvViaPowerShell 不支持当前平台: %s", runtime.GOOS)
}

func GetAllUserEnvVars() (map[string]string, error) {
	return nil, fmt.Errorf("GetAllUserEnvVars 不支持当前平台: %s", runtime.GOOS)
}
