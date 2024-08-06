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

func (dbc *DBWithCache) getCacheVal(k string) []byte {
	dbc.RLock()
	defer dbc.RUnlock()
	if v, found := dbc.cache[k]; found {
		return v
	}
	return nil
}

func (dbc *DBWithCache) getOrSet(b, k, v []byte) []byte {
	ck := cacheKey(k, b)
	ckV := dbc.getCacheVal(ck)
	if ckV != nil {
		if v == nil || bytes.Equal(ckV, v) {
			return ckV
		}
	}

	q := &Query{Bucket: b, KV: []Pair{NewPair(k, v)}}
	dbc.MustDo(dbc.DB.GetOrPut, q)
	dbc.udpateCache(k, q.KV[0].Val, b)

	return dbc.getCacheVal(ck)
}

func (dbc *DBWithCache) GetOrPut(bucket, k, v, defaultV []byte) []byte {
	if v == nil {
		return dbc.getOrSet(bucket, k, defaultV)
	}
	return dbc.getOrSet(bucket, k, v)
}
