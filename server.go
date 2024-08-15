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
	errorLog *log.Logger
	infoLog  *log.Logger
	serveMux http.ServeMux
	buckets  *SignalBuckets
	sc       *ServerConf
	nc       *NodeConf

	sync.RWMutex
}

type SignalBuckets struct {
	Record  *RecordBucket
	Value   *ValueBucket
	Signal  *SignalBucket
	Address *AddressBucket
	Rank    *RankBucket
	db      *DB
}

func (ss *SignalServer) setHandlers() {
	fs, _ := fs.Sub(static, string(ss.sc.UiDir(nil)))
	ck := ss.CheckAPIKey
	jw := JSONResponseHeadersWrapper

	ss.serveMux.Handle("/", http.FileServer(http.FS(fs)))

	// authenticated apis (settings)
	ss.serveMux.Handle("/verify/token", jw(ck(http.HandlerFunc(ss.verifyHandler))))

	// public
	ss.serveMux.Handle("/new/record", jw(http.HandlerFunc(ss.newRecord)))
	ss.serveMux.Handle("/get/page", jw(http.HandlerFunc(ss.getPage)))
	ss.serveMux.Handle("/message/template", jw(http.HandlerFunc(ss.getMessageTemplate)))
}

func newSignalServer() (*SignalServer, error) {
	errLog := log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	db := MustOpenAndWrapDB(ServerConfFullPath.Default(), errLog, infoLog)
	sdb := MustOpenAndWrapDB(SignalDataDBFullPath.Default(), errLog, infoLog)
	dbwc := &DBWithCache{
		cache: make(map[string][]byte),
		DB:    db,
	}
	sc := &ServerConf{dbwc}
	nc := &NodeConf{dbwc}
	sb := &SignalBuckets{
		Record:  &RecordBucket{sdb},
		Value:   &ValueBucket{sdb},
		Signal:  &SignalBucket{sdb},
		Address: &AddressBucket{sdb},
		Rank:    &RankBucket{sdb},
		db:      sdb,
	}

	// create if not exists
	sc.CreateBucket(&Query{Bucket: sc.ConfBucketName()})
	nc.CreateBucket(&Query{Bucket: nc.ConfBucketName()})
	log.SetFlags(log.Lshortfile)
	ss := &SignalServer{
		sc:       sc,
		nc:       nc,
		buckets:  sb,
		infoLog:  infoLog,
		errorLog: errLog,
	}

	// should start bg goroutine to listen for close signal
	// defer ss.closeNodeDBs()
	// sc.Close()

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

/*
func (ss *SignalServer) closeNodeDB(id string) {
	ss.RLock()
	defer ss.RUnlock()
	ss.nodes[id].node.Close()
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
*/

func (ss *SignalServer) Run() {
	ss.infoLog.Printf("\nVisit https://0.0.0.0%s\n", ss.sc.Port(nil))

	s := &http.Server{
		Addr:           ByteSlice2String(ss.sc.Port(nil)),
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
		ss.errorLog.Printf("failed to serve: %v", err)
	case sig := <-sigs:
		ss.infoLog.Printf("terminating: %v", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		ss.errorLog.Fatal(err)
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

/*
struct Record type {

}

	struct AddrMeta type {
		records []Record
		satsHeld uint
		satsAllocated uint
	}

func (ss *SignalServer) lookupBtcAddress(addr string) (){
}

	func (ss *SignalServer) addNode(name string) ([]byte, error) {
		if id, err := ss.GenID(name); err != nil {
			return nil, err
		}
		q := &Query{
			Bucket:                  NodeLookupBucket,
			q.KV:                    []Pair{NewPair(id, name)},
			CreateBucketIfNotExists: true,
		}

		if err := ss.nc.Put(q); err != nil {
			return nil, err
		}

		// now create a bucket for this new node
		// this bucket will be used as a node conf to store
		// the settings of the new node being created and
		// the nodeconf struct will be used to interact with
		// this bucket going forward
		q = &Query{
			Bucket: id,
			KV: []Pair{
				NewPair(Name, name),
				NewPair(NodeID, id),
			},
			CreateBucketIfNotExists: true,
		}
		if err := ss.nc.Put(q); err != nil {
			return nil, err
		}
		ss.nodes[name] = Node{}

		// put returns err if one
		return id, nil
	}

	func (ss *SignalServer) ListNodes(startFrom KV) { //q *PageQuery) {
		q := &Query{
			Bucket:                  NodeLookupBucket,
			CreateBucketIfNotExists: true,
		}
		pq := &PageQuery{
			Query:     q,
			KV:        make([]Pair, 10), // arbitrary 10 items per page
			Ascending: true,
		}
		if startFrom != nil {
			pq.StartFrom = startFrom
		}
		ss.nc.GetPage(pq)
	}

	func (ss *SignalServer) CreateNode(id string, n *Node) {
		ss.Lock()
		ss.nodes[id] = struct{ node *Node }{node: n}
		ss.Unlock()
	}

// func (ss *SignalServer)
*/
