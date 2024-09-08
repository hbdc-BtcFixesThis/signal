package main

import (
	"log"
)

type recordTracker struct {
	record   *Record
	oldTotal uint64
}

type SignalProcessor struct {
	// read only used to generate updates
	updates *DataUpdates

	records   map[string]*recordTracker
	addresses map[string]*Address

	signalsToAdd    []Signal
	signalsToDelete []Signal

	addSigIds    []KV
	removeSigIds []KV

	buckets  *SignalBuckets
	errorLog *log.Logger
	infoLog  *log.Logger
}

func (sp *SignalProcessor) getAddress(signal Signal) error {
	if _, found := sp.addresses[signal.BtcAddress.String()]; !found {
		address, err := sp.buckets.Address.GetAddress(signal.BtcAddress)
		if err != nil {
			sp.errorLog.Println(err)
			return err
		}
		sp.addresses[signal.BtcAddress.String()] = &address
	}
	return nil
}

func (sp *SignalProcessor) getRecord(signal Signal) error {
	if _, found := sp.records[signal.RecID.String()]; !found {
		r, err := sp.buckets.Record.GetRecordWithSignalsById(signal.RecID)
		if err != nil {
			sp.errorLog.Println(err)
			return err
		}
		sp.records[signal.RecID.String()] = &recordTracker{record: &r, oldTotal: r.TotalSats()}
	}
	return nil
}

func (sp *SignalProcessor) getAddrAndRecFromSigal(signal Signal) error {
	if err := sp.getAddress(signal); err != nil {
		sp.errorLog.Println(err)
		return err
	}
	if err := sp.getRecord(signal); err != nil {
		sp.errorLog.Println(err)
		return err
	}
	return nil
}

func (sp *SignalProcessor) AddSignals(signals []Signal) error {
	for i := 0; i < len(signals); i++ {
		sp.addSigIds[i] = signals[i].ID()
		if err := sp.getAddrAndRecFromSigal(signals[i]); err != nil {
			sp.errorLog.Println(err)
			return err
		}
		// recomputes record sats total in the record (can now be used to update rank)
		sp.records[signals[i].RecID.String()].record.AddSignal(signals[i])
	}
	return nil
}

func (sp *SignalProcessor) DeleteSignals(signals []Signal) error {
	for i := 0; i < len(signals); i++ {
		sp.removeSigIds[i] = signals[i].ID()
		if err := sp.getAddrAndRecFromSigal(signals[i]); err != nil {
			sp.errorLog.Println(err)
			return err
		}
		// recomputes record sats total (can now be used to update rank)
		sp.records[signals[i].RecID.String()].record.RemoveSignal(sp.removeSigIds[i])
	}
	return nil
}

func (sp *SignalProcessor) UpdateAddresses() error {
	for addrId := range sp.addresses {
		sp.addresses[addrId].DeleteIds(sp.removeSigIds)
		sp.addresses[addrId].PutIds(sp.addSigIds)
		if len(sp.addresses[addrId].Signals) > 0 {
			putQuery, err := sp.buckets.Address.PutAddr(KV(addrId), *sp.addresses[addrId])
			if err != nil {
				sp.errorLog.Println(err)
				return err
			}
			sp.updates.AddPutQuery(putQuery)
		} else {
			sp.updates.AddDeleteQuery(sp.buckets.Address.DeleteAddr(KV(addrId)))
		}
	}
	return nil
}

func (sp *SignalProcessor) SignalUpdates() error {
	if len(sp.removeSigIds) > 0 {
		sp.updates.AddDeleteQuery(sp.buckets.Signal.DeleteIds(sp.removeSigIds))
	}
	signalUpdates, sigUpdatesErr := sp.buckets.Signal.PutSignals(sp.signalsToAdd)
	if sigUpdatesErr != nil {
		sp.errorLog.Println(sigUpdatesErr)
		return sigUpdatesErr
	}
	sp.updates.AddPutQuery(signalUpdates)
	return nil
}

func (sp *SignalProcessor) UpdateRankAndRecord() error {
	// update records and ranks
	for recId := range sp.records {
		record := sp.records[recId].record
		if record.TotalSats() > 0 {
			putQuery, putRecErr := sp.buckets.Record.PutRec(*record)
			if putRecErr != nil {
				sp.errorLog.Println(putRecErr)
				return putRecErr
			}
			sp.updates.AddPutQuery(putQuery)
		} else {
			sp.updates.AddDeleteQuery(sp.buckets.Record.DeleteRec(record.ID()))
			// share the same record  id
			sp.updates.AddDeleteQuery(sp.buckets.Value.DeleteV(record.ID()))
		}

		oldRank := record.RankForSatCountB(sp.records[recId].oldTotal)
		if record.TotalSats() > uint64(0) {
			rankUpdates, reRankErr := sp.buckets.Rank.ReRankRec(
				// (old rank, new rank, record id)
				oldRank, record.RankB(), record.ID(),
			)
			if reRankErr != nil {
				sp.errorLog.Println(reRankErr)
				return reRankErr
			}
			sp.updates.AppendUpdates(rankUpdates)
		} else {
			updates, err := sp.buckets.Rank.deleteRecFromRank(oldRank, record.ID())
			if err != nil {
				sp.errorLog.Println(err)
				return err
			}
			sp.updates.AppendUpdates(updates)
		}
	}
	return nil
}
