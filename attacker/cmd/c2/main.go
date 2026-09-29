package main

import (
	"log"

	"github.com/dotping-me/simulated-ddos/c2"
)

func main() {
	m := c2.NewMaster()
	go func() { // Starts HTTP server in non-blocking fashion
		if err := m.Serve("0.0.0.0:8080"); err != nil {
			log.Fatal(err)
		}
	}()

	if err := c2.RunCLI(m); err != nil {
		log.Fatal(err)
	}
}
