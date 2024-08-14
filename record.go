package main

import (
	"encoding/binary"
	"encoding/json"
)

const RecordBucketName = "Record"

// content of record
// the result of RecordHash below is used as the id for these records
type Record struct {
	Sats      uint64   `json:"sats"`
	Name      string   `json:"name"`
	Value     string   `json:"value,omitempty"` // not stored
	VBytes    uint64   `json:"vbytes"`          // combined
	VHash     string   `json:"vhash"`           // hash of value
	Signals   []Signal `json:"signals,omitempty"`
	SignalIds []KV     `json:"sids"` // ids for signals in SignalBucket
}

type SerializedRecord struct {
	Sats      uint64 `json:"sats"`
	Name      string `json:"name"`
	VBytes    uint64 `json:"vbytes"`
	VHash     string `json:"vhash"`
	SignalIds []KV   `json:"sids"`
}

func (r *Record) ID() []byte { return String2ByteSlice(r.Hash()) }

func (r *Record) Hash() string {
	NHash := SHA256(String2ByteSlice(r.Name))
	return SHA256(String2ByteSlice(NHash + "::" + r.VHash))
}

// --------------------------------------------
// NOTE
// cant rely on value being here used as helper
//
//	helper for incoming records only (otherwise VHash should used)
func (r *Record) vHash() string { return SHA256(String2ByteSlice(r.Value)) }

func (r *Record) vBytes() uint64 {
	if len(r.Value) > 0 {
		nSize := binary.Size(String2ByteSlice(r.Name))
		vSize := binary.Size(String2ByteSlice(r.Value))
		return uint64(nSize + vSize)
	}
	return r.VBytes
}

//  --------------------------------------------

func (r *Record) Rank() float64 {
	return float64(float64(r.TotalSats()) / float64(r.VBytes))
}

func (r *Record) RankB() []byte {
	return F64tb(r.Rank())
}

func (r *Record) RankForSatCount(sats uint64) float64 {
	return float64(float64(sats) / float64(r.VBytes))
}
func (r *Record) RankForSatCountB(sats uint64) []byte {
	return F64tb(r.RankForSatCount(sats))
}

func (r *Record) toSerializedRecord() SerializedRecord {
	return SerializedRecord{
		Sats:      r.Sats,
		Name:      r.Name,
		VBytes:    r.VBytes,
		VHash:     r.VHash,
		SignalIds: r.SignalIds,
	}
}

func (r *Record) toBytes() ([]byte, error) {
	return json.Marshal(r.toSerializedRecord())
}

func (r *Record) TotalSats() uint64 {
	total := uint64(0)
	for i := 0; i < len(r.Signals); i++ {
		total += r.Signals[i].Sats
	}
	return total
}

func (r *Record) RemoveSignal(id KV) {
	for i := len(r.SignalIds) - 1; i >= 0; i-- {
		if r.SignalIds[i].String() == id.String() {
			r.SignalIds = append(r.SignalIds[:i], r.SignalIds[i+1:]...)
			i--
		}
	}
	for i := len(r.Signals) - 1; i >= 0; i-- {
		if r.Signals[i].Hash() == id.String() {
			r.Signals = append(r.Signals[:i], r.Signals[i+1:]...)
			i--
		}
	}
	// recomputes the sat count of record in memory (not whats saved to disk)
	r.Sats = r.TotalSats()
}

func (r *Record) AddSignal(s Signal) {
	dup := false
	for i := len(r.SignalIds) - 1; i >= 0; i-- {
		if r.SignalIds[i].String() == s.Hash() {
			dup = true
		}
	}
	if !dup {
		r.SignalIds = append(r.SignalIds, s.ID())
	}

	dup = false
	for i := len(r.Signals) - 1; i >= 0; i-- {
		if r.Signals[i].Hash() == s.Hash() {
			dup = true
		}
	}
	if !dup {
		r.Signals = append(r.Signals, s)
	}

	// recomputes the sat count of record in memory (not whats saved to disk)
	r.Sats = r.TotalSats()
}

func (r *Record) AddSignals(signals []Signal) {
	for _, signal := range signals {
		r.AddSignal(signal)
	}
}

// key = r.RecordHash(); val = serialize record
type RecordBucket struct{ *DB }

func (r *RecordBucket) Name() []byte { return []byte(RecordBucketName) }

func (r *RecordBucket) GetId(id []byte) ([]byte, error) {
	query := &Query{
		Bucket:                  r.Name(),
		KV:                      []Pair{NewPair(id, nil)},
		CreateBucketIfNotExists: true,
	}
	err := r.Get(query)
	return query.KV[0].Val, err
}

func (r *RecordBucket) GetRecordById(id []byte) (Record, error) {
	var signalRec Record

	rBytes, err := r.GetId(id)
	r.infoLog.Printf("GetRecordById rbytes: %s", rBytes)
	if err != nil {
		r.errorLog.Println(err)
		return signalRec, err
	}
	if len(rBytes) == 0 {
		return Record{}, nil
	}
	err = json.Unmarshal(rBytes, &signalRec)

	return signalRec, err
}

func (r *RecordBucket) GetRecordWithSignalsById(id []byte) (Record, error) {
	record, err := r.GetRecordById(id)
	if err != nil {
		r.errorLog.Println(err)
		return record, err
	}
	signals, err := (&SignalBucket{r.DB}).GetSignalsByIds(record.SignalIds)
	if err != nil {
		r.errorLog.Println(err)
		return record, err
	}
	record.Signals = signals
	return record, nil
}

func (r *RecordBucket) PutRec(rec Record) (*Query, error) {
	b, err := rec.toBytes()
	if err != nil {
		return &Query{}, err
	}
	return r.PutRecB(rec.ID(), b), nil
}

func (r *RecordBucket) PutRecB(key, val []byte) *Query {
	return &Query{
		Bucket:                  r.Name(),
		KV:                      []Pair{NewPair(key, val)},
		CreateBucketIfNotExists: true,
	}
}

func (r *RecordBucket) DeleteRec(id []byte) *Query {
	return &Query{
		Bucket:                  r.Name(),
		KV:                      []Pair{NewPair(id, nil)},
		CreateBucketIfNotExists: true,
	}
}

func (r *RecordBucket) DeleteSignalsFromRecord(id []byte, signalIds [][]byte) (*DataUpdates, uint64, error) {
	updates := &DataUpdates{}
	query := &Query{
		Bucket:                  r.Name(),
		KV:                      []Pair{NewPair(id, nil)},
		CreateBucketIfNotExists: true,
	}
	if err := r.Get(query); err != nil {
		return updates, uint64(0), err
	}

	var sigRec Record
	if err := json.Unmarshal(query.KV[0].Val, sigRec); err != nil {
		return updates, uint64(0), err
	}

	for i := len(sigRec.SignalIds) - 1; i >= 0; i-- {
		for j := 0; j < len(signalIds); j++ {
			if sigRec.SignalIds[i].String() == ByteSlice2String(signalIds[j]) {
				sigRec.SignalIds = append(sigRec.SignalIds[:i], sigRec.SignalIds[i+1:]...)
				i--
			}
		}
	}

	if len(sigRec.SignalIds) == 0 {
		// a record can not exist without signals to spread it
		updates.AddDeleteQuery(query)
		return updates, uint64(0), nil
	}

	signals, err := (&SignalBucket{r.DB}).GetSignalsByIds(sigRec.SignalIds)
	if err != nil {
		return updates, uint64(0), err
	}
	sigRec.Signals = signals
	sigRec.Sats = SatsSignedFor(sigRec.Signals)

	putQuery, err := r.PutRec(sigRec)
	if err != nil {
		return updates, uint64(0), err
	}
	updates.AddPutQuery(putQuery)
	return updates, sigRec.Sats, nil
}
