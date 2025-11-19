package main

import (
	"flag"
	"fmt"
	"log"
)

func init() {
	flag.BoolVar(&enableDebug, "debug", false, "Enable debug mode")
	flag.Parse() // Parse the command-line flags
	fmt.Println("Debug mode set to: ", enableDebug)
}

func main() {
	log.SetFlags(0)
	ss, err := newSignalServer()
	if err != nil {
		panic(err)
	}
	ss.Run()
}
