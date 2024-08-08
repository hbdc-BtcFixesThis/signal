package main

import (
	"cmp"
	"slices"

	"encoding/json"
	"net/http"
)

func (ss *SignalServer) newRecord(w http.ResponseWriter, r *http.Request) {
	var rec Record
	if err := json.NewDecoder(r.Body).Decode(&rec); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	updates, err := ss.IngestRecord(rec)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := ss.buckets.db.MultiWrite(updates); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(rec.ID())
}

func (ss *SignalServer) getPage(w http.ResponseWriter, r *http.Request) {
	recordIds, err := ss.buckets.Rank.GetPageRecordIds([]byte{}, 10)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var records []Record
	for i := 0; i < len(recordIds); i++ {
		record, err := ss.buckets.Record.GetRecordById(recordIds[i])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		records = append(records, record)
	}
	json.NewEncoder(w).Encode(records)
}

func (ss *SignalServer) IngestRecord(r Record) (*DataUpdates, error) {
	// checking the actual size of the database will not really
	// be possible because space gets preallocated to be used for
	// new data. If the db is set to double when some threshold
	// is met so new data can populate the new portion of space;
	// the file size will not tell you how much unallocated space
	// is left. Not sure which approach to take yet.

	spaceTaken := uint64(0)                                       // TODO check get size of db
	spaceLeft := Btoi(MaxStorageSize.DefaultBytes()) - spaceTaken // Btoi Byte to uint64
	if spaceLeft < r.VBytes() {
		satsPerByte := float64(r.Sats) / float64(r.VBytes())
		if satsPerByte < ss.buckets.Rank.GetLowestRank() {
			return &DataUpdates{}, ErrSignalTooWeak
		}
	}

	if Btoi(MaxRecordSize.DefaultBytes()) < r.VBytes() {
		return &DataUpdates{}, ErrorRecordTooLarge
	}
	isNewRec := false
	if v, err := ss.buckets.Record.GetId(r.ID()); err != nil {
		return &DataUpdates{}, err
	} else if len(v) == 0 {
		// temp remove them when savin
		signalIds := r.SignalIds
		r.SignalIds = []KV{}

		// might need a re-write at some point. this is a dirty write
		// to avoid checks in the validation code to fail due to not
		// having sufficient information to make a descision; resulting
		// in an error being thrown. maybe pass a temp record around?
		// maybe have a global in mem cache for new records and check
		// there too when validating.
		q, err := ss.buckets.Record.PutRec(r)
		if err != nil {
			return &DataUpdates{}, err
		}
		if err := ss.buckets.Record.Put(q); err != nil {
			return &DataUpdates{}, err
		}
		r.SignalIds = signalIds
		isNewRec = true
	}

	updates := &DataUpdates{}
	addrPendingSats := map[string]uint64{}
	for i := 0; i < len(r.Signals); i++ {
		// new signal writes updates for all buckets (rank, address,
		// signal and record). This allows for the update to happen
		// in a signgle transaction since all changes with respect
		// to the incoming signals
		r.Signals[i].RecID = r.ID()

		// hash is the id of the signal as a string (.ID() returns a KV)
		sid := r.Signals[i].Hash()
		if _, found := addrPendingSats[sid]; !found {
			addrPendingSats[sid] = r.Signals[i].Sats
		} else {
			addrPendingSats[sid] += r.Signals[i].Sats
		}

		// should delete the record if issue found
		treeUpdates, updatesErr := ss.NewSignal(r.Signals[i], addrPendingSats[sid])
		if updatesErr != nil && isNewRec {
			deleteErr := ss.buckets.Record.Delete(ss.buckets.Record.DeleteRec(r.ID()))
			if deleteErr != nil {
				return &DataUpdates{}, deleteErr
			}
			return &DataUpdates{}, updatesErr
		}
		updates.AppendUpdates(treeUpdates)
	}
	return updates, nil
}

func (ss *SignalServer) NewSignal(s Signal, pendingSats uint64) (*DataUpdates, error) {
	// NOTE
	// pending sats is a hack. When a record comes in it may have multiple singals
	// that have not yet been written. This is because updates to the state of the
	// signal db are done in a transaction where queries are batched to avoid dirty
	// writes that may leave unwanted records behind; or worse, invalid states.

	updates := &DataUpdates{}
	if len(ss.buckets.Signal.GetId(s.ID())) > 0 {
		// nothing to do here
		return updates, nil
	}

	// TODO configurable? no intervention outside of > 0? any downsides to making
	// this configurable?
	if s.Sats < 1 {
		return updates, ErrorNeedMoreSats
	}
	if err := s.CheckSignature(); err != nil {
		return updates, err
	}

	onChain, chainCheckErr := BtcAddressTotal(s.BtcAddress.String())
	if chainCheckErr != nil {
		return updates, chainCheckErr
	}
	if onChain < s.Sats {
		return updates, ErrInsufficientFunds
	}

	// when a new signal comes in it has not been stored yet (assuming it's a valid
	// signal) so checking the address for signed messages will not contain the new
	// signal
	signals, retrieveSignalsErr := ss.buckets.Address.GetSignals(s.BtcAddress)
	if retrieveSignalsErr != nil {
		return updates, retrieveSignalsErr
	}
	for i := 0; i < len(signals); i++ {
		if signals[i].Hash() == s.Hash() {
			// since the id is a hash of the signature, the signal already exists
			return updates, nil
		}
	}

	satsLeft := uint64(onChain) - SatsSignedFor(signals)
	if (satsLeft - s.Sats - pendingSats) < s.Sats {
		// pass in onchain total to avoid reaching out for it again
		// TODO maybe cache these and update whenever blockheight
		// changes to avoid unnecessary questions
		return ss.insufficientFundsSignalReorg(s, signals, onChain)
	}

	return ss.updateRecordSignals([]Signal{s}, []Signal{})
}

func (ss *SignalServer) insufficientFundsSignalReorg(s Signal, addrSignals []Signal, onchainAddrTotal uint64) (*DataUpdates, error) {
	// steps:
	// 1) check that the incoming signal record has only one signal
	// 2) find and sort all signals for address by sats per byte
	// 3) Check to see if reorg is applicable by comparing the new
	//    incoming sigal in sr to the first signal in the ordered signals
	// 4) Go through the other signals and add up the sats and bytes
	//    to check if the new signal should replace one or more signals;
	//    determined by sats/byte

	newSignalSatsPerByte := float64(s.Sats) / float64(s.VBytes)
	slices.SortFunc(addrSignals, func(a, b Signal) int {
		return cmp.Compare(
			float64(a.Sats)/float64(a.VBytes),
			float64(b.Sats)/float64(b.VBytes),
		)
	})

	if newSignalSatsPerByte < float64(addrSignals[0].Sats)/float64(addrSignals[0].VBytes) {
		return &DataUpdates{}, ErrInsufficientFundsNoReorg
	}
	satsCount := uint64(0)
	vByteCount := uint64(0)
	replaceUpTo := 0
	for i := 1; i < len(addrSignals); i++ {
		satsCount += addrSignals[i].Sats
		vByteCount += addrSignals[i].VBytes
		if newSignalSatsPerByte > float64(satsCount)/float64(vByteCount) {
			if onchainAddrTotal-satsCount > s.Sats {
				replaceUpTo = i
				break
			}
		} else {
			// if the sats per byte for the new record is smaller
			// going forward will only increase gap; since the
			// loop is over an ordered slice
			break
		}
	}
	if replaceUpTo == 0 {
		return &DataUpdates{}, ErrSignalTooWeak
	}

	u, err := ss.updateRecordSignals([]Signal{s}, addrSignals[0:replaceUpTo])
	return u, err
}

func (ss *SignalServer) updateRecordSignals(signalsToAdd []Signal, signalsToDelete []Signal) (*DataUpdates, error) {
	// its long but this is done here to be able to keep track of the totals
	// without leaving data in a dirty state. This allows for updates to the
	// rankings and record before save. That way if we fail in the middle no
	// writes have yet been done and there is no need for a rollback. All
	// writes will be saved for the last step to allow for the updates to
	// happen in a transaction which helps ensure integrity (ie sats per
	// signal found in a record will add up to whats in the record, signals
	// are valid and acounted for in the rec ids, sigals in the address bucket
	// are up to date with changes and the record has been re-ranked based
	// on sat/vbyte.
	type recordTracker struct {
		record   *Record
		oldTotal uint64
	}

	updates := &DataUpdates{}

	records := map[string]recordTracker{}
	addresses := map[string]*Address{}

	addSigIds := make([]KV, len(signalsToAdd))
	removeSigIds := make([]KV, len(signalsToDelete))

	updateStructs := func(signal Signal) error {
		// for address updates
		if _, found := addresses[signal.BtcAddress.String()]; !found {
			address, err := ss.buckets.Address.GetAddress(signal.BtcAddress)
			if err != nil {
				return err
			}
			addresses[signal.BtcAddress.String()] = &address
		}

		// for rank updates
		if _, found := records[signal.RecID.String()]; !found {
			r, err := ss.buckets.Record.GetRecordWithSignalsById(signal.RecID)
			if err != nil {
				return err
			}
			records[signal.RecID.String()] = recordTracker{record: &r, oldTotal: r.TotalSats()}
		}
		return nil

	}
	for i := 0; i < len(signalsToAdd); i++ {
		addSigIds[i] = signalsToAdd[i].ID()
		if err := updateStructs(signalsToAdd[i]); err != nil {
			return &DataUpdates{}, err
		}
		// recomputes record sats total in the record (can now be used to update rank)
		records[signalsToAdd[i].RecID.String()].record.AddSignal(signalsToAdd[i])
		addresses[signalsToAdd[i].BtcAddress.String()].PutIds([]KV{addSigIds[i]})
	}

	for i := 0; i < len(signalsToDelete); i++ {
		removeSigIds[i] = signalsToDelete[i].ID()
		if err := updateStructs(signalsToDelete[i]); err != nil {
			return &DataUpdates{}, err
		}
		// recomputes record sats total (can now be used to update rank)
		records[signalsToDelete[i].RecID.String()].record.RemoveSignal(removeSigIds[i])
		addresses[signalsToDelete[i].BtcAddress.String()].DeleteIds([]KV{removeSigIds[i]})
	}

	// addresses updates
	for addrId := range addresses {
		if len(addresses[addrId].Signals) > 0 {
			putQuery, err := ss.buckets.Address.PutAddr(KV(addrId), *addresses[addrId])
			if err != nil {
				return &DataUpdates{}, err
			}
			updates.AddPutQuery(putQuery)
		} else {
			updates.AddDeleteQuery(ss.buckets.Address.DeleteAddr(KV(addrId)))
		}
	}

	// signal updates
	if len(removeSigIds) > 0 {
		updates.AddDeleteQuery(ss.buckets.Signal.DeleteIds(removeSigIds))
	}
	signalUpdates, sigUpdatesErr := ss.buckets.Signal.PutSignals(signalsToAdd)
	if sigUpdatesErr != nil {
		return &DataUpdates{}, sigUpdatesErr
	}
	updates.AddPutQuery(signalUpdates)

	// update records and ranks
	for recId := range records {
		record := records[recId].record
		if record.TotalSats() > 0 {
			putQuery, putRecErr := ss.buckets.Record.PutRec(*record)
			if putRecErr != nil {
				return &DataUpdates{}, putRecErr
			}
			updates.AddPutQuery(putQuery)
		} else {
			updates.AddDeleteQuery(ss.buckets.Record.DeleteRec(record.ID()))
		}

		rankUpdates, reRankErr := ss.buckets.Rank.ReRankRec(
			// (old rank, new rank, record id)
			record.RankForSatCountB(records[recId].oldTotal), record.RankB(), record.ID(),
		)
		if reRankErr != nil {
			return &DataUpdates{}, reRankErr
		}
		updates.AppendUpdates(rankUpdates)
	}
	return updates, nil
}
