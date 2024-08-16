package main

import (
	"cmp"
	"fmt"
	"slices"

	"encoding/json"
	"net/http"

	bolt "go.etcd.io/bbolt"
)

func (ss *SignalServer) getMessageTemplate(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, Signal{}.MessageTemplate(), "{numSats}", "{bitcion address}", "{record id}")
}

func (ss *SignalServer) getPage(w http.ResponseWriter, r *http.Request) {
	recordIds, err := ss.buckets.Rank.GetPageRecordIds([]byte{}, 10)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var records []SerializedRecord
	for i := 0; i < len(recordIds); i++ {
		record, err := ss.buckets.Record.GetRecordById(recordIds[i])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		rec := record.toSerializedRecord()
		rec.ID = record.ID()
		records = append(records, rec)
	}

	json.NewEncoder(w).Encode(records)
}

func (ss *SignalServer) newRecord(w http.ResponseWriter, r *http.Request) {
	var rec Record
	if err := json.NewDecoder(r.Body).Decode(&rec); err != nil {
		ss.errorLog.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	updates, err := ss.IngestRecord(rec)
	if err != nil {
		ss.errorLog.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// debug
	/*
		ss.infoLog.Println("--------------------------")
		ss.infoLog.Println("------------Put-----------")
		for update := range updates.Put {
			ss.infoLog.Printf("\n\n-------\nBucket: %s\n", updates.Put[update].Bucket)
			for kv := range updates.Put[update].KV {
				ss.infoLog.Printf("\nkey: %s\nValue: %s", updates.Put[update].KV[kv].Key, updates.Put[update].KV[kv].Val)
			}
			ss.infoLog.Println("\n\n-------\n")
		}
		ss.infoLog.Println("--------------------------")
		ss.infoLog.Println("----------Delete----------")
		for update := range updates.Delete {
			ss.infoLog.Printf("\n\n-------\nBucket: %s\n", updates.Delete[update].Bucket)
			for kv := range updates.Put[update].KV {
				ss.infoLog.Printf("\nkey: %s\nValue: %s", updates.Delete[update].KV[kv].Key, updates.Put[update].KV[kv].Val)
			}
			ss.infoLog.Printf("\n\n-------\n")
		}
		ss.infoLog.Println("--------------------------")
	*/
	if err := ss.buckets.db.MultiWrite(updates); err != nil {
		ss.errorLog.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(rec.ID())
}

func (ss *SignalServer) checkSizeLimits(r Record) error {
	// checking the actual size of the database will not really
	// be possible because space gets preallocated to be used for
	// new data. If the db is set to double when some threshold
	// is met so new data can populate the new portion of space;
	// the file size will not tell you how much unallocated space
	// is left. Not sure which approach to take yet.
	spaceTaken := uint64(0) // TODO check get size of db
	maxDbSize, _ := Btoi(MaxStorageSize.DefaultBytes())
	spaceLeft := maxDbSize - spaceTaken
	if spaceLeft < r.vBytes() {
		satsPerByte := float64(r.Sats) / float64(r.vBytes())
		lr, err := ss.buckets.Rank.GetLowestRank()
		if err != nil {
			ss.errorLog.Println(err)
			return err
		}
		if satsPerByte < lr {
			return ErrSignalTooWeak
		}
	}
	// check happens in db but why not fail early
	if len(r.Name) > bolt.MaxKeySize {
		return bolt.ErrKeyTooLarge
	}
	if int64(len(r.Value)) > bolt.MaxValueSize {
		// TODO (maybe shouldnt fail here and instead
		// if its a valid signal with enough sats to
		// make the cut then store in fail?
		// NOTE if the todo was done it would probably
		// be best to set the location of the files
		// stored in a user specified location that
		// they can update through the ui
		return bolt.ErrValueTooLarge
	}
	m, _ := Btoi(MaxRecordSize.DefaultBytes())
	if m < r.vBytes() {
		return ErrorRecordTooLarge
	}
	return nil
}

func (ss *SignalServer) dirtyPrepNewRec(r Record) error {
	// might need a re-write at some point. this is a dirty write
	// to avoid checks in the validation code to fail due to not
	// having sufficient information to make a descision; resulting
	// in an error being thrown. maybe pass a temp record around?
	// maybe have a global in mem cache for new records and check
	// there too when validating.
	r.VBytes = r.vBytes()

	// temp remove them when savin
	r.SignalIds = []KV{}

	if len(r.Value) > 0 {
		r.VHash = r.vHash()
	}

	if q, err := ss.buckets.Record.PutRec(r); err != nil {
		ss.errorLog.Println(err)
		return err
	} else {
		if err := ss.buckets.Record.Put(q); err != nil {
			ss.errorLog.Println(err)
			return err
		}
	}

	// generate put query
	putVQuery, putQueryErr := ss.buckets.Value.PutRecV(RecordValue{
		RecID: r.ID(),
		Value: r.Value,
	})
	if putQueryErr != nil {
		if deleteErr := ss.buckets.Record.Delete(ss.buckets.Record.DeleteRec(r.ID())); deleteErr != nil {
			ss.errorLog.Println(deleteErr)
			return deleteErr
		}
		ss.errorLog.Println(putQueryErr)
		return putQueryErr
	}
	if putErr := ss.buckets.Value.Put(putVQuery); putErr != nil {
		if deleteErr := ss.buckets.Record.Delete(ss.buckets.Record.DeleteRec(r.ID())); deleteErr != nil {
			ss.errorLog.Println(deleteErr)
			return deleteErr
		}
		ss.errorLog.Println(putErr)
		return putErr
	}
	return nil
}

func (ss *SignalServer) IngestRecord(r Record) (*DataUpdates, error) {
	if err := ss.checkSizeLimits(r); err != nil {
		return &DataUpdates{}, err
	}

	if len(r.Value) > 0 {
		// to check the record id
		r.VHash = r.vHash()
	}

	isNewRec := false
	internalRecord, recErr := ss.buckets.Record.GetRecordWithSignalsById(r.ID())
	if recErr != nil {
		ss.errorLog.Println(recErr)
		return &DataUpdates{}, recErr
	} else if internalRecord.Sats == 0 {
		if len(r.Value) > 0 {
			r.VBytes = r.vBytes()
			r.VHash = r.vHash()
			r.Sats = r.TotalSats()
		}
		if err := ss.dirtyPrepNewRec(r); err != nil {
			ss.errorLog.Println(err)
			return &DataUpdates{}, err
		}
		isNewRec = true
	} else {
		r.AddSignals(internalRecord.Signals)
		r.Sats = internalRecord.Sats
		r.VHash = internalRecord.VHash
	}

	updates := &DataUpdates{}
	addrPendingSats := map[string]uint64{}
	for i := 0; i < len(r.Signals); i++ {
		// new signal writes updates for all buckets (rank, address,
		// signal and record). This allows for the update to happen
		// in a signgle transaction since all changes with respect
		// to the incoming signals
		r.Signals[i].RecID = r.ID()
		r.Signals[i].VBytes = r.VBytes

		// hash is the id of the signal as a string (.ID() returns a KV)
		signalFound := len(ss.buckets.Signal.GetId(r.Signals[i].ID())) > 0
		sid := r.Signals[i].Hash()
		if !signalFound {
			if _, found := addrPendingSats[sid]; !found {
				addrPendingSats[sid] = 0
			}
			addrPendingSats[sid] += r.Signals[i].Sats

			// should delete the record if issue found
			treeUpdates, updatesErr := ss.NewSignal(r.Signals[i], addrPendingSats[sid])
			if updatesErr != nil && isNewRec {
				if err := ss.buckets.Record.Delete(ss.buckets.Record.DeleteRec(r.ID())); err != nil {
					ss.errorLog.Println(err)
					return &DataUpdates{}, err
				}
				if err := ss.buckets.Value.Delete(ss.buckets.Value.DeleteV(r.ID())); err != nil {
					ss.errorLog.Println(err)
					return &DataUpdates{}, err
				}
				ss.errorLog.Println(updatesErr)
				return &DataUpdates{}, updatesErr
			}
			updates.AppendUpdates(treeUpdates)
		}
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
		ss.errorLog.Println(err)
		return updates, err
	}

	onChain, chainCheckErr := BtcAddressTotal(s.BtcAddress.String())
	if chainCheckErr != nil {
		return updates, chainCheckErr
	}
	// onChain := uint64(1000000)
	if onChain < s.Sats {
		return updates, ErrInsufficientFunds
	}

	// when a new signal comes in it has not been stored yet (assuming it's a valid
	// signal) so checking the address for signed messages will not contain the new
	// signal
	signals, retrieveSignalsErr := ss.buckets.Address.GetSignals(s.BtcAddress)
	if retrieveSignalsErr != nil {
		ss.errorLog.Println(retrieveSignalsErr)
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

	recordUpdater := SignalProcessor{
		updates: &DataUpdates{},

		records:   map[string]*recordTracker{},
		addresses: map[string]*Address{},

		signalsToAdd:    signalsToAdd,
		signalsToDelete: signalsToDelete,

		addSigIds:    make([]KV, len(signalsToAdd)),
		removeSigIds: make([]KV, len(signalsToDelete)),

		buckets: ss.buckets,

		errorLog: ss.errorLog,
		infoLog:  ss.infoLog,
	}

	if err := recordUpdater.AddSignals(signalsToAdd); err != nil {
		ss.errorLog.Println(err)
		return &DataUpdates{}, err
	}
	if err := recordUpdater.DeleteSignals(signalsToDelete); err != nil {
		ss.errorLog.Println(err)
		return &DataUpdates{}, err
	}
	if err := recordUpdater.UpdateAddresses(); err != nil {
		ss.errorLog.Println(err)
		return &DataUpdates{}, err
	}
	if err := recordUpdater.SignalUpdates(); err != nil {
		ss.errorLog.Println(err)
		return &DataUpdates{}, err
	}
	if err := recordUpdater.UpdateRankAndRecord(); err != nil {
		ss.errorLog.Println(err)
		return &DataUpdates{}, err
	}

	return recordUpdater.updates, nil
}
