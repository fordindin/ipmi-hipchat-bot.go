package main

// Wrapper for execution ipmi-related system commands
func IpmiExec(host string, command string) []execReturn {
	var ret []execReturn
	var cmdArray []string
	cmdArray = append(cmdArray,
		"-U", config.Ipmiusername,
		"-P", config.Ipmipassword,
		"-H", host)
	for _, c := range commands[command] {
		cArray := append(cmdArray, c...)
		ret = append(ret, execCommand(ipmitool, cArray...))
	}
	return ret
}
