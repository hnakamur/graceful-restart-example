package main

import (
	"flag"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/lomik/go-daemon"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/graceful"
	"github.com/zenazn/goji/web"
)

func getSleep(w http.ResponseWriter, r *http.Request) {
	time.Sleep(3 * time.Second)
	fmt.Fprintf(w, "Hello, %q\n", html.EscapeString(r.URL.Path))
}

func route(m *web.Mux) {
	m.Get("/sleep", getSleep)
}

func main() {
	var addr string
	flag.StringVar(&addr, "addr", ":7777", "listen address (host:port)")
	var pidFileName string
	flag.StringVar(&pidFileName, "pidfile", "/tmp/exampleapp.pid", "pid file path")
	var isDaemon bool
	flag.BoolVar(&isDaemon, "daemon", false, "background runnning daemon")
	var logPath string
	flag.StringVar(&logPath, "log-path", "/tmp/examplewebapp.log", "log file path")
	flag.Parse()

	log.Printf("examplewebapp start. daemon=%v", isDaemon)

	if isDaemon == true {
		dmn := new(daemon.Context)
		child, _ := dmn.Reborn()

		if child != nil {
			return
		} else {
			_, err := daemon.CreatePidFile(pidFileName, 0644)
			if err != nil {
				log.Println("pidfile create error")
				panic(err)
			}
			log.Println("start daemon child prossess")
		}
	}

	logFile, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.Printf("examplewebapp start#2. daemon=%v", isDaemon)

	route(goji.DefaultMux)
	log.Println("Start goji.Serve")
	graceful.AddSignal(syscall.SIGTERM)
	goji.Serve()
}
