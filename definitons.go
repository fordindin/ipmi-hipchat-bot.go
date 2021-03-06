package main

/*
 * definitions of all static values, internal types and global variables
 */

import (
	"flag"
	"fmt"
	"net"
)

var signal = flag.String("s", "", `send signal to the daemon
		quit — graceful shutdown
		stop — fast shutdown
		reload — reloading the configuration file`)
var configPath string

var ipmitool string
var ipmitoolBinErr error

var listener net.Listener

var daemonName = "ipmi-hipchat-gobot"
var processName = []string{fmt.Sprintf("[%s]", daemonName)}

type Config struct {
	Address      string
	Port         int
	Pidfile      string
	Logfile      string
	Workdir      string
	Ipmiusername string
	Ipmipassword string
	Dbpath       string
	Ipmitoolpath string
}

var config Config = Config{
	Address:      "",
	Port:         8000,
	Pidfile:      fmt.Sprintf("%s.pid", daemonName),
	Logfile:      fmt.Sprintf("%s.log", daemonName),
	Workdir:      ".",
	Ipmiusername: "ADMIN",
	Ipmipassword: "ADMIN",
	Dbpath:       "./ipmibot.sqlite3",
	Ipmitoolpath: "/usr/local/bin/ipmitool",
}

var commands = map[string][][]string{
	"status": [][]string{[]string{"chassis", "power", "status"}},
	"off":    [][]string{[]string{"chassis", "power", "off"}},
	"on":     [][]string{[]string{"chassis", "power", "on"}},
	"cycle":  [][]string{[]string{"chassis", "power", "cycle"}},
	"lanboot": [][]string{[]string{"chassis", "bootdev", "pxe"},
		[]string{"chassis", "power", "cycle"}},
}

var dbversion = 0

type execReturn struct {
	output        string
	commandstring string
	err           error
}

type hipchatMessage struct {
	node       string
	name       string
	command    string
	chatstring string
	args       []string
}

type dbLogEntry struct {
	timestamp     int
	caller        string
	chatstring    string
	chatout       string
	systemcommand string
	systemout     string
	systemerror   error
}

type dbAliasEntry struct {
	name  string
	owner string
	host  string
}

type httpHandler struct{}

var helpstr string = `Valid commands are:
		help - list of help topics, for more information type /ipmi help <topic>
		reboot <ip or alias>
		off <ip or alias>
		on <ip or alias>
		lanboot <ip or alias>
		status <ip or alias>
		alias [ - | add | del | show ] [<alias name>]
		last [<number>]`

var topics map[string]string = map[string]string{
	"help":    "Shows help messages, for detailed help by topics type /ipmi help <topic>",
	"reboot":  "Usage:\n/ipmi reboot <ip or alias>\nReboots host by ip-address or alias",
	"off":     "Usage:\n/ipmi off <ip or alias>\nSwitches off power for host by ip-address or alias",
	"on":      "Usage:\n/ipmi on <ip or alias>\nTurns on host with ip-address or alias",
	"lanboot": "Usage:\n/ipmi lanboot <ip or alias>\nSets network boot for host or alias and reboots it",
	"status":  "Usage:\n/ipmi status <ip or alias>\nShows current chassis power status for host or alias",
	"last":    "Usage:\n/ipmi last [<number>]\nShows last <number> executed commands, default is ten commands",
	"alias": `Usage:
ipmi alias [ - | add | del | show ] [<alias name>] [<ip address>]

/ipmi alias
Shows list of aliases for current user

/ipmi alias add <alias name> <ip address>
Adds alias <alias name> for <ip address>

/ipmi alias del <alias name>
Deletes alias <alias name>

/ipmi alias show <alias name>
Shows ip address for alias <alias name>
`}
