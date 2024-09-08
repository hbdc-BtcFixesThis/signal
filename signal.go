package main

import (
	"fmt"

	"encoding/json"

	verifier "github.com/bitonicnl/verify-signed-message/pkg"
	"github.com/btcsuite/btcd/chaincfg"
)

const SignalBucketName = "Signals"

// signal sent by bitcoiner who signed a message to spread a record
type Signal struct {
	Sats       uint64 `json:"sats"`
	BtcAddress KV     `json:"btc_address"`
	Signature  KV     `json:"signature"`
	VBytes     uint64 `json:"vbytes"`
	RecID      KV     `json:"rid"`
	// NOTE VBytes is added to be able to compare signals in order to
	// make descisions about state changes without incuring another
	// trip to the db to lookup the record (should this be optimized
	// for storage instead?)

	// ::TODO:: answer question
	// should signature and sats be slices?
	// why the question? Let's say a record is signed for an address.
	// Later, the owner of the address to add to the signal without
	// but wants to do so without moving funds. To accomplish this a
	// signal (signed message) from  the same address would be broadcasted.
	// If this is done twice for the same record then storing all the
	// signals for the same entry would avoid duplicating the address.
	// Hash would need to hash all the signatures to be used as an ID.
	// The size of a btc address is 25 bytes. Let's say this grows and
	// somehow gets 10,000,000 address that have signed messages. And
	// let's say every address signed two signals per record then the
	// duplicate btc address would jump from 1/4 to 1/2 a gb. So there's
	// time before this needs to be address and future migrations should
	// be able to transform the data if needed. The other option here is
	// to only allow one signal per record/address pair. That would mean
	// the only metric that can rerliably be used to determinie the
	// outcome would be the strongest signal (sats/byte). That is due
	// to the fact that the cost to produce a message is 0 (different
	// from the cost to acquire the funds needed to make a valid signal).
	// If a date is required any date can be enterd and signed.
}

func (s *Signal) Hash() string { return SHA256(s.Signature) }

func (s *Signal) ID() KV { return String2ByteSlice(s.Hash()) }

func (s Signal) MessageTemplate() string {
	return `This is not a bitcoin transaction!

For as long as the there are at least

SATS: %v

unspent in Bitcoin

ADDRESS: %s

may they be used to spread

RECORD ID: %s


Peace and love freaks`
}

func (s *Signal) MessageToSign() string {
	return fmt.Sprintf(s.MessageTemplate(), s.Sats, s.BtcAddress, s.RecID)
}

func (s *Signal) CheckSignature() error {
	valid, err := verifier.VerifyWithChain(verifier.SignedMessage{
		Address:   s.BtcAddress.String(),
		Message:   s.MessageToSign(),
		Signature: s.Signature.String(),
	}, &chaincfg.MainNetParams) // &chaincfg.TestNet3Params)
	if err != nil {
		return err
	}
	if !valid {
		return ErrorInvalidSignature
	}
	return nil
}

// hash signature to save some space
// key = SHA256(Signature); val = serialized signal
type SignalBucket struct{ *DB }

func (sb *SignalBucket) Name() []byte { return []byte(SignalBucketName) }

func (sb *SignalBucket) GetId(id KV) KV {
	query := &Query{
		Bucket:                  sb.Name(),
		KV:                      []Pair{NewPair(id, nil)},
		CreateBucketIfNotExists: true,
	}
	sb.Get(query)
	return query.KV[0].Val
}

func (sb *SignalBucket) PutSignals(signals []Signal) (*Query, error) {
	// this only accounts for the signal bucket
	// rank, and record buckets also need to update state
	// but this happes elsewhere. Just a warning
	pairs := make([]Pair, len(signals))
	for i := 0; i < len(signals); i++ {
		b, err := json.Marshal(signals[i])
		if err != nil {
			return &Query{}, err
		}
		pairs[i] = NewPair(signals[i].ID(), b)
	}
	return &Query{
		Bucket:                  sb.Name(),
		KV:                      pairs,
		CreateBucketIfNotExists: true,
	}, nil
}

func (sb *SignalBucket) DeleteIds(ids []KV) *Query {
	pairs := make([]Pair, len(ids))
	for i := 0; i < len(ids); i++ {
		pairs[i] = NewPair(ids[i], nil)
	}
	return &Query{
		Bucket:                  sb.Name(),
		KV:                      pairs,
		CreateBucketIfNotExists: true,
	}
}

func (sb *SignalBucket) GetIds(ids []KV) (*Query, error) {
	kv := make([]Pair, len(ids))
	for i := 0; i < len(ids); i++ {
		kv[i] = NewPair(ids[i], nil)
	}

	query := &Query{
		Bucket:                  sb.Name(),
		KV:                      kv,
		CreateBucketIfNotExists: true,
	}
	err := sb.Get(query)
	return query, err
}

func (sb *SignalBucket) GetSignalsByIds(ids []KV) ([]Signal, error) {
	signals := make([]Signal, len(ids))
	results, err := sb.GetIds(ids)
	if err != nil {
		return []Signal{}, err
	}

	for i := 0; i < len(results.KV); i++ {
		// this should not happen
		if uErr := json.Unmarshal(results.KV[i].Val, &signals[i]); uErr != nil {
			return []Signal{}, uErr
		}
	}
	return signals, nil
}
