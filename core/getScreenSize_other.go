//go:build !windows

package core

// getScreenSize 非 Windows 平台返回默认值。
// 如需精确屏幕尺寸，可在 Mac 上使用 github.com/kbinani/screen 等第三方库。
func GetScreenSize() (width, height int) {
	return 1920, 1080
}
