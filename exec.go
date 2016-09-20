package main

import (
	//"fmt"
	//"log"
	"os/exec"
)

type execReturn struct {
	output string
	err    error
}

func execCommand(command string, args ...string) execReturn {
	cmd := exec.Command(command, args...)
	err := cmd.Wait()
	stdout, err := cmd.CombinedOutput()
	return execReturn{
		output: string(stdout),
		err:    err,
	}
}
