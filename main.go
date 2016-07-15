package main

import (
	"flag"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lomik/go-daemon"
)

func sleepHandleFunc(w http.ResponseWriter, r *http.Request) {
	time.Sleep(3 * time.Second)
	fmt.Fprintf(w, "Hello, %q\n", html.EscapeString(r.URL.Path))
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

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGUSR1, syscall.SIGUSR2)
	go func() {
		log.Println("Start signal handler goroutine")
		for {
			s := <-c
			switch s {
			case syscall.SIGINT, syscall.SIGHUP, syscall.SIGUSR1, syscall.SIGUSR2:
				log.Println("Got signal:", s)
			case syscall.SIGTERM:
				log.Println("Got sigterm signal:", s)
				os.Exit(0)
			default:
				log.Println("Got unexpected signal:", s)
			}
		}
	}()

	http.HandleFunc("/sleep", sleepHandleFunc)
	log.Println("Start ListenAndServe")
	log.Fatal(http.ListenAndServe(addr, nil))
}
