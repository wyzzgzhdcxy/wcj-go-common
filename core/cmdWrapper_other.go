//go:build !windows

package core

func ApplySysProcAttr(cmd *Cmd) {}

func ApplySysProcAttrVisible(cmd *Cmd) {}
