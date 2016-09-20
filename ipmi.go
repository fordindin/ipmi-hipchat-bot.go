package main

import (
	//"fmt"
	"os/exec"
	"strings"
)

type ipmiExecResult struct {
	command string
	result  execReturn
}

var ipmiUserame string = "ADMIN"
var ipmiPassword string = "ADMIN"

var ipmitool, ipmitoolBinErr = exec.LookPath("ipmitool")

var commands = map[string][]string{
	"status": []string{"chassis", "power", "status"},
	"off":    []string{"chassis", "power", "off"},
	"on":     []string{"chassis", "power", "on"},
	"cycle":  []string{"chassis", "power", "cycle"},
}

func IpmiExec(host string, command string) ipmiExecResult {
	var ret ipmiExecResult
	var cmdArray []string
	cmdArray = append(cmdArray,
		"-U", ipmiUserame,
		"-P", ipmiPassword,
		"-H", host)
	cmdArray = append(cmdArray, commands[command]...)
	ret.command = strings.Join(cmdArray, " ")
	ret.result = execCommand(ipmitool, cmdArray...)
	return ret
}
