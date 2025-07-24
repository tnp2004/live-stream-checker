package main

import (
	"log"

	"github.com/tnp2004/live-stream-checker/terminal"
)

func main() {
	if err := terminal.Run(); err != nil {
		log.Fatal(err)
	}
}
