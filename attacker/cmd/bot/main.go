package main

import (
	"log"

	"github.com/dotping-me/simulated-ddos/bot"
)

func main() {
	if err := bot.NewCLI().Execute(); err != nil {
		log.Fatal(err)
	}
}
