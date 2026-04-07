package main

import (
	"fmt"
	"net/http"
)

func (ss *SignalServer) getNodeMaxStorageSize(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", ss.nc.MaxStorageSize(nil))
}
