package main

import (
	"encoding/json"
)

const ValueBucketName = "Value"

type RecordValue struct {
	RecID []byte `json:"rid,omitempty"`
	Value string `json:"value"`
}

type ValueBucket struct{ *DB }

func (v *ValueBucket) Name() []byte { return []byte(ValueBucketName) }

func (v *ValueBucket) GetId(id []byte) ([]byte, error) {
	query := &Query{
		Bucket:                  v.Name(),
		KV:                      []Pair{NewPair(id, nil)},
		CreateBucketIfNotExists: true,
	}
	err := v.Get(query)
	return query.KV[0].Val, err
}

func (v *ValueBucket) GetRecordById(id []byte) (RecordValue, error) {
	var recV RecordValue

	rvBytes, err := v.GetId(id)
	if err != nil {
		v.errorLog.Println(err)
		return recV, err
	}
	if err := json.Unmarshal(rvBytes, &recV); err != nil {
		v.errorLog.Println(err)
		return recV, err
	}

	return recV, nil
}

func (v *ValueBucket) PutRecV(recV RecordValue) (*Query, error) {
	recID := recV.RecID
	recV.RecID = []byte{}
	b, err := json.Marshal(recV)
	if err != nil {
		v.errorLog.Println(err)
		return &Query{}, err
	}
	q := v.PutRecB(recID, b)
	recV.RecID = recID
	return q, nil
}

func (v *ValueBucket) PutRecB(key, val []byte) *Query {
	return &Query{
		Bucket:                  v.Name(),
		KV:                      []Pair{NewPair(key, val)},
		CreateBucketIfNotExists: true,
	}
}

func (v *ValueBucket) DeleteV(id []byte) *Query {
	return &Query{
		Bucket:                  v.Name(),
		KV:                      []Pair{NewPair(id, nil)},
		CreateBucketIfNotExists: true,
	}
}
