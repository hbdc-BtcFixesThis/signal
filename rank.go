package main

import (
	"encoding/json"
)

// this is an index to look records up by name
// since only the records name plus value are required to be unique
type NameBucket struct{ *DB }

type Rank struct {
	Records []KV `json:"record_ids"`
}

const RankBucketName = "Rank" // key = r.TotalSats();  val = [list of records]

// key = sats/byte; val = list of rec hashes
type RankBucket struct{ *DB }

func (r *RankBucket) Name() []byte { return []byte(RankBucketName) }

func (r *RankBucket) GetLowestRank() (float64, error)  { return r.getLast(true) }
func (r *RankBucket) GetHighestRank() (float64, error) { return r.getLast(false) }
func (r *RankBucket) getLast(ascending bool) (float64, error) {
	pq := &PageQuery{
		Ascending: ascending,
		Query: &Query{
			Bucket:                  r.Name(),
			KV:                      make([]Pair, 1),
			CreateBucketIfNotExists: true,
		},
	}

	r.GetPage(pq)
	if len(pq.KV[0].Val) == 0 {
		return 0.0, nil
	}

	return F64fb(pq.KV[0].Key)
}

func (r *RankBucket) GetId(id []byte) ([]byte, error) {
	q := &Query{
		Bucket:                  r.Name(),
		KV:                      []Pair{NewPair(id, nil)},
		CreateBucketIfNotExists: true,
	}
	err := r.Get(q)
	return q.KV[0].Val, err
}

func (r *RankBucket) PutId(id, val []byte) *Query {
	return &Query{
		Bucket:                  r.Name(),
		KV:                      []Pair{NewPair(id, val)},
		CreateBucketIfNotExists: true,
	}
}

func (r *RankBucket) PutSignalRec(sr Record) (*Query, error) {
	key := F64tb(float64(sr.Sats) / float64(sr.VBytes()))
	rankB, getErr := r.GetId(key)
	if getErr != nil {
		return &Query{}, getErr
	}

	var rank Rank
	if err := json.Unmarshal(rankB, &rank); err != nil {
		return &Query{}, err
	}
	for i := 0; i < len(rank.Records); i++ {
		if rank.Records[i].String() == sr.Hash() {
			// TODO make sure this case is handled when doing batch updates
			// record exists in rank; nothing to do
			return &Query{}, nil
		}
	}
	rank.Records = append(rank.Records, KV(sr.ID()))

	b, err := json.Marshal(rank)
	if err != nil {
		return &Query{}, err
	}
	return r.PutId(key, b), nil
}

func (r *RankBucket) deleteRecFromRank(currentRankB, recId []byte) (*DataUpdates, error) {
	updates := &DataUpdates{}
	qCurrentRank := &Query{
		Bucket:                  r.Name(),
		KV:                      []Pair{NewPair(currentRankB, nil)},
		CreateBucketIfNotExists: true,
	}
	if err := r.Get(qCurrentRank); err != nil {
		return updates, err
	}
	if len(qCurrentRank.KV[0].Val) == 0 {
		// for new records do nothing
		return updates, nil
	}

	var currentRank Rank
	if err := json.Unmarshal(qCurrentRank.KV[0].Val, &currentRank); err != nil {
		return updates, err
	}
	found := false
	for i := 0; i < len(currentRank.Records); i++ {
		if currentRank.Records[i].String() == ByteSlice2String(recId) {
			found = true
			currentRank.Records = append(currentRank.Records[:i], currentRank.Records[i+1:]...)
			i--
		}
	}
	if len(currentRank.Records) == 0 {
		updates.AddDeleteQuery(qCurrentRank)
		return updates, nil
	}
	if found {
		if b, err := json.Marshal(currentRank); err != nil {
			return updates, err
		} else {
			qCurrentRank.KV[0].Val = b
			updates.AddPutQuery(qCurrentRank)
			return updates, nil
		}
	}
	return updates, nil
}

func (r *RankBucket) ReRankRec(currentRankB []byte, newRankB []byte, recId KV) (*DataUpdates, error) {
	updates, err := r.deleteRecFromRank(currentRankB, recId)
	if err != nil {
		return updates, err
	}
	qNewRank := &Query{
		Bucket:                  r.Name(),
		KV:                      []Pair{NewPair(newRankB, nil)},
		CreateBucketIfNotExists: true,
	}
	var newRank Rank
	if err := r.Get(qNewRank); err != nil {
		return &DataUpdates{}, err
	}
	if len(qNewRank.KV[0].Val) > 0 {
		if err := json.Unmarshal(qNewRank.KV[0].Val, newRank); err != nil {
			return &DataUpdates{}, err
		}
		for i := 0; i < len(newRank.Records); i++ {
			if newRank.Records[i].String() == ByteSlice2String(recId) {
				return updates, nil // rec id already in specified rank
			}
		}
	}
	newRank.Records = append(newRank.Records, recId)
	b, err := json.Marshal(newRank)
	if err != nil {
		return &DataUpdates{}, err
	}
	qNewRank.KV[0].Val = b
	updates.AddPutQuery(qNewRank)
	return updates, nil
}

func (r *RankBucket) GetPageRecordIds(last []byte, size int) ([][]byte, error) {
	pq := &PageQuery{
		Query: &Query{
			Bucket:                  r.Name(),
			KV:                      make([]Pair, 1),
			CreateBucketIfNotExists: true,
		},
	}

	results := [][]byte{}
	var rank Rank
out:
	for len(results) < size {
		if len(last) > 0 {
			// might be a problem
			pq.StartFrom = KV(last)
		}
		r.GetPage(pq)
		if len(pq.KV[0].Val) == 0 {
			break
		}

		err := json.Unmarshal(pq.KV[0].Val, &rank)
		if err != nil {
			return results, err
		}
		newRecFound := true
		for i := 0; i < len(rank.Records); i++ {
			// NOTE
			// including this condition means there will
			// need to be a way to determine the last
			// record the user lands on. Multiple records
			// can be in a rank and if the size of the
			// results is full while the rank has records
			// left in it then as of right now they will
			// be missed for the next page load if this
			// condition is re-introduced
			//
			// if len(results) == size {
			//	break out
			// }

			for j := 0; j < len(results); j++ {
				if ByteSlice2String(results[j]) == ByteSlice2String(rank.Records[i]) {
					newRecFound = false
				}
			}
			if newRecFound {
				results = append(results, rank.Records[i])
			} else {
				// if there are less records then the size
				// specified and this condition is not put
				// in then it will continue to add the same
				// record till the results reach the size
				// specified. If there are more records then
				// the size specified but
				// # records % size != 0
				// and the last page is being queried then
				// it will go back to the beginning which
				// is probably fine for now
				break out
			}
		}
		last = pq.Query.KV[0].Key
	}

	return results, nil
}
