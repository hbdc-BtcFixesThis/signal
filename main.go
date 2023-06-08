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

func (ss *SignalServer) initBuckets() error {
	// Start a writable transaction.
	tx, err := ss.DB.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Use the transaction...
	_, err = tx.CreateBucket([]byte("records"))
	if err != nil {
		return err
	}

	// Commit the transaction and check for error.
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func main() {
	log.SetFlags(0)

	sc, err := NewServerConf()
	if err != nil {
		log.Fatal("ERROR:", err)
	}

	ss := newSignalServer(sc)
	ss.Run()
}
