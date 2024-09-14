package main

import (
	"bytes"
	"cmp"
	"context"
	"encoding/json"
	"slices"
	"time"
)

func (ss *SignalServer) runAddressMonitor(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			ss.errorLog.Println("Stopping runAddressMonitor!")
			return
		case <-time.After(30 * time.Second):
			reCheckAddresses := false
			height, err := LatestBtcBlockHeight()
			if err == nil {
				ss.Lock()
				if ss.lastBlock != height {
					ss.infoLog.Printf(
						"[Block height changed], ss.lastBlock: %v, height: %v",
						ss.lastBlock, height,
					)
					ss.lastBlock = height
					reCheckAddresses = true
				}
				ss.Unlock()
			}
			if reCheckAddresses {
				ss.checkForAddressChanges()
			}
		}
	}
}

func (ss *SignalServer) checkForAddressChanges() {
	pq := &PageQuery{
		Query: &Query{
			Bucket:                  ss.buckets.Address.Name(),
			KV:                      make([]Pair, 1),
			CreateBucketIfNotExists: true,
		},
	}

	var addr Address
	var last []byte
	for {
		if len(last) > 0 {
			pq.StartFrom = KV(last)
		}
		ss.buckets.Address.GetPage(pq)
		if bytes.Equal(pq.KV[0].Key, last) {
			ss.infoLog.Printf("pq.KV: %+v", pq.KV)
			ss.infoLog.Println("Done checking all addresses!")
			break
		}
		last = pq.Query.KV[0].Key

		if len(pq.KV[0].Val) == 0 {
			ss.errorLog.Printf(
				"Found address with nothing??? key: %v, val: %v",
				pq.KV[0].Key, pq.KV[0].Val,
			)
			continue
		}

		if err := json.Unmarshal(pq.KV[0].Val, &addr); err != nil {
			ss.errorLog.Println(err)
			continue
		}

		onchainTotal, onchainErr := BtcAddressTotal(pq.KV[0].Key.String())
		if onchainErr != nil {
			// signals, sigRetrieveErr := ss.buckets.Signal.GetSignalsByIds(addr.Signals)
			// if sigRetrieveErr {
			//		ss.errorLog.Println(err)
			//	}
			// ss.updateRecordSignals([]Signal{}, signals)
			ss.errorLog.Printf("Error getting onchain total for addr %s", pq.KV[0].Key)
			ss.errorLog.Println(onchainErr)
			continue
		}
		if onchainTotal != addr.SatsSinceLastChecked {
			signals, sigRetrieveErr := ss.buckets.Signal.GetSignalsByIds(addr.Signals)
			if sigRetrieveErr != nil {
				ss.errorLog.Printf("Error retrieving signals %+v", addr)
				ss.errorLog.Println(sigRetrieveErr)
				continue
			}

			var updateErr error
			updates := &DataUpdates{}

			totalSignedFor := SatsSignedFor(signals)
			if totalSignedFor > onchainTotal {
				slices.SortFunc(signals, func(a, b Signal) int {
					return cmp.Compare(
						float64(a.Sats)/float64(a.VBytes),
						float64(b.Sats)/float64(b.VBytes),
					)
				})
				replaceUpTo := 1
				newTotal := totalSignedFor
				for i := 0; i < len(signals); i++ {
					replaceUpTo = i + 1
					newTotal -= signals[i].Sats
					if newTotal >= onchainTotal {
						break
					}
				}
				updates, updateErr = ss.updateRecordSignals([]Signal{}, signals[0:replaceUpTo])
				if updateErr != nil {
					ss.errorLog.Println(updateErr)
				}
				addr.Signals = addr.Signals[replaceUpTo:len(addr.Signals)]
			} else {
				addr.SatsSinceLastChecked = onchainTotal
				q, err := ss.buckets.Address.PutAddr(pq.KV[0].Key, addr)
				if err != nil {
					ss.errorLog.Println(err)
				}
				updates.AddPutQuery(q)
			}
			// ss.logNewRecordAndOrSignalUpdates(updates)
			if err := ss.buckets.db.MultiWrite(updates); err != nil {
				ss.errorLog.Println(err)
			}
		}
	}
}
