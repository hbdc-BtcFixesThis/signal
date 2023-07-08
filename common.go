package main

import (
	"errors"
	"os"

	"encoding/json"
	"net/http"
	"path/filepath"
)

const (
	SignalDirName = ".signal"
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

func SignalHomeDir() string {
	path, err := GetCurrentUserHomeDir()
	if err != nil {
		return path
	}
	signalHomePath := filepath.Join(path, SignalDirName)
	if _, err := os.Stat(signalHomePath); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(signalHomePath, os.ModePerm)
		if err != nil {
			return path
		}
	}
	return signalHomePath
}
