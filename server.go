package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

var address string = ""
var port int = 8000

func ipmi(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	data := parseInputJson(body)
	_, _ = err, data
	var output string
	var color string = "green"
	switch data.command {
	case "help":
		output = help(data.args)
		color = "blue"
	case "reboot", "on", "off", "lanboot", "status":
		log.Println(data.node)
		if len(data.args) == 0 || data.args[0] == "" {
			output = help([]string{data.command})
			color = "purple"
		} else {
			out := IpmiExec(data.args[0], data.command)
			output = out.result.output
		}
	case "alias":
	}

	io.WriteString(w, jsonFormatReply(color, output))
}

var mux map[string]func(http.ResponseWriter, *http.Request)

func main() {
	var addr string
	addr = fmt.Sprintf("%s:%d", address, port)

	server := http.Server{
		Addr:    addr,
		Handler: &myHandler{},
	}

	mux = make(map[string]func(http.ResponseWriter, *http.Request))
	mux["/ipmi"] = ipmi

	server.ListenAndServe()
	log.Println("Server listening on address", addr)
}

type myHandler struct{}

func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
		last`

	topics := map[string]string{
		"help":    "Shows help messages, for detailed help by topics type /ipmi help <topic>",
		"reboot":  "Usage:\n/ipmi reboot <ip or alias>\nReboots host by ip-address or alias",
		"off":     "Usage:\n/ipmi off <ip or alias>\nSwitches off power for host by ip-address or alias",
		"on":      "Usage:\n/ipmi on <ip or alias>\nTurns on host with ip-address or alias",
		"lanboot": "Usage:\n/ipmi lanboot <ip or alias>\nSets network boot for host or alias and reboots it",
		"status":  "Usage:\n/ipmi status <ip or alias>\nShows current chassis power status for host or alias",
		"alias": `/Usage:
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
