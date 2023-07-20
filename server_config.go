package main

import (
	"crypto/tls"
	"encoding/json"
	fp "path/filepath"
)

// Maybe a bit unorthodox but for the ability to update settings
// via the ui, a seperate db is used to access, update, and persist
// settings set by users. This also makes it
// easy to update settings in real time and
// allows users to manage multiple data sets
// that are inherently decoupled; yet managed by
// a single running server instance. It also
// allows for other data to be stored wherever
// users would like it to live because the
// location of
type ServerConf struct {
	*DBWithCache
}

func (sc *ServerConf) gop(k ServerConfKey, v []byte) []byte {
	return sc.GetOrPut(sc.ConfBucketName(), k, v)
}

func (sc ServerConf) ConfBucketName() []byte { return ServerConfBucket.DefaultBytes() }
func (sc ServerConf) ConfFname() string      { return ServerConfFileName.Default() }
func (sc ServerConf) FullPath() string       { return ServerConfFullPath.Default() }
func (sc *ServerConf) Port() []byte          { return sc.gop(Port, nil) }
func (sc *ServerConf) UiDir() []byte         { return sc.gop(UiDir, nil) }
func (sc *ServerConf) Admin() []byte         { return sc.gop(Admin, nil) }
func (sc *ServerConf) nodeIds() []byte       { return sc.gop(NodeIds, nil) }
func (sc *ServerConf) PassHash() []byte      { return sc.gop(PassHash, nil) }
func (sc *ServerConf) DefaultNode() []byte   { return sc.gop(DefaultNode, nil) }
func (sc *ServerConf) TlsCrtFname() []byte   { return sc.gop(TlsCrtFname, nil) }
func (sc *ServerConf) TlsKeyFname() []byte   { return sc.gop(TlsKeyFname, nil) }
func (sc *ServerConf) TlsHosts() []byte      { return sc.gop(TlsHosts, nil) }

func (sc *ServerConf) TlsCrtPath() string {
	return fp.Join(SignalHomeDir(), ByteSlice2String(sc.TlsCrtFname()))
}
func (sc *ServerConf) TlsKeyPath() string {
	return fp.Join(SignalHomeDir(), ByteSlice2String(sc.TlsKeyFname()))
}

func (sc *ServerConf) NodeIds() ([]string, error) {
	var nn []string
	names := sc.nodeIds()
	return nn, json.Unmarshal(names, &nn)
}

func (sc *ServerConf) GenNewRandAdminPassHash() []byte {
	ph := MustGenNewAdminPW(NewPwMsg, ByteSlice2String(sc.Port()))
	return sc.gop(PassHash, String2ByteSlice(ph))
}

func (sc *ServerConf) genNewCertsIfNotFound() (string, string, error) {
	// gen certs if none are found (updates to cert files get picked up automatically)
	crtPath := sc.TlsCrtPath()
	keyPath := sc.TlsKeyPath()
	hosts := ByteSlice2String(sc.TlsHosts())

	if err := CheckTLSKeyCertPath(crtPath, keyPath); err != nil {
		if err = GenerateTLSKeyCert(crtPath, keyPath, hosts); err != nil {
			return "", "", err
		}
	}

	return crtPath, keyPath, nil
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
			// Always get latest crt/key pair
			if crtPath, keyPath, err := sc.genNewCertsIfNotFound(); err != nil {
				return nil, err
			} else if cert, err := tls.LoadX509KeyPair(crtPath, keyPath); err != nil {
				return nil, err
			} else {
				return &cert, nil
			}
		},
	}
}
