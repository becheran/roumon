package client

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/becheran/roumon/internal/model"
)

type Client struct {
	c      *http.Client
	server string
}

func NewClient(ip string, port int) *Client {
	server := fmt.Sprintf("http://%s:%d/debug/pprof/goroutine?debug=2", ip, port)
	log.Printf("Attach to server %s\n", server)
	c := &http.Client{}
	return &Client{
		c:      c,
		server: server,
	}
}

func (client *Client) Run(terminate chan<- error, routineUpdate chan<- []model.Goroutine) {
	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()

	for {
		resp, err := client.c.Get(client.server)
		if err != nil {
			terminate <- fmt.Errorf("Failed to list go routines. Err: %s", err.Error())
			return
		}
		defer resp.Body.Close()

		goroutines, err := model.ParseStackFrame(resp.Body)
		if err != nil {
			log.Printf("Error while parsing stack: %s", err.Error())
		} else {
			routineUpdate <- goroutines
		}
		<-ticker.C
	}
}
