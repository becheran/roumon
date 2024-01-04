package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime/debug"

	"github.com/becheran/roumon/internal/client"
	"github.com/becheran/roumon/internal/model"
	"github.com/becheran/roumon/internal/ui"
)

func main() {
	var host string
	var dbgFile string
	var port int
	flag.StringVar(&host, "host", "localhost", "The pprof server IP or hostname")
	flag.IntVar(&port, "port", 6060, "The pprof server port")
	flag.StringVar(&dbgFile, "debug", "", "Path to debug file")
	flag.Parse()

	version := "dev"
	if info, ok := debug.ReadBuildInfo(); ok {
		version = info.Main.Version
	}

	if len(dbgFile) > 0 {
		f, err := os.OpenFile(dbgFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer f.Close()
		log.SetOutput(f)
	} else {
		log.SetOutput(ioutil.Discard)
	}

	log.Printf("Start roumon (%s)", version)

	c := client.NewClient(host, port)
	ui := ui.NewUI()

	terminate := make(chan error)

	routinesUpdate := make(chan []model.Goroutine)
	go c.Run(terminate, routinesUpdate)
	go ui.Run(terminate, routinesUpdate)

	err := <-terminate
	ui.Stop()

	if err != nil {
		fmt.Println(err.Error())
		log.Printf(err.Error())
	}

	log.Print("Stopped")
}
