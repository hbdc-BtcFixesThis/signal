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
	ss.infoLog.Printf("Serving embedded file directory: %s", ss.sc.UiDir(nil))

	// authenticated apis (settings)
	ss.serveMux.Handle("/verify/token", jw(ck(http.HandlerFunc(ss.verifyHandler))))

	// public
	ss.serveMux.Handle("/new/record", jw(http.HandlerFunc(ss.newRecordAndOrSignal)))
	ss.serveMux.Handle("/new/signal", jw(http.HandlerFunc(ss.newRecordAndOrSignal)))
	ss.serveMux.Handle("/get/page", jw(http.HandlerFunc(ss.getPage)))
	ss.serveMux.Handle("/record/value", jw(http.HandlerFunc(ss.getRecordValue)))
	ss.serveMux.Handle("/record/signals", jw(http.HandlerFunc(ss.getRecordSignals)))
	ss.serveMux.Handle("/message/template", jw(http.HandlerFunc(ss.getMessageTemplate)))
}

func newSignalServer() (*SignalServer, error) {
	errLog := log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	infoLog.Printf("Opening conf db; filepath: %s", ServerConfFullPath.Default())
	db := MustOpenAndWrapDB(ServerConfFullPath.Default(), errLog, infoLog)
	infoLog.Printf("Opening signal db; filepath: %s", SignalDataDBFullPath.Default())
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
