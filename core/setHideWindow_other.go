//go:build !windows

package core

import "os/exec"

func SetHideWindow(cmd *exec.Cmd) {
}
