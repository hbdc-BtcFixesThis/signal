package main

import (
	"sync"
	"time"

	"crypto/tls"

	bolt "go.etcd.io/bbolt"
)

// maybe a bit unorthodox but for be able to
// update settings via the ui a seperate db
// that lives in the same dir as the server
// will be used to access, update, and persist
// settings set by users. This also makes it
// easier to update settings in real time
type ServerConf struct {
	DB       *bolt.DB
	settings *SignalSettings
	sync.RWMutex
}

var defaultPathToServerConf = "signal_conf.db"

func NewServerConf() (*ServerConf, error) {
	// gen certs if none are found (updates to cert files get picked up automatically)
	db, err := bolt.Open(defaultPathToServerConf, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	settings, err := LoadSettings(db)
	if err != nil {
		return nil, err
	}
	return &ServerConf{settings: settings, DB: db}, nil
}

func (sc *ServerConf) Port() string { return sc.settings.Port }

func (sc *ServerConf) PathToWebUI() string { return sc.settings.UiDir }

func (sc *ServerConf) PathToDB() string { return sc.settings.DbPath }

func (sc *ServerConf) genNewCertsIfNotFound() error {
	defaultCrtPath := sc.settings.TlsCrtPath
	defaultKeyPath := sc.settings.TlsKeyPath
	err := CheckTLSKeyCertPath(defaultCrtPath, defaultKeyPath)
	if err != nil {
		err = GenerateTLSKeyCert(defaultCrtPath, defaultKeyPath, "0.0.0.0,localhost,127.0.0.1,::1")
		if err != nil {
			return err
		}
	}
	return nil
}

func (sc ServerConf) TLSConf() *tls.Config {
	return &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
		GetCertificate: func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
			// Always get latest signal.crt and signal.key
			if err := sc.genNewCertsIfNotFound(); err != nil {
				return nil, err
			}
			cert, err := tls.LoadX509KeyPair(sc.settings.TlsCrtPath, sc.settings.TlsKeyPath)
			if err != nil {
				return nil, err
			}
			return &cert, nil
		},
	}
}
