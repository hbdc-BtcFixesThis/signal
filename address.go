package main

import (
	"encoding/json"
	//"fmt"
)

type Address struct {
	Signals []KV `json:"sig_ids"` // ids of signals (signed messages)
}

func (a *Address) DeleteIds(removeIds []KV) {
	var updatedIds []KV
	var include bool
	for i := 0; i < len(a.Signals); i++ {
		include = true
	out:
		for j := 0; j < len(removeIds); j++ {
			if a.Signals[i].String() == removeIds[j].String() {
				include = false
				break out
			}
		}

		if include {
			updatedIds = append(updatedIds, a.Signals[i])
		}
	}
	a.Signals = updatedIds
}

func (a *Address) PutIds(ids []KV) {
	for i := 0; i < len(ids); i++ {
		include := true
	out:
		for j := 0; j < len(a.Signals); j++ {
			if a.Signals[j].String() == ids[i].String() {
				include = false
				break out
			}
		}

		if include {
			a.Signals = append(a.Signals, ids[i])
		}
	}
}

// one to many where addresses can sign for multiple records which is the signal
// key == bitcoin address
type AddressBucket struct{ *DB }

func (ab *AddressBucket) Name() []byte { return []byte("Address") }

func (ab *AddressBucket) GetAddressB(addr KV) ([]byte, error) {
	query := &Query{
		Bucket:                  ab.Name(),
		KV:                      []Pair{NewPair(addr, nil)},
		CreateBucketIfNotExists: true,
	}
	if err := ab.Get(query); err != nil {
		return []byte{}, err
	}

	return query.KV[0].Val, nil
}

func (ab *AddressBucket) GetAddress(addr KV) (Address, error) {
	var address Address
	addrB, err := ab.GetAddressB(addr)
	if err != nil {
		return Address{Signals: []KV{}}, err
	}

	if len(addrB) == 0 {
		return Address{Signals: []KV{}}, nil
	}

	// the caller cann check the error since something
	// terribly wrong might have happened
	jsonErr := json.Unmarshal(addrB, &address)
	return address, jsonErr
}

func (ab *AddressBucket) UpdateAddrSigs(addrId KV, newIds, removeIds []KV) (*DataUpdates, error) {
	var addr Address
	query := &Query{
		Bucket:                  ab.Name(),
		KV:                      []Pair{NewPair(addrId, nil)},
		CreateBucketIfNotExists: true,
	}

	updates := &DataUpdates{}
	if err := ab.Get(query); err != nil {
		return updates, err
	}
	if err := json.Unmarshal(query.KV[0].Val, &addr); err != nil {
		return updates, err
	}

	addr.DeleteIds(removeIds)
	addr.PutIds(newIds)
	if len(addr.Signals) > 0 {
		q, err := ab.PutAddr(addrId, addr)
		if err != nil {
			return updates, err
		}
		updates.AddPutQuery(q)
	} else {
		updates.AddDeleteQuery(ab.DeleteAddr(addrId))
	}

	return updates, nil
}

func (ab *AddressBucket) PutSignals(s []Signal) (*DataUpdates, error) {
	// assume checks for duplicates and validity has been done
	updates := &DataUpdates{}
	addrs := map[string]Address{}
	for i := 0; i < len(s); i++ {
		addr, found := addrs[s[i].BtcAddress.String()]
		if !found {
			if a, err := ab.GetAddress(s[i].BtcAddress); err != nil {
				return updates, err
			} else {
				addr = a
			}
		}
		addr.PutIds([]KV{s[i].ID()})
	}
	for a := range addrs {
		if query, err := ab.PutAddr(KV(String2ByteSlice(a)), addrs[a]); err != nil {
			// json error
			return &DataUpdates{}, err
		} else {
			updates.AddPutQuery(query)
		}
	}
	return updates, nil
}

func (ab *AddressBucket) PutAddr(addrId KV, address Address) (*Query, error) {
	addrBytes, err := json.Marshal(&address)
	if err != nil {
		return &Query{}, err
	}
	return &Query{
		Bucket:                  ab.Name(),
		KV:                      []Pair{NewPair(addrId, addrBytes)},
		CreateBucketIfNotExists: true,
	}, nil
}

func (ab *AddressBucket) DeleteAddr(addrId KV) *Query {
	return &Query{
		Bucket:                  ab.Name(),
		KV:                      []Pair{NewPair(addrId, nil)},
		CreateBucketIfNotExists: true,
	}
}

func (ab *AddressBucket) GetSignalIdsB(address []byte) KV {
	query := &Query{
		Bucket:                  ab.Name(),
		KV:                      []Pair{NewPair(address, nil)},
		CreateBucketIfNotExists: true,
	}
	ab.Get(query)
	return query.KV[0].Val
}

func (ab *AddressBucket) GetSignals(address []byte) ([]Signal, error) {
	signalBytes := ab.GetSignalIdsB(address)
	if len(signalBytes) == 0 {
		return []Signal{}, nil
	}

	var a Address
	if err := json.Unmarshal(signalBytes, &a); len(a.Signals) == 0 || err != nil {

		return []Signal{}, err
	}

	sb := SignalBucket{ab.DB}
	signals, err := sb.GetSignalsByIds(a.Signals)
	if err != nil {
		return []Signal{}, err
	}

	return signals, nil
}

// func (ab *AddressBucket) SatsSignedFor(address []byte) (uint64, error) {
func SatsSignedFor(signals []Signal) uint64 {
	total := uint64(0)
	for i := 0; i < len(signals); i++ {
		total += signals[i].Sats
	}
	return total
}
