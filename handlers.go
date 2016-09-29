package main

import (
	"fmt"
	"net"
	"strconv"
)

// function gets HTTP request body as []byte and returns handled JSON as string
// basicaly, this function is central point of data processing
func processIpmi(body []byte) string {
	data := parseInputJson(body)
	var logEntry dbLogEntry
	var output string
	var color string = "green"
	var out []execReturn
	output = fmt.Sprintf("@%s,\n", data.name)

	switch data.command {
	case "help":
		output += help(data.args)
		color = "purple"
	case "last":
		var ncommands int = 10
		var err error
		if len(data.args) > 0 {
			ncommands, err = strconv.Atoi(data.args[0])
			_ = ncommands
			if err != nil {
				output += fmt.Sprintf("'%s' doesn't look like a number to me", data.args[0])
				color = "red"
			}
		}
		if color != "red" {
			output += last(ncommands)
		}
	case "reboot", "on", "off", "lanboot", "status":
		switch {
		case len(data.args) == 0:
			output += help([]string{data.command})
			color = "purple"
		default:
			out = IpmiExec(unwrapAlias(data.args[0], data.name), data.command)
			for _, o := range out {
				output += o.output
			}
		}
	case "alias":
		switch {
		case len(data.args) == 0:
			aliasEntry := mkDbAliasEntry("", data.name, "0.0.0.0")
			aliasEntries := showAlias(*aliasEntry)
			output += fmt.Sprintf("Aliases for you:\n")
			for _, a := range aliasEntries {
				output += aliasToString(a)
				output += "\n"
			}
		case len(data.args) == 2 && data.args[0] == "del":
			aliasEntry := mkDbAliasEntry(data.args[1], data.name, "0.0.0.0")
			aliasEntries := showAlias(*aliasEntry)
			if len(aliasEntries) == 0 {
				output += fmt.Sprintf("There is no alias '%s' for %s\n", data.args[1], data.name)
				color = "red"
				break
			} else {
				delAlias(*aliasEntry)
				aliasEntries := showAlias(*aliasEntry)
				if len(aliasEntries) == 0 {
					output += fmt.Sprintf("Alias '%s' removed (owner  %s)", data.args[1], data.name)
				}
			}
		case len(data.args) == 2 && data.args[0] == "show":
			aliasEntry := mkDbAliasEntry(data.args[1], data.name, "0.0.0.0")
			aliasEntries := showAlias(*aliasEntry)
			if len(aliasEntries) == 0 {
				output += fmt.Sprintf("There is no alias '%s' for %s\n", data.args[1], data.name)
				color = "red"
				break
			}
			output += aliasToString(aliasEntries[0])
		case len(data.args) == 3 && data.args[0] == "add":
			aliasEntry := mkDbAliasEntry(data.args[1], data.name, data.args[2])
			ip := net.ParseIP(data.args[2])
			if ip == nil {
				output += fmt.Sprintf("'%s' doesn't look like IP-address to me\n", data.args[2])
				color = "red"
				break
			}
			updateAlias(*aliasEntry)
			aliasEntries := showAlias(*aliasEntry)
			if len(aliasEntries) == 0 {
				output += "Something wrong, I cannot add alias =("
				color = "red"
				break
			}
			output += aliasToString(aliasEntries[0])
		default:
			output += help([]string{data.command})
		}
	default:
		output += help([]string{data.command})
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

// Returns help topic or general help, depending on input argument
func help(args []string) string {
	if len(args) == 1 {
		if val, ok := topics[args[0]]; ok {
			return val
		}
	}
	return helpstr
}

// Returns last executed commands from the log
func last(params ...int) string {
	nentries := 10
	if len(params) > 0 {
		nentries = params[0]
	}
	lastEntries := lastFromDB(nentries)
	out := "Last executed commands:\n"
	for i, entry := range lastEntries {
		out += fmt.Sprintf("% 2d % 12s: %s\n", i+1, entry.caller, entry.chatstring)
	}
	return out
}
