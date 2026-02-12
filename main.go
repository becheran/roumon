package main

import (
	"flag"
	"fmt"
	"io"
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
	var versionFlag bool
	flag.StringVar(&host, "host", "localhost", "The pprof server IP or hostname")
	flag.IntVar(&port, "port", 6060, "The pprof server port")
	flag.StringVar(&dbgFile, "debug", "", "Path to debug file")
	flag.BoolVar(&versionFlag, "v", false, "Print version of roumon and exit")
	flag.Parse()

	version := "dev"
	if info, ok := debug.ReadBuildInfo(); ok {
		version = info.Main.Version
	}

	if versionFlag {
		fmt.Println(version)
		return
	}

	if len(dbgFile) > 0 {
		f, err := os.OpenFile(dbgFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer func() {
			if err := f.Close(); err != nil {
				log.Printf("error closing file: %v", err)
			}
		}()
		log.SetOutput(f)
	} else {
		log.SetOutput(io.Discard)
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
		log.Print(err.Error())
	}

	log.Print("Stopped")
}
