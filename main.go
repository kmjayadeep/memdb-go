package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kmjayadeep/memdb-go/server"
)

func main() {
	fmt.Println("Starting memdb")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	srv := server.NewServer()

	select {
	case <-stop:
		srv.Stop()
	}
}
