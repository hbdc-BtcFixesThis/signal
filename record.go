package main

import (
	"encoding/binary"
	"encoding/json"
)

// content of record
// the result of RecordHash below is used as the id for these records
type Record struct {
	Sats      uint64   `json:"sats"`
	Name      KV       `json:"key"`
	Value     KV       `json:"value"`
	Signals   []Signal `json:"signals"`
	SignalIds []KV     `json:"signal_ids"` // ids for signals in SignalBucket
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

func (r *Record) VBytes() uint64 { return uint64(binary.Size(r.Name) + binary.Size(r.Value)) }

func (r *Record) Rank() float64 {
	return float64(float64(r.TotalSats()) / float64(r.VBytes()))
}

func (r *Record) RankB() []byte {
	return F64tb(r.Rank())
}

func (r *Record) RankForSatCount(sats uint64) float64 {
	return float64(float64(sats) / float64(r.VBytes()))
}
func (r *Record) RankForSatCountB(sats uint64) []byte {
	return F64tb(r.RankForSatCount(sats))
}

func (r *Record) Hash() string {
	return SHA256(String2ByteSlice(SHA256(r.Name) + "::" + SHA256(r.Value)))
}

func (r *Record) ID() []byte { return String2ByteSlice(r.Hash()) }

// key = r.RecordHash(); val = serialize record
type RecordBucket struct{ *DB }

func (r *RecordBucket) Name() []byte { return []byte("Record") }

func (r *RecordBucket) GetId(id []byte) ([]byte, error) {
	query := &Query{
		Bucket:                  r.Name(),
		KV:                      []Pair{NewPair(id, nil)},
		CreateBucketIfNotExists: true,
	}
	err := r.Get(query)
	return query.KV[0].Val, err
}

func (r *RecordBucket) PutRec(rec Record) (*Query, error) {
	// hacky. should omit signals if they exist in the record
	rec.Signals = []Signal{}

	b, err := json.Marshal(rec)
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

func (r *RecordBucket) GetRecordById(id []byte) (Record, error) {
	var signalRec Record

	rBytes, err := r.GetId(id)
	if err != nil {
		return signalRec, err
	}
	err = json.Unmarshal(rBytes, &signalRec)

	return signalRec, err
}

func (r *RecordBucket) GetRecordWithSignalsById(id []byte) (Record, error) {
	record, err := r.GetRecordById(id)
	if err != nil {
		return record, err
	}
	signals, err := (&SignalBucket{r.DB}).GetSignalsByIds(record.SignalIds)
	if err != nil {
		return record, err
	}
	record.Signals = signals
	return record, nil
}

func (r *Record) MessageToBeSigned() string {
	return `This is not a bitcoin transaction!

For as long as the there are funds left
unspent in bitcoin wallet address

{bitcion address}

{amount} sats of the balance will be
used to spread the following record

{name label}:
{name}
{content label}:
{content}

Peace and love freaks`
}
