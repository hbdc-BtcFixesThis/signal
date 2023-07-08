package main

import (
	"log"
)

type Record struct {
	Name      string `json:"name"`
	Type      int    `json:"type"`
	Content   string `json:"content"`
	PublicKey string `json:"pub_key"`
	Signature string `json:"signature"`
}

func main() {
	log.SetFlags(0)
	ss, err := newSignalServer()
	if err != nil {
		panic(err)
	}
	ss.Run()
}
