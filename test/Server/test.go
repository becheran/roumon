package main

import (
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"
)

// Start a test server which can be used to monitor go routines
func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	go func() {
		for {
			time.Sleep(1 * time.Second)
			go func() {
				min := 10
				max := 5000
				randSleepMSec := rand.Intn(max-min) + min
				time.Sleep(time.Duration(randSleepMSec) * time.Millisecond)
				println("update")
				randSleepMSec = rand.Intn(max-min) + min
				time.Sleep(time.Duration(randSleepMSec) * time.Millisecond)
			}()
		}
	}()
	c := make(chan os.Signal, 1)
	<-c
}
