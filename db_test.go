package main

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"
)

func TestGetAndPut(t *testing.T) {
	open()
	td := tData(100)

	for i, pair := range td {
		rVal := MustGenRandBytes(10)
		q := &Query{
			Bucket: []byte("bucket"),
			KV: []Pair{Pair{
				Key: pair.Key,
				// start with diff value
				Val: MustGenRandBytes(10),
			}},
		}
		// check get with default val set to another set of rand bytes
		TestDB.MustDo(TestDB.Get, q)
		if bytes.Equal(q.KV[0].Val, rVal) {
			t.Fatalf("Get returned unexpected value %v", q.KV)
		}

		// check get with original val still in pair
		q.KV = []Pair{Pair{Key: pair.Key, Val: pair.Val}}
		TestDB.MustDo(TestDB.Put, q)
		TestDB.MustDo(TestDB.Get, q)
		if !bytes.Equal(q.KV[0].Val, pair.Val) {
			t.Fatalf("Get returned unexpected value %v", q.KV)
		}

		// change default passed in as the default return val
		// and check get still returns val stored from put
		q.KV[0].Val = rVal
		TestDB.MustDo(TestDB.Get, q)
		if !bytes.Equal(q.KV[0].Val, String2ByteSlice(strconv.Itoa(i))) {
			t.Fatalf("Get returned unexpected value %v", q.KV)
		}
		if !bytes.Equal(q.KV[0].Key, String2ByteSlice(strconv.Itoa(i))) {
			t.Fatalf("Get returned unexpected key %v", q.KV)
		}
	}

	if (&PageQuery{}).Size() != 10 {
		t.Fatalf("Page query returned unexpected size: %v", (&PageQuery{}).Size())
	}

	pqCap := 5
	pq := &PageQuery{
		Query: Query{
			Bucket: []byte("bucket"),
			KV:     make([]Pair, pqCap),
		},
	}
	TestDB.GetPage(pq)
	if pqCap != pq.Size() {
		t.Fatalf("Page query returned unexpected size: %v", pq.Size())
	}
	if pqCap != len(pq.KV) {
		t.Fatalf("Page query returned unexpected size: %v", pq.Size())
	}
	for i, kv := range pq.KV {
		if !bytes.Equal(kv.Key, td[len(td)-1-i].Key) {
			t.Fatalf("Get returned unexpected key %v", kv)
		}
		if !bytes.Equal(kv.Val, td[len(td)-1-i].Val) {
			t.Fatalf("Get returned unexpected Val %v", kv)
		}
	}

	pq.Ascending = true
	pq2 := &PageQuery{Query: Query{Bucket: []byte("bucket"), KV: make([]Pair, 100)}}
	TestDB.GetPage(pq)
	TestDB.GetPage(pq2)
	fmt.Printf("gpq.KV: %v", pq2.KV)
	for i, kv := range pq.KV {
		if !bytes.Equal(kv.Key, td[i].Key) {
			t.Fatalf("Get returned unexpected key %v", kv)
		}
		if !bytes.Equal(kv.Val, td[i].Val) {
			t.Fatalf("Get returned unexpected Val %v", kv)
		}
	}
}
