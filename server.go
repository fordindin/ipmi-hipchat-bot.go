package main

import (
	"flag"
	"fmt"
	"github.com/sevlyar/go-daemon"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"syscall"
)

// http server reader-writer
func serveIpmi(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	output := processIpmi(body)
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

func worker() {
	var addr string
	var err error
	initDB()
	addr = fmt.Sprintf("%s:%d", config.Address, config.Port)
	listener, err = net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	handler := http.NewServeMux()
	handler.HandleFunc("/ipmi", serveIpmi)

	http.Serve(listener, handler)
	log.Println("Server listening on address", addr)
}

func main() {
	flag.StringVar(&configPath, "c", "/usr/local/etc/ipmi-hipchat-bot.cfg", "path to the configuration file")
	flag.Parse()
	daemon.AddCommand(daemon.StringFlag(signal, "quit"), syscall.SIGQUIT, termHandler)
	daemon.AddCommand(daemon.StringFlag(signal, "stop"), syscall.SIGTERM, termHandler)
	daemon.AddCommand(daemon.StringFlag(signal, "reload"), syscall.SIGHUP, reloadHandler)

	readConfig(configPath)

	if ipmitoolBinErr != nil {
		log.Fatal("ipmitool error", ipmitoolBinErr)
	}

	cntxt := &daemon.Context{
		PidFileName: config.Pidfile,
		PidFilePerm: 0644,
		LogFileName: config.Logfile,
		LogFilePerm: 0640,
		WorkDir:     config.Workdir,
		Umask:       027,
		Args:        processName,
	}

	if len(daemon.ActiveFlags()) > 0 {
		d, err := cntxt.Search()
		if err != nil {
			log.Fatalln("Unable send signal to the daemon:", err)
		}
		daemon.SendCommands(d)
		return
	}

	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatalln(err)
	}
	if d != nil {
		return
	}
	defer cntxt.Release()

	log.Println("- - - - - - - - - - - - - - -")
	log.Println("daemon started")
	log.Printf("Configuration file: %s\n", configPath)
	log.Printf("Pidfile: %s\n", config.Pidfile)
	log.Printf("Database: %s\n", config.Dbpath)

	go worker()

	err = daemon.ServeSignals()
	if err != nil {
		log.Println("Error:", err)
	}
	log.Println("daemon terminated")
}

var (
	stop = make(chan struct{})
	done = make(chan struct{})
)

func termHandler(sig os.Signal) error {
	log.Println("terminating...")
	//stop <- struct{}{}
	if sig == syscall.SIGQUIT || sig == syscall.SIGTERM || sig == os.Interrupt || sig == os.Kill {
		DB.Close()
		listener.Close()
		os.Exit(0)
	}
	return daemon.ErrStop
}

func reloadHandler(sig os.Signal) error {
	log.Println("configuration reloaded")
	return nil
}
