package main

import (
	"encoding/json"
)

type NodeConf struct {
	*DBWithCache
}

func (nc *NodeConf) DataPath(bucket []byte) []byte { return nc.GetOrPut(bucket, Path, nil) }
func (nc *NodeConf) Name(bucket []byte) []byte     { return nc.GetOrPut(bucket, Name, nil) }
func (nc *NodeConf) Type(bucket []byte) []byte     { return nc.GetOrPut(bucket, Type, nil) }
func (nc *NodeConf) Peers(bucket []byte) []byte    { return nc.GetOrPut(bucket, Peers, nil) }

func (nc *NodeConf) MaxRecordSize(bucket []byte) []byte {
	return nc.GetOrPut(bucket, MaxRecordSize, nil)
}

func (nc *NodeConf) MaxStorageSize(bucket []byte) []byte {
	return nc.GetOrPut(bucket, MaxStorageSize, nil)
}

func (nc *NodeConf) NodeType(bucket []byte) NodeType {
	return NodeTypeFromString(string(nc.Type(bucket)))
}

func (nc *NodeConf) PeersSlice(bucket []byte) ([]string, error) {
	var p []string
	return p, json.Unmarshal(nc.Peers(bucket), &p)
}
