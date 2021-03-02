package client_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/becheran/roumon/internal/client"
	"github.com/becheran/roumon/internal/model"
)

func TestPrint(t *testing.T) {
	c := client.NewClient("localhost", 6060)
	done := make(chan error)
	routines := make(chan []model.Goroutine)

	go c.Run(done, routines)
	select {
	case r := <-routines:
		fmt.Println(r)
	case <-done:
		log.Fatal("Failed")
	}
}
