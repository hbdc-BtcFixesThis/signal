package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"os"
	"time"

	"encoding/json"
	"io/fs"
	"net/http"
	"os/signal"

	bolt "go.etcd.io/bbolt"
)

//go:embed static
var static embed.FS

type SignalServer struct {
	// subscriberMessageBuffer controls the max number
	// of messages that can be queued for a subscriber
	// before it is kicked.
	//
	// Defaults to 16.
	// subscriberMessageBuffer int

	// publishLimiter controls the rate limit applied to the publish endpoint.
	//
	// Defaults to one publish every 100ms with a burst of 8.
	//publishLimiter *rate.Limiter

	// logf controls where logs are sent.
	// Defaults to log.Printf.
	logf func(f string, v ...interface{})

	// serveMux routes the various endpoints to the appropriate handler.
	serveMux http.ServeMux
	DB       *bolt.DB
	SC       *ServerConfig

	// subscribersMu sync.Mutex
	// subscribers   map[*subscriber]struct{}
}

// newChatServer constructs a chatServer with the defaults.
func newSignalServer(sc *ServerConfig) *SignalServer {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open("my.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	ss := &SignalServer{
		logf: log.Printf,
		DB:   db,
		SC:   sc,
		// subscriberMessageBuffer: 16,
		// subscribers:             make(map[*subscriber]struct{}),
		// publishLimiter:          rate.NewLimiter(rate.Every(time.Millisecond*100), 8),
	}

	// ss.serveMux.Handle("/", http.FileServer(http.Dir(".")))
	fs, _ := fs.Sub(static, sc.PathToWebUI())
	ss.serveMux.Handle("/", http.FileServer(http.FS(fs)))

	ss.serveMux.HandleFunc("/data", ss.getPageHandler)

	return ss
}

func (ss *SignalServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ss.serveMux.ServeHTTP(w, r)
}

func (ss *SignalServer) Run() {
	fmt.Println("path to config:", ss.SC.PathToConfig())
	fmt.Println("Visit http://0.0.0.0" + ss.SC.Port())

	defer ss.DB.Close()
	ss.initBuckets()

	s := &http.Server{
		Addr:           ss.SC.Port(),
		Handler:        ss,
		ReadTimeout:    time.Second * 10,
		WriteTimeout:   time.Second * 10,
		MaxHeaderBytes: 1 << 20,
	}

	errc := make(chan error, 1)
	go func() {
		errc <- s.ListenAndServe()
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	select {
	case err := <-errc:
		log.Printf("failed to serve: %v", err)
	case sig := <-sigs:
		log.Printf("terminating: %v", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

}

func (ss *SignalServer) getPageHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: pagination
	ss.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("records"))
		c := b.Cursor()
		pageResults := []Record{}
		for k, v := c.First(); k != nil; k, v = c.Next() {
			record := Record{}
			err := json.Unmarshal(v, &record)
			if err != nil {
				fmt.Println(err)
			} else {
				pageResults = append(pageResults, record)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(pageResults)

		return nil
	})
}