package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	bolt "go.etcd.io/bbolt"
)

const (
	DbTimeout = 2 * time.Second
)

type KV []byte

func (kv KV) String() string { return ByteSlice2String(kv) }

func (kv KV) MarshalJSON() ([]byte, error) {
	return String2ByteSlice(fmt.Sprintf(`"%s"`, kv.String())), nil
}

func (kv *KV) UnmarshalJSON(b []byte) error {
	var data string
	fmt.Println("----------------------------------")
	fmt.Println("UnmarshalJSON data: ", data)
	fmt.Println("UnmarshalJSON data: ", string(b))
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	fmt.Println("UnmarshalJSON data: ", data)
	fmt.Println("----------------------------------")
	*kv = String2ByteSlice(data)
	return nil
}

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
		if k, _ := c.Seek(q.StartFrom); k == nil {
			return NewPair(c.Last())
		}
		// return new pair after calling next since a valid
		// key was sent and a value for it was found. That
		// indicates the client already having the StartFrom
		return NewPair(q.Direction(c)())
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
