package main

import "os/exec"

var address string = ""
var port int = 8000

var ipmiUserame string = "ADMIN"
var ipmiPassword string = "ADMIN"

var commands = map[string][]string{
	"status": []string{"chassis", "power", "status"},
	"off":    []string{"chassis", "power", "off"},
	"on":     []string{"chassis", "power", "on"},
	"cycle":  []string{"chassis", "power", "cycle"},
}

var ipmitool, ipmitoolBinErr = exec.LookPath("ipmitool")

var dbaddr string = "./ipmibot.sqlite3"
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
