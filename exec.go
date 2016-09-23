package main

import (
	"fmt"
	//"log"
	"os/exec"
	"strings"
)

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
