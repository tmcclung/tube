package utils

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"
)

// FileExists ...
func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// CmdExists ...
func CmdExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// RunCmd ...
func RunCmd(timeout int, command string, args ...string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, command, args...)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("cmd.CombinedOutput error: %w\n%s", err, out)
	}

	return nil
}
