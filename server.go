package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func ipmi(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	data := parseInputJson(body)
	_, _ = err, data

	var logEntry dbLogEntry
	var output string
	var color string = "green"
	var out execReturn

	log.Printf("Command %s received from %s\n", data.command, data.name)
	switch data.command {
	case "help":
		output = help(data.args)
		color = "blue"
	case "reboot", "on", "off", "lanboot", "status":
		if len(data.args) == 0 || data.args[0] == "" {
			output = help([]string{data.command})
			color = "purple"
		} else {
			out = IpmiExec(data.args[0], data.command)
			output = out.output
		}

	case "alias":
		switch {
		case len(data.args) == 0:
			aliasEntry := mkDbAliasEntry("", data.name, "")
			aliasEntries := showAlias(*aliasEntry)
			output = fmt.Sprintf("Aliases for you:\n")
			for _, a := range aliasEntries {
				output += aliasToString(a)
			}
		case len(data.args) == 2 && data.args[0] == "del":
			output = data.args[0]
		case len(data.args) == 2 && data.args[0] == "show":
			aliasEntry := mkDbAliasEntry(data.args[1], data.name, "")
			aliasEntries := showAlias(*aliasEntry)
			if len(aliasEntries) == 0 {
				output = fmt.Sprintf("There is no alias '%s' for %s\n", data.args[2], data.name)
				break
			}
			output = aliasToString(aliasEntries[0])
		case len(data.args) == 3 && data.args[0] == "add":
			aliasEntry := mkDbAliasEntry(data.args[1], data.name, data.args[2])
			ip := net.ParseIP(data.args[2])
			if ip == nil {
				output = fmt.Sprintf("'%s' doesn't look like IP-address to me\n", data.args[2])
				break
			}
			updateAlias(*aliasEntry)
			aliasEntries := showAlias(*aliasEntry)
			if len(aliasEntries) == 0 {
				output = "Something wrong, I cannot add alias =("
				break
			}
			output = aliasToString(aliasEntries[0])
		default:
			output = help([]string{data.command})
			color = "purple"
		}
	}

	output = jsonFormatReply(color, output)

	logEntry.caller = data.name
	logEntry.chatstring = data.chatstring
	logEntry.chatout = output
	logEntry.systemcommand = string(out.commandstring)
	logEntry.systemout = string(out.output)
	logEntry.systemerror = out.err

	logToDB(logEntry)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, output)
}

var mux map[string]func(http.ResponseWriter, *http.Request)

func (*httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h, ok := mux[r.URL.String()]; ok {
		h(w, r)
		return
	}

	io.WriteString(w, help(make([]string, 0)))
}

func help(args []string) string {
	helpstr := `Valid commands are:
		help - list of help topics, for more information type /ipmi help <topic>
		reboot <ip or alias>
		off <ip or alias>
		on <ip or alias>
		lanboot <ip or alias>
		status <ip or alias>
		alias [ - | add | del | show ] [<alias name>]
		last [<number>]`

	topics := map[string]string{
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
	if len(args) == 1 {
		return topics[args[0]]
	}
	return helpstr
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func() {
		sig := <-c
		DB.Close()
		log.Println("Signal received, exiting:", sig)
		os.Exit(0)
	}()

	var addr string
	initDB()
	addr = fmt.Sprintf("%s:%d", address, port)

	server := http.Server{
		Addr:    addr,
		Handler: &httpHandler{},
	}

	mux = make(map[string]func(http.ResponseWriter, *http.Request))
	mux["/ipmi"] = ipmi

	server.ListenAndServe()
	log.Println("Server listening on address", addr)
}
