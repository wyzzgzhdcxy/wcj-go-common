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
	// 打开当前用户的环境变量注册表键
	k, err := registry.OpenKey(
		registry.CURRENT_USER,
		"Environment",
		registry.READ,
	)
	if err != nil {
		return nil, fmt.Errorf("无法打开注册表键: %v", err)
	}
	defer k.Close()

	// 获取所有值名称
	names, err := k.ReadValueNames(0)
	if err != nil {
		return nil, fmt.Errorf("无法读取值名称: %v", err)
	}

	// 读取所有环境变量
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
