package ipmibot

import (
	"io"
	"io/ioutil"
	//"log"
	"net/http"
)

func ipmi(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	data := parseInputJson(body)
	_, _ = err, data
	out := IpmiExec("10.20.30.13", "off")
	io.WriteString(w, jsonFormatReply("green", out.result.output))
}

var mux map[string]func(http.ResponseWriter, *http.Request)

func main() {
	server := http.Server{
		Addr:    ":8000",
		Handler: &myHandler{},
	}

	mux = make(map[string]func(http.ResponseWriter, *http.Request))
	mux["/ipmi"] = ipmi

	server.ListenAndServe()
}

type myHandler struct{}

func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h, ok := mux[r.URL.String()]; ok {
		h(w, r)
		return
	}

	io.WriteString(w, help())
}

func help() string {
	hstr := `Valid commands are:
		help - list of help topics, for more information type /ipmi help <topic>
		reboot <ip or alias>
		off < ip or alias>
		on < ip or alias>
		lanboot < ip or alias>
		status < ip or alias>
		alias [ - | add | del | show ] [<alias name>]
		last`
	return hstr
}
