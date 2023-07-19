package main

import (
	"encoding/json"
)

type NodeConf struct {
	*DBWithCache
	// bucket []byte
}

func (nc *NodeConf) getOrPut(bucket []byte, k NodeConfKey, v []byte) []byte {
	if v == nil {
		return nc.getOrSet(k.Bytes(), k.DefaultBytes(), bucket)
	}
	return nc.getOrSet(k.Bytes(), v, bucket)
}

func (nc *NodeConf) DataPath(bucket []byte) []byte { return nc.getOrPut(bucket, Path, nil) }
func (nc *NodeConf) Name(bucket []byte) []byte     { return nc.getOrPut(bucket, Name, nil) }
func (nc *NodeConf) Type(bucket []byte) []byte     { return nc.getOrPut(bucket, Type, nil) }
func (nc *NodeConf) Peers(bucket []byte) []byte    { return nc.getOrPut(bucket, Peers, nil) }

func (nc *NodeConf) MaxRecordSize(bucket []byte) []byte {
	return nc.getOrPut(bucket, MaxRecordSize, nil)
}

func (nc *NodeConf) MaxStorageSize(bucket []byte) []byte {
	return nc.getOrPut(bucket, MaxStorageSize, nil)
}

func (nc *NodeConf) NodeType(bucket []byte) NodeType {
	return NodeTypeFromString(string(nc.Type(bucket)))
}

func (nc *NodeConf) PeersSlice(bucket []byte) ([]string, error) {
	var p []string
	return p, json.Unmarshal(nc.Peers(bucket), &p)
}
