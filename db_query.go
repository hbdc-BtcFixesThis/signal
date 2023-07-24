package main

import (
	"bytes"
	"time"

	bolt "go.etcd.io/bbolt"
)

const (
	DbTimeout = 2 * time.Second
)

type KV []byte

func (kv KV) String() string { return ByteSlice2String(kv) }

type Pair struct {
	Key KV
	Val KV
}

func NewPair(k, v []byte) Pair {
	return Pair{Key: KV(k), Val: KV(v)}
}

func (p *Pair) IsEqual(p2 *Pair) bool {
	if !bytes.Equal(p.Key, p2.Key) {
		return false
	}
	if !bytes.Equal(p.Val, p2.Val) {
		return false
	}
	return true
}

type Query struct {
	Bucket []byte
	KV     []Pair

	CreateBucketIfNotExists bool
}

func (q *Query) IsEqual(q2 *Query) bool {
	if !bytes.Equal(q.Bucket, q2.Bucket) {
		return false
	}
	if len(q.KV) != len(q2.KV) {
		return false
	}
	for i, kv := range q.KV {
		if !kv.IsEqual(&q2.KV[i]) {
			return false
		}
	}
	return true
}

type PageQuery struct {
	*Query
	StartFrom KV
	Ascending bool
}

func (q *PageQuery) SeekFrom(c *bolt.Cursor) Pair {
	// sets the cursor position and returns the first k/v
	switch q.StartFrom.String() {
	case "first":
		return NewPair(c.First())
	case "last":
		return NewPair(c.Last())
	case "":
		if q.Ascending {
			return NewPair(c.First())
		}
		return NewPair(c.Last())
	default:
		k, v := c.Seek(q.StartFrom)
		if k == nil {
			return NewPair(c.Last())
		}
		return NewPair(k, v)
	}
}

type nextFunc func() ([]byte, []byte)

func (q *PageQuery) Direction(c *bolt.Cursor) nextFunc {
	if q.Ascending {
		return c.Next
	}
	return c.Prev
}

func (q *PageQuery) Size() int {
	if c := cap(q.KV); c != 0 {
		return c
	}
	return 10
}
