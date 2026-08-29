package core

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"log"

	"github.com/gogf/gf/os/gproc"
	"github.com/gogf/gf/text/gstr"
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
		return GetWindowsHostsPath()
	}
	return GetUnixHostsPath()
}

// getWindowsHostsPath 获取 Windows 系统的 hosts 文件路径
func GetWindowsHostsPath() string {
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
func GetUnixHostsPath() string {
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

// UpdateSystemDate 通过系统命令更新本机的系统日期/时间。
// 注意：调用方需要具备管理员/root 权限，否则将返回 false。
func UpdateSystemDate(dateTime string) bool {
	_, err1 := gproc.ShellExec(`date  ` + gstr.Split(dateTime, " ")[0])
	_, err2 := gproc.ShellExec(`time  ` + gstr.Split(dateTime, " ")[1])
	if err1 != nil && err2 != nil {
		log.Printf("更新系统时间错误:请用管理员身份启动程序!")
		return false
	}
	log.Printf("更新成功")
	return true
}
