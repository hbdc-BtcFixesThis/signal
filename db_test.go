package main

import (
	"bytes"
	"strconv"
	"testing"
)

func TestGetDefault(t *testing.T) {
	q := tDataQuery(1, false)
	// &Query{Bucket: []byte("bucket"), KV: tData(1, false)}
	TestDB.MustDo(TestDB.Get, q)
	if !bytes.Equal(q.KV[0].Val, String2ByteSlice(strconv.Itoa(0))) {
		t.Fatalf("Get returned unexpected value %v", q.KV)
	}
	t2 := []byte("another default")
	if bytes.Equal(q.KV[0].Val, t2) {
		t.Fatalf("What in the fuck? q.KV: %v", q.KV)
	}
	q.KV[0].Val = t2
	TestDB.MustDo(TestDB.Get, q)
	if bytes.Equal(q.KV[0].Val, String2ByteSlice(strconv.Itoa(0))) {
		t.Fatalf("Get returned unexpected value %v", q.KV)
	}
}

func TestPutOne(t *testing.T) {
	qAftPut := tDataQuery(1, true)
	qNoPut := tDataQuery(1, false)
	qAftPut.IsEqual(qNoPut)
	if !qAftPut.IsEqual(qNoPut) {
		t.Fatalf("qAftPut != qNoPut wher qAftPut = %v, qNoPut = %s", qAftPut.KV, qNoPut)
	}

	// get random data with the same key
	qAftPut = tDataQueryRandVals(1, false)
	if bytes.Equal(qAftPut.KV[0].Val, F64tb(float64(0))) {
		t.Fatalf("This is supposed to be random %v", qAftPut.KV)
	}

	// retrieve and check to make sure everything worked
	TestDB.MustDo(TestDB.Get, qAftPut)
	if !bytes.Equal(qAftPut.KV[0].Key, F64tb(float64(0))) {
		t.Fatalf("Get returned unexpected key %v", qAftPut.KV)
	}
	if !bytes.Equal(qAftPut.KV[0].Val, String2ByteSlice(strconv.Itoa(0))) {
		t.Fatalf("Get returned unexpected value %v", qAftPut.KV)
	}
	if !qAftPut.IsEqual(qNoPut) {
		t.Fatalf("Get returned unexpected value %v", qAftPut.KV)
	}

	// cleanup
	TestDB.MustDo(TestDB.Delete, qAftPut)
}

func TestGetOrPutMany(t *testing.T) {
	qAftPut := tDataQuery(10, false)
	TestDB.MustDo(TestDB.GetOrPut, qAftPut)

	qGetWithRandDefault := tDataQueryRandVals(10, false)
	if qAftPut.IsEqual(qGetWithRandDefault) {
		t.Fatalf("qAftPut.IsEqual(qGetWithRandDefault)! %v, %v", qAftPut, qGetWithRandDefault)
	}
	TestDB.MustDo(TestDB.Get, qGetWithRandDefault)
	if !qAftPut.IsEqual(qGetWithRandDefault) {
		t.Fatalf("qAftPut.IsEqual(qGetWithRandDefault)! %v, %v", qAftPut, qGetWithRandDefault)
	}

	// cleanup
	TestDB.MustDo(TestDB.Delete, qAftPut)
}

func TestPageQuerySize(t *testing.T) {
	if (&PageQuery{Query: &Query{}}).Size() != 10 {
		t.Fatalf("Page query returned unexpected size: %v", (&PageQuery{}).Size())
	}

	pqCap := 5
	pq := &PageQuery{
		Query: &Query{
			Bucket: []byte("bucket"),
			KV:     make([]Pair, pqCap),
		},
	}
	if pq.Size() != 5 {
		t.Fatalf("Page query returned unexpected size: %v", pq.Size())
	}
}

func TestGetPage(t *testing.T) {
	q := tDataQuery(100, true)

	// some size other then the default
	q20Aes := tDataQueryRandVals(20, false)
	q20Des := tDataQueryRandVals(20, false)

	for i, _ := range q20Des.KV {
		q20Des.KV[i].Key = q.KV[len(q.KV)-i-1].Key
	}
	if q20Aes.IsEqual(q20Des) {
		t.Fatalf("qAes.IsEqual(q20.DesQuery); where  q20Aes =  %v, q20Des = %v", q20Aes, q20Des)
	}

	pq20Des := &PageQuery{Query: q20Des}
	TestDB.GetPage(pq20Des)
	for i, kv := range pq20Des.KV {
		if !kv.IsEqual(&q.KV[len(q.KV)-1-i]) {
			t.Fatalf("Get page returned unexpected key %v %v", kv, q.KV[len(q.KV)-1-i])
		}
	}

	pq20Aes := &PageQuery{Query: q20Aes, Ascending: true}
	TestDB.GetPage(pq20Aes)
	for i, kv := range pq20Aes.KV {
		if !kv.IsEqual(&q.KV[i]) {
			t.Fatalf("Get page returned unexpected key %v", kv)
		}
	}

	// cleanup
	TestDB.MustDo(TestDB.Delete, q)
}
