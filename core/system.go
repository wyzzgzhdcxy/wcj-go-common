package core

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func IsWin() bool {
	return runtime.GOOS == "windows"
}

func GetTempDir() string {
	dir, _ := os.UserCacheDir()
	return dir + "/wtools"
}

func GetUserHomeDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("获取用户主目录失败: %v\n", err)
		return ""
	}
	return homeDir
}

// GetHostsPath 返回当前系统的 hosts 文件路径
func GetHostsPath() string {
	if runtime.GOOS == "windows" {
		return getWindowsHostsPath()
	}
	return getUnixHostsPath()
}

// getWindowsHostsPath 获取 Windows 系统的 hosts 文件路径
func getWindowsHostsPath() string {
	// 尝试从 SystemRoot 环境变量获取路径
	systemRoot := os.Getenv("SystemRoot")
	if systemRoot != "" {
		return filepath.Join(systemRoot, "System32", "drivers", "etc", "hosts")
	}

	// 回退方案：尝试常见路径
	possiblePaths := []string{
		`C:\Windows\System32\drivers\etc\hosts`, // Win7/10/11
		`C:\WINNT\System32\drivers\etc\hosts`,   // WinXP/2000
		`D:\Windows\System32\drivers\etc\hosts`, // 可能安装在D盘
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	// 最终回退
	return `C:\Windows\System32\drivers\etc\hosts`
}

// getUnixHostsPath 获取 Unix-like 系统的 hosts 文件路径
func getUnixHostsPath() string {
	// 标准 Unix/Linux 路径
	unixPaths := []string{
		"/etc/hosts",         // 大多数Linux/Unix
		"/private/etc/hosts", // macOS
	}

	for _, path := range unixPaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	// 最终回退
	return "/etc/hosts"
}
