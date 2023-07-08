package main

import (
	"fmt"
	"path/filepath"
)

// ///////////////////////////ServerConf///////////////////////
// Probably wont have more then
// 255 keys but easy to update
// if at some point it will
type ServerConfKey uint8

const (
	ServerConfFullPath ServerConfKey = iota
	ServerConfBucket
	ServerConfFileName
	Admin
	Port
	UiDir
	PassHash
	TlsCrtFname
	TlsKeyFname
	DefaultNode
	NodeIds
	TlsHosts
)

func (sck ServerConfKey) Keys() []string {
	return []string{
		"ServerConfFullPath",
		"ServerConfBucket",
		"ServerConfFileName",
		"Admin",
		"Port",
		"UiDir",
		"PassHash",
		"TlsCrtFname",
		"TlsKeyFname",
		"DefaultNode",
		"NodeIds",
		"TlsHosts",
	}
}

func (sck ServerConfKey) Defaults() []string {
	return []string{
		filepath.Join(SignalHomeDir(), "conf.db"),
		"conf",
		"conf.db",
		"admin",
		":8888",
		"static",
		SHA256([]byte("pass")),
		"signal.crt.pem",
		"signal.key.pem",
		SIGNAL.String(),
		fmt.Sprintf(`["%s"]`, SIGNAL.String()),
		"0.0.0.0,localhost,127.0.0.1,::1",
	}
}

func (sck ServerConfKey) String() string       { return sck.Keys()[sck] }
func (sck ServerConfKey) Bytes() []byte        { return []byte(sck.Keys()[sck]) }
func (sck ServerConfKey) Default() string      { return sck.Defaults()[sck] }
func (sck ServerConfKey) DefaultBytes() []byte { return []byte(sck.Default()) }

// func (sck ServerConfKey) Bytes() []byte        { return String2ByteSlice(sck.Keys()[sck]) }
// func (sck ServerConfKey) DefaultBytes() []byte { return String2ByteSlice(sck.Default()) }
