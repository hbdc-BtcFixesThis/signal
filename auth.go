package main

import (
	"time"

	"encoding/json"
	"net/http"
)

func verifyToken(token string, hash string) bool {
	// token create by the client should consist of a hash of
	// the password plus the utc date of the token creation.
	// The token is valid for a day. This seems good enough
	// for now and if needed allowing users to dictate the
	// lifespan of a session should be fairly easy by adding
	// a setting and updating the check here
	return token == SHA256(hash+" "+time.Now().UTC().String()[:10])
}

func checkAPIKey(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(r.URL.Query().Get("key")) == 0 {
			http.Error(w, "missing key", http.StatusUnauthorized)
			return // don't call original handler
		}
		h.ServeHTTP(w, r)
	})
}

func (ss *SignalServer) verifyHandler(w http.ResponseWriter, r *http.Request) {
	// ::TODO:: headers wrapper
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{
		"success": verifyToken(r.URL.Query().Get("key"), ss.SC.PassHash()),
	})
}
