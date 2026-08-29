package core

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/sys/windows"
)

// FormatUsageShow 优化命令行描述展示更加整齐
func FormatUsageShow(usageText string) {
	_, _ = fmt.Fprintf(os.Stderr, "%s\n", usageText)
	// 遍历所有命令行参数
	fmt.Println("命令:")
	flag.VisitAll(func(f *flag.Flag) {
		if !(f.Usage[0:1] == "1") {
			return
		}
		line := fmt.Sprintf("  -%-10s %s [%s]", f.Name, strings.ReplaceAll(f.Usage[2:], "\n", " "), f.DefValue)
		_, _ = fmt.Fprintln(os.Stderr, line)
	})
	fmt.Println()
	fmt.Println("参数:")
	flag.VisitAll(func(f *flag.Flag) {
		if f.Usage[0:1] == "1" {
			return
		}
		line := fmt.Sprintf("  -%-10s %s [%s]", f.Name, strings.ReplaceAll(f.Usage[2:], "\n", " "), f.DefValue)
		_, _ = fmt.Fprintln(os.Stderr, line)
	})
}

func Open(url string) {
	cmd := exec.Command("explorer", url)
	err := cmd.Start()
	if err != nil {
		log.Printf(err.Error())
	}
}

func ExecuteCmderLine(cmdStr string) {
	Execute(cmdStr, ExecPath())
}

func Execute(cmdStr string, dir string) int {
	cmd := CompileCmd(cmdStr)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if StrIsNotEmpty(dir) {
		cmd.Dir = dir
	}
	//cmd.SysProcAttr = &syscall.SysProcAttr{
	//	HideWindow: true, // 这将创建一个新的会话组
	//}
	err := cmd.Start()
	pid := cmd.Process.Pid
	if err != nil {
		fmt.Printf("执行命令错误：%s,%s\n", cmdStr, err)
	}
	// 等待命令执行完成
	err = cmd.Wait()
	if err != nil {
		fmt.Printf("执行命令错误：%s,%s\n", cmdStr, err)
	}
	return pid
}

func CompileCmd(cmdStr string) *exec.Cmd {
	if IsWin() {
		return exec.Command("cmd", "/c", cmdStr)
	} else {
		return exec.Command("bash", "-c", cmdStr)
	}
}

func ExecuteCommand(name string, arg ...string) string {
	return string(*ExecuteCommandReturnByte(name, arg...))
}

func ExecuteCommandReturnByte(name string, arg ...string) *[]byte {
	// 构建FFmpeg命令
	cmd := exec.Command(name, arg...)
	cmd.SysProcAttr = &windows.SysProcAttr{
		HideWindow: true,
	}
	// 运行命令并捕获输出和错误
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error running: %v\n", err)
		fmt.Printf("output:\n%s\n", string(output))
	}
	return &output
}

func ExecuteCommandByTargetDir(targetDir, name string, arg ...string) *[]byte {
	// 构建FFmpeg命令
	cmd := exec.Command(name, arg...)
	cmd.SysProcAttr = &windows.SysProcAttr{
		HideWindow: true,
	}
	cmd.Dir = targetDir
	// 运行命令并捕获输出和错误
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error running: %v\n", err)
		fmt.Printf("output:\n%s\n", string(output))
	}
	return &output
}
