package main

func IpmiExec(host string, command string) execReturn {
	var ret execReturn
	var cmdArray []string
	cmdArray = append(cmdArray,
		"-U", ipmiUserame,
		"-P", ipmiPassword,
		"-H", host)
	cmdArray = append(cmdArray, commands[command]...)
	ret = execCommand(ipmitool, cmdArray...)
	return ret
}
