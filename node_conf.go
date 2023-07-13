package main

import (
	"encoding/json"
)

type NodeConf struct {
	*DBWithCache
	// bucket []byte
}

func (nc *NodeConf) getOrSet(bucket []byte, k NodeConfKey) []byte {
	q := &Query{
		Bucket: bucket,
		KV:     []Pair{NewPair(k.Bytes(), k.DefaultBytes())},
	}
	nc.MustDo(nc.GetOrPut, q)
	return q.KV[0].Val
}

func (nc *NodeConf) DataPath(bucket []byte) []byte       { return nc.getOrSet(bucket, Path) }
func (nc *NodeConf) Name(bucket []byte) []byte           { return nc.getOrSet(bucket, Name) }
func (nc *NodeConf) Type(bucket []byte) []byte           { return nc.getOrSet(bucket, Type) }
func (nc *NodeConf) Peers(bucket []byte) []byte          { return nc.getOrSet(bucket, Peers) }
func (nc *NodeConf) MaxRecordSize(bucket []byte) []byte  { return nc.getOrSet(bucket, MaxRecordSize) }
func (nc *NodeConf) MaxStorageSize(bucket []byte) []byte { return nc.getOrSet(bucket, MaxStorageSize) }

func (nc *NodeConf) NodeType(bucket []byte) NodeType {
	return NodeTypeFromString(string(nc.Type(bucket)))
}

func (nc *NodeConf) PeersSlice(bucket []byte) ([]string, error) {
	var p []string
	return p, json.Unmarshal(nc.Peers(bucket), &p)
}
