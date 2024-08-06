package main

import (
	"path/filepath"
)

// /////////////////////////////////////////////////////////////
// ///////////////////////////NodeConf//////////////////////////
type NodeConfKey uint8

const (
	Name NodeConfKey = iota
	Type
	Path
	Peers
	MaxRecordSize
	MaxStorageSize
	NodeConfBucket
	NodeID
	// RequiresPayment ?
	// Users ?
	// UiDir ?
)

func (nck NodeConfKey) Keys() []string {
	return []string{
		"name", "type",
		"path", "peers",
		"max_record_size",
		"max_storage_size",
		"NodeConfBucket",
		"ID",
	}
}

func (nck NodeConfKey) Defaults() []string {
	return []string{
		SIGNAL.String(),
		SIGNAL.String(),
		filepath.Join(SignalHomeDir(), SIGNAL.String()+".db"),
		"[]",
		"1000000000",     // 1gb max record default (change via api; will expose in ui)
		"10000000000",    // 10gb max record storage default
		"NodeConfBucket", // the name of the bucket thats used to lookup existing nodes
		"id",
	}
}

func (nck NodeConfKey) String() string       { return nck.Keys()[nck] }
func (nck NodeConfKey) Bytes() []byte        { return []byte(nck.Keys()[nck]) }
func (nck NodeConfKey) Default() string      { return nck.Defaults()[nck] }
func (nck NodeConfKey) DefaultBytes() []byte { return []byte(nck.Default()) }

// func (nck NodeConfKey) Bytes() string        { return String2ByteSlice(nck.Keys()[nck]) }
// func (nck NodeConfKey) DefaultBytes() []byte { return String2ByteSlice(nck.Default()) }
