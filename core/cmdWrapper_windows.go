//go:build windows

package core

import "syscall"

func ApplySysProcAttr(cmd *Cmd) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		return
	}
	cmd.SysProcAttr.HideWindow = true
}

func ApplySysProcAttrVisible(cmd *Cmd) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
		return
	}
	cmd.SysProcAttr.HideWindow = false
}
