// Package cmdWrapper 提供子进程执行与 Go runtime 的统一封装。
// cmdWrapper.go 是整个项目中唯一允许直接 import os/exec 与 runtime 标准库的
// 文件，其他代码一律通过本包的导出方法与类型别名使用，便于集中管控平台差异。
package core

import (
	"bytes"
	"os"
	"os/exec"
	"runtime"
)

// ========== 类型别名：调用方无需直接 import os/exec / runtime ==========

// Cmd 是 exec.Cmd 的类型别名，供调用方声明、传递命令对象。
type Cmd = exec.Cmd

// ExitError 是 exec.ExitError 的类型别名，供调用方做退出错误断言。
type ExitError = exec.ExitError

// MemStats 是 runtime.MemStats 的类型别名，配合 ReadMemStats 使用。
type MemStats = runtime.MemStats

// ========== runtime 封装 ==========

// GOOS 返回目标操作系统标识（windows/darwin/linux...）。
func GOOS() string { return runtime.GOOS }

// IsWindows 报告当前是否为 Windows 平台。
func IsWindows() bool { return runtime.GOOS == "windows" }

// NumCPU 返回本机 CPU 逻辑核心数。
func NumCPU() int { return runtime.NumCPU() }

// GOMAXPROCS 设置并返回可同时执行的 CPU 数上限。
func GOMAXPROCS(n int) int { return runtime.GOMAXPROCS(n) }

// ReadMemStats 将当前内存统计写入 m。
func ReadMemStats(m *MemStats) { runtime.ReadMemStats(m) }

// NumGoroutine 返回当前存活的 goroutine 数量。
func NumGoroutine() int { return runtime.NumGoroutine() }

// SetFinalizer 在对象 obj 被 GC 回收时调用 finalizer。
func SetFinalizer(obj, finalizer interface{}) { runtime.SetFinalizer(obj, finalizer) }

// ========== os/exec 封装 ==========

// LookPath 在 PATH 中查找可执行文件的完整路径。
func LookPath(name string) (string, error) { return exec.LookPath(name) }

// Options holds configuration for command execution
type Options struct {
	Dir        string        // Working directory
	HideWindow bool          // Hide command window on Windows
	Stdout     *bytes.Buffer // Custom stdout buffer (nil uses os.Stdout)
	Stderr     *bytes.Buffer // Custom stderr buffer (nil uses os.Stderr)
}

// defaultOptions returns default options
func defaultOptions() Options {
	return Options{
		Dir:        "",
		HideWindow: IsWindows(), // Hide window by default on Windows
		Stdout:     nil,
		Stderr:     nil,
	}
}

// Option is a function that modifies Options
type Option func(*Options)

// WithDir sets the working directory
func WithDir(dir string) Option {
	return func(o *Options) {
		o.Dir = dir
	}
}

// WithHideWindow sets whether to hide the command window on Windows
func WithHideWindow(hide bool) Option {
	return func(o *Options) {
		o.HideWindow = hide
	}
}

// WithStdout sets a custom stdout buffer
func WithStdout(buf *bytes.Buffer) Option {
	return func(o *Options) {
		o.Stdout = buf
	}
}

// WithStderr sets a custom stderr buffer
func WithStderr(buf *bytes.Buffer) Option {
	return func(o *Options) {
		o.Stderr = buf
	}
}

// Command 返回一个预配置好的 *Cmd：Windows 上自动设置 HideWindow，
// 避免 GUI 程序打包后子进程弹出黑色控制台窗口。调用方仍可自行设置
// Dir/Stdout/Stderr/ExtraFiles 等字段。
func Command(name string, args ...string) *Cmd {
	cmd := exec.Command(name, args...)
	ApplySysProcAttr(cmd)
	return cmd
}

// CommandVisible 返回一个不强制隐藏窗口的 *Cmd。
// 用于 explorer、cmd /c start 等需要显示窗口（资源管理器/浏览器）的场景。
// 隐藏窗口会导致这些 GUI 子进程静默退出而不弹出窗口。
func CommandVisible(name string, args ...string) *Cmd {
	cmd := exec.Command(name, args...)
	ApplySysProcAttrVisible(cmd)
	return cmd
}

// createCommand creates an exec.Cmd with the given options
func createCommand(name string, args []string, opts Options) *Cmd {
	cmd := exec.Command(name, args...)

	if opts.Dir != "" {
		cmd.Dir = opts.Dir
	}

	// Set stdout/stderr
	if opts.Stdout != nil {
		cmd.Stdout = opts.Stdout
	} else {
		cmd.Stdout = os.Stdout
	}

	if opts.Stderr != nil {
		cmd.Stderr = opts.Stderr
	} else {
		cmd.Stderr = os.Stderr
	}

	// Hide window on Windows if requested
	if IsWindows() && opts.HideWindow {
		ApplySysProcAttr(cmd)
	}

	return cmd
}

// Run executes a command and returns any error
func Run(name string, args ...string) error {
	return RunWithOptions(name, args)
}

// RunWithOptions executes a command with options and returns any error
func RunWithOptions(name string, args []string, opts ...Option) error {
	options := defaultOptions()
	for _, opt := range opts {
		opt(&options)
	}

	cmd := createCommand(name, args, options)
	return cmd.Run()
}

// RunWithOutput executes a command and returns combined output
func RunWithOutput(name string, args ...string) (string, error) {
	return RunWithOutputAndOptions(name, args)
}

// RunWithOutputAndOptions executes a command with options and returns combined output
func RunWithOutputAndOptions(name string, args []string, opts ...Option) (string, error) {
	options := defaultOptions()
	for _, opt := range opts {
		opt(&options)
	}

	// For output capture, we need custom buffers
	var stdout, stderr bytes.Buffer
	options.Stdout = &stdout
	options.Stderr = &stderr

	cmd := createCommand(name, args, options)
	err := cmd.Run()

	// Combine stdout and stderr
	output := stdout.String() + stderr.String()
	return output, err
}

// Start starts a command but does not wait for it to complete
func Start(name string, args ...string) error {
	return StartWithOptions(name, args)
}

// StartWithOptions starts a command with options but does not wait for it to complete
func StartWithOptions(name string, args []string, opts ...Option) error {
	options := defaultOptions()
	for _, opt := range opts {
		opt(&options)
	}

	cmd := createCommand(name, args, options)
	return cmd.Start()
}

// RunWithDir is a convenience function to run a command in a specific directory
func RunWithDir(dir, name string, args ...string) error {
	return RunWithOptions(name, args, WithDir(dir))
}

// RunWithDirAndOutput is a convenience function to run a command in a specific directory and get output
func RunWithDirAndOutput(dir, name string, args ...string) (string, error) {
	return RunWithOutputAndOptions(name, args, WithDir(dir))
}
