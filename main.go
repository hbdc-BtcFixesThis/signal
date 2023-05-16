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

type Record struct {
	Name      string `json:"name"`
	Type      int    `json:"type"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
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

func (ss *SignalServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ss.serveMux.ServeHTTP(w, r)
}

func main() {
	// http.HandleFunc("/", handlers.RootHandler)
	// http.HandleFunc("/chat", b.ChatHandler)

	log.SetFlags(0)

	sc, err := NewServerConfig()
	if err != nil {
		log.Fatal("ERROR:", err)
	}

	ss := newSignalServer(sc)
	defer ss.DB.Close()
	ss.initBuckets()

	s := &http.Server{
		Addr:           sc.Port(),
		Handler:        ss,
		ReadTimeout:    time.Second * 10,
		WriteTimeout:   time.Second * 10,
		MaxHeaderBytes: 1 << 20,
	}

	// http.Handle("/", http.FileServer(http.FS(fs)))

	fmt.Println("path to config:", sc.PathToConfig())
	fmt.Println("Visit http://0.0.0.0" + sc.Port())

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

	if s.Shutdown(ctx) != nil {
		log.Fatal(err)
	}
}
