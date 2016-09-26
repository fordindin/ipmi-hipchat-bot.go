package main

import (
	"fmt"
	"net"
)

func processIpmi(body []byte) string {

	data := parseInputJson(body)
	var logEntry dbLogEntry
	var output string
	var color string = "green"
	var out []execReturn

	//log.Printf("Command %s received from %s\n", data.command, data.name)
	switch data.command {
	case "help":
		output = help(data.args)
		color = "purple"
	case "last":
	case "reboot", "on", "off", "lanboot", "status":
		switch {
		case len(data.args) == 0:
			output = help([]string{data.command})
			color = "purple"
		default:
			out = IpmiExec(data.args[0], data.command)
			for _, o := range out {
				output += o.output
			}
		}
	case "alias":
		switch {
		case len(data.args) == 0:
			aliasEntry := mkDbAliasEntry("", data.name, "0.0.0.0")
			aliasEntries := showAlias(*aliasEntry)
			output = fmt.Sprintf("Aliases for you:\n")
			for _, a := range aliasEntries {
				output += aliasToString(a)
				output += "\n"
			}
		case len(data.args) == 2 && data.args[0] == "del":
			aliasEntry := mkDbAliasEntry(data.args[1], data.name, "0.0.0.0")
			aliasEntries := showAlias(*aliasEntry)
			if len(aliasEntries) == 0 {
				output = fmt.Sprintf("There is no alias '%s' for %s\n", data.args[1], data.name)
				color = "red"
				break
			} else {
				delAlias(*aliasEntry)
				aliasEntries := showAlias(*aliasEntry)
				if len(aliasEntries) == 0 {
					output = fmt.Sprintf("Alias '%s' removed (owner  %s)", data.args[1], data.name)
				}
			}
		case len(data.args) == 2 && data.args[0] == "show":
			aliasEntry := mkDbAliasEntry(data.args[1], data.name, "0.0.0.0")
			aliasEntries := showAlias(*aliasEntry)
			if len(aliasEntries) == 0 {
				output = fmt.Sprintf("There is no alias '%s' for %s\n", data.args[1], data.name)
				color = "red"
				break
			}
			output = aliasToString(aliasEntries[0])
		case len(data.args) == 3 && data.args[0] == "add":
			aliasEntry := mkDbAliasEntry(data.args[1], data.name, data.args[2])
			ip := net.ParseIP(data.args[2])
			if ip == nil {
				output = fmt.Sprintf("'%s' doesn't look like IP-address to me\n", data.args[2])
				color = "red"
				break
			}
			updateAlias(*aliasEntry)
			aliasEntries := showAlias(*aliasEntry)
			if len(aliasEntries) == 0 {
				output = "Something wrong, I cannot add alias =("
				color = "red"
				break
			}
			output = aliasToString(aliasEntries[0])
		default:
			output = help([]string{data.command})
		}
	}

	output = jsonFormatReply(color, output)

	logEntry.caller = data.name
	logEntry.chatstring = data.chatstring
	logEntry.chatout = output
	for _, o := range out {
		logEntry.systemcommand += o.commandstring
		logEntry.systemout += string(o.output)
		logEntry.systemerror = o.err
	}

	logToDB(logEntry)
	return output
}
