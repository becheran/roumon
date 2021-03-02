package client_test

import (
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/becheran/roumon/internal/client"
	"github.com/becheran/roumon/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestEmptyResponse(t *testing.T) {
	const testport = 6062

	// test server
	go func() {
		err := http.ListenAndServe(fmt.Sprintf("localhost:%d", testport), nil)
		assert.Nil(t, err)
	}()

	testClient := client.NewClient("localhost", testport)

	done := make(chan error)
	routines := make(chan []model.Goroutine)

	go testClient.Run(done, routines)
	select {
	case r := <-routines:
		assert.Empty(t, r)
	case <-done:
		log.Fatal("Failed")
	}
}
