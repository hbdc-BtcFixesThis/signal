package main

import (
	"crypto/tls"
	"fmt"
)

var (
	defaultPort    = ":8888"
	defaultCrtPath = "signal.crt.pem"
	defaultKeyPath = "signal.key.pem"
)

type ServerConfig struct {
	port                string
	uiDir               string
	key                 string
	isPublicNode        bool
	maxStorageSizeBytes uint
}

func NewServerConfig() (*ServerConfig, error) {
	// gen certs if none are found (updates to cert files get picked up automatically)
	return &ServerConfig{
		port:  defaultPort,
		uiDir: "static",
	}, nil
}

func (sc *ServerConfig) Port() string { return sc.port }

func (sc *ServerConfig) PathToWebUI() string { return sc.uiDir }

func genNewCertsIfNotFound() error {
	err := CheckTLSKeyCertPath(defaultCrtPath, defaultKeyPath)
	if err != nil {
		fmt.Println(err)
		err = GenerateTLSKeyCert(defaultCrtPath, defaultKeyPath, "0.0.0.0"+defaultPort)
		if err != nil {
			return err
		}
	}
	return nil
}

func (sc ServerConfig) TLSConf() *tls.Config {
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
			if err := genNewCertsIfNotFound(); err != nil {
				return nil, err
			}
			cert, err := tls.LoadX509KeyPair(defaultCrtPath, defaultKeyPath)
			if err != nil {
				return nil, err
			}
			return &cert, nil
		},
	}
}
