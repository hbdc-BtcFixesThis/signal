package main

import (
	"fmt"
	"net/http"
)

func (ss *SignalServer) getTlsCrtFname(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", ss.sc.TlsCrtFname(nil))
}

func (ss *SignalServer) getTlsKeyFname(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", ss.sc.TlsKeyFname(nil))
}

func (ss *SignalServer) getTlsHosts(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", ss.sc.TlsHosts(nil))
}
