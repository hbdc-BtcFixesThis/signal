package main

import (
	"context"
	"embed"
	"log"
	"os"
	"sync"
	"time"

	"crypto/tls"
	"io/fs"
	"net/http"
	"os/signal"
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

	serveMux http.ServeMux
	sc       *ServerConf
	nodes    map[string]struct {
		node *Node
		conf *NodeConf
	}

	sync.RWMutex
}

func (ss *SignalServer) SetNode(id string, n *Node, c *NodeConf) {
	ss.Lock()
	ss.nodes[id] = struct {
		node *Node
		conf *NodeConf
	}{node: n, conf: c}
	ss.Unlock()
}

// func (ss *SignalServer)

func (ss *SignalServer) setHandlers() {
	fs, _ := fs.Sub(static, string(ss.sc.UiDir()))
	ck := ss.CheckAPIKey
	jw := JSONResponseHeadersWrapper

	ss.serveMux.Handle("/", http.FileServer(http.FS(fs)))

	// authenticated apis (settings)
	ss.serveMux.Handle("/verify/token", jw(ck(http.HandlerFunc(ss.verifyHandler))))

	// public
	// ss.serveMux.Handle("/data", jw(http.HandlerFunc(ss.getPageHandler)))
}

func newSignalServer() (*SignalServer, error) {
	db := &DB{MustOpenDB(ServerConfFullPath.Default())}
	sc := &ServerConf{db}
	sc.CreateBucket(sc.ConfBucketName())
	defer sc.Close()

	ss := &SignalServer{logf: log.Printf}
	defer ss.closeNodeDBs()

	if nIDs, err := sc.NodeIds(); err != nil {
		return nil, err
	} else {
		ss.sc = sc
		ss.nodes = make(map[string]struct {
			node *Node
			conf *NodeConf
		})
		for _, id := range nIDs {
			nc := &NodeConf{db}
			ss.SetNode(id, &Node{
				MustOpenAndWrapDB(ByteSlice2String(nc.DataPath(String2ByteSlice(id)))),
			}, nc)
		}
	}

	ss.setHandlers()
	return ss, nil
}

func (ss *SignalServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ss.serveMux.ServeHTTP(w, r)
}

func (ss *SignalServer) Respond(w http.ResponseWriter, r Response) {
	if err := JSONResponse(w, r); err != nil {
		status := http.StatusInternalServerError
		JSONResponse(w, Response{StatusCode: status, Err: err.Error()})
		return
	}
}

func (ss *SignalServer) closeNodeDB(id string) {
	ss.RLock()
	defer ss.RUnlock()
	ss.nodes[id].node.Close()
	ss.nodes[id].conf.Close()
}

func (ss *SignalServer) closeNodeDBs() {
	ss.RLock()
	for id, _ := range ss.nodes {
		ss.RUnlock()
		ss.closeNodeDB(id)
		ss.RLock()
	}
	ss.RUnlock()
}

func (ss *SignalServer) Run() {
	ss.logf("\nVisit http://0.0.0.0%s\n", ss.sc.Port())

	s := &http.Server{
		Addr:           ByteSlice2String(ss.sc.Port()),
		Handler:        ss,
		ReadTimeout:    time.Second * 10,
		WriteTimeout:   time.Second * 10,
		MaxHeaderBytes: 1 << 20,
		TLSConfig:      ss.sc.TLSConf(),
		TLSNextProto:   make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}

	errc := make(chan error, 1)
	go func() {
		errc <- s.ListenAndServeTLS("", "")
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	select {
	case err := <-errc:
		ss.logf("failed to serve: %v", err)
	case sig := <-sigs:
		ss.logf("terminating: %v", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}

/* func (ss *SignalServer) getPageHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: pagination
	ss.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("records"))
		c := b.Cursor()
		pageResults := []Record{}
		for k, v := c.First(); k != nil; k, v = c.Next() {
			record := Record{}
			err := json.Unmarshal(v, &record)
			if err != nil {
				s.logf(err)
			} else {
				pageResults = append(pageResults, record)
			}
		}

		// w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(pageResults)

		return nil
	})
}*/
