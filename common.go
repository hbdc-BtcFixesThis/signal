package main

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	StatusCode int         `json:"status_code"`
	Err        string      `json:"error,omitempty"`
	Msg        string      `json:"message,omitempty"`
	Data       interface{} `json:"data,omitempty"`
}

func JSONResponse(w http.ResponseWriter, r Response) error {
	w.WriteHeader(r.StatusCode)
	return json.NewEncoder(w).Encode(r)
}

func JSONResponseHeadersWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		h.ServeHTTP(w, r)
	})
}
