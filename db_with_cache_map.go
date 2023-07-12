package main

import (
	"bytes"
	"sync"
)

type DBWithCache struct {
	cache map[string][]byte

	*DB
	sync.RWMutex
}

func cacheKey(k, b []byte) string {
	return ByteSlice2String(bytes.Join([][]byte{k, b}, []byte("::")))
}

func (dbc *DBWithCache) udpateCache(k, v, b []byte) {
	dbc.Lock()
	defer dbc.Unlock()
	dbc.cache[cacheKey(k, b)] = v
}

func (dbc *DBWithCache) getOrSet(k, v, b []byte) []byte {
	dbc.RLock()
	if val, found := dbc.cache[cacheKey(k, b)]; found {
		return val
	}
	dbc.RUnlock()

	q := &Query{Bucket: b, KV: []Pair{NewPair(k, v)}}
	dbc.MustDo(dbc.GetOrPut, q)

	dbc.udpateCache(k, q.KV[0].Val, b)
	return q.KV[0].Val
}
