//go:build windows

package utils

import (
	"fmt"
	"os/exec"

	"golang.org/x/sys/windows/registry"
)

func SetEnvViaPowerShell(key, value string, systemWide bool) error {
	scope := "User"
	if systemWide {
		scope = "Machine"
	}
	cmd := exec.Command("powershell", "-Command",
		fmt.Sprintf("[Environment]::SetEnvironmentVariable('%s', '%s', '%s')", key, value, scope))

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("执行 PowerShell 失败: %v, 输出: %s", err, output)
	}
	return nil
}

func GetAllUserEnvVars() (map[string]string, error) {
	k, err := registry.OpenKey(
		registry.CURRENT_USER,
		"Environment",
		registry.READ,
	)
	if err != nil {
		return nil, fmt.Errorf("无法打开注册表键: %v", err)
	}
	defer k.Close()

	names, err := k.ReadValueNames(0)
	if err != nil {
		return nil, fmt.Errorf("无法读取值名称: %v", err)
	}

	envs := make(map[string]string)
	for _, name := range names {
		value, _, err := k.GetStringValue(name)
		if err != nil {
			continue
		}
		envs[name] = value
	}

	return envs, nil
}
