package main

import (
	"flag"
	"fmt"
	"html"
	"log"
	"net"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/lestrrat/go-server-starter/listener"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/graceful"
	"github.com/zenazn/goji/web"
)

func getSleep(w http.ResponseWriter, r *http.Request) {
	time.Sleep(3 * time.Second)
	fmt.Fprintf(w, "Hello, %q\n", html.EscapeString(r.URL.Path))
}

func route(m *web.Mux) {
	m.Get("/sleep/:sleepid", getSleep)
}

func main() {
	var addr string
	flag.StringVar(&addr, "addr", ":7777", "listen address (host:port)")
	var logPath string
	flag.StringVar(&logPath, "log-path", "/tmp/graceful-restart-example.log", "log file path")
	flag.Parse()

	logFile, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	route(goji.DefaultMux)

	graceful.AddSignal(syscall.SIGTERM)

	var l net.Listener
	if os.Getenv("SERVER_STARTER_PORT") != "" {
		listeners, err := listener.ListenAll()
		if err != nil {
			log.Printf("error in ListenAll. err=%s", err)
			return
		}
		if len(listeners) == 0 {
			log.Printf("no listeners")
			return
		}
		log.Printf("listeners count=%d", len(listeners))
		l = listeners[0]
	} else {
		l, err = net.Listen("tcp", addr)
		if err != nil {
			log.Printf("error in listen. err=%s", err)
		}
	}
	goji.ServeListener(l)
}
