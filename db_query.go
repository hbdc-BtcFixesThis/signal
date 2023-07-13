package main

import (
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

type Query struct {
	Bucket []byte
	KV     []Pair
}

type PageQuery struct {
	Query
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
