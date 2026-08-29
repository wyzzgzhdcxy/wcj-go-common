//go:build !windows

package core

import (
	"fmt"
	"os/exec"
)

func ExecuteCommandReturnByte(name string, arg ...string) *[]byte {
	cmd := exec.Command(name, arg...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error running: %v\n", err)
		fmt.Printf("output:\n%s\n", string(output))
	}
	return &output
}

func ExecuteCommandByTargetDir(targetDir, name string, arg ...string) *[]byte {
	cmd := exec.Command(name, arg...)
	cmd.Dir = targetDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error running: %v\n", err)
		fmt.Printf("output:\n%s\n", string(output))
	}
	return &output
}
