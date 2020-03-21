package media

import (
	"bytes"
	"fmt"
	"os/exec"
	"time"
)

func cmdExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func runCmd(timeout int, command string, args ...string) error {
	// instantiate new command
	cmd := exec.Command(command, args...)

	// get pipe to standard output
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("cmd.StdoutPipe() error: %w", err)
	}

	// start process via command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("cmd.Start() error: %w", err)
	}

	// setup a buffer to capture standard output
	var buf bytes.Buffer

	// create a channel to capture any errors from wait
	done := make(chan error)
	go func() {
		if _, err := buf.ReadFrom(stdout); err != nil {
			panic("buf.Read(stdout) error: " + err.Error())
		}
		done <- cmd.Wait()
	}()

	// block on select, and switch based on actions received
	select {
	case <-time.After(time.Duration(timeout) * time.Second):
		if err := cmd.Process.Kill(); err != nil {
			return fmt.Errorf("failed to kill: %w", err)
		}
		return fmt.Errorf("timeout reached, process killed")
	case err := <-done:
		if err != nil {
			close(done)
			return fmt.Errorf("process done, with error: %w", err)
		}
		return nil
	}
}
