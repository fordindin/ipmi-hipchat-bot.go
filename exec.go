package main

/*
 * Wrapper for os/exec
 */

import (
	"fmt"
	"os/exec"
	"strings"
)

// execCommand executes single system command, waits for it's completion
// and returns result as struct execReturn
func execCommand(command string, args ...string) execReturn {
	cmd := exec.Command(command, args...)
	err := cmd.Wait()
	stdout, err := cmd.CombinedOutput()
	return execReturn{
		commandstring: fmt.Sprintf("%s %s", command, strings.Join(args, " ")),
		output:        string(stdout),
		err:           err,
	}
}
