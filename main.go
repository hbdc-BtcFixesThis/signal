package main

import (
	"log"
)

func main() {
	log.SetFlags(0)
	ss, err := newSignalServer()
	if err != nil {
		panic(err)
	}
	ss.Run()
}
