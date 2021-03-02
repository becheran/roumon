package main

import (
	"flag"
	"log"
	"os"

	"github.com/becheran/roumon/internal/client"
	"github.com/becheran/roumon/internal/model"
	"github.com/becheran/roumon/internal/ui"
)

func main() {
	var logFile string
	var host string
	var port int
	flag.StringVar(&logFile, "log", "", "Path to logfile")
	flag.StringVar(&host, "host", "localhost", "The pprof server IP or hostname. Default is 'localhost'")
	flag.IntVar(&port, "port", 6060, "The pprof server port. Default is 6060")
	flag.Parse()

	if len(logFile) > 0 {
		f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer f.Close()
		log.SetOutput(f)
	}

	log.Print("Start")

	c := client.NewClient(host, port)
	ui := ui.NewUI()
	defer ui.Stop()
	terminate := make(chan error)

	routinesUpdate := make(chan []model.Goroutine)
	go c.Run(terminate, routinesUpdate)
	go ui.Run(terminate, routinesUpdate)

	err := <-terminate

	if err != nil {
		log.Printf("Error: %s\n", err.Error())
	}

	log.Print("Stopped")
}
