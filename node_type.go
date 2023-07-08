package main

import (
	"fmt"

	"encoding/json"
)

type NodeType uint8

const (
	// To aquire a slot here anyone can compete.
	// This acts as more of a free market for
	// storage on a redundant db replicated on
	// all the machines runnning this code. At
	// its core the record is really split into
	// two parts; key and value. The record can
	// be any size and neither part on its own
	// required to be unique. As long as there
	// is at least one sat in a signal the
	// record is valid. If a node runs out of
	// space only the srongest signals survive.
	// A signal is like a vote. To broacast one
	// you sign a message with the keys to some
	// bitcoin.
	SIGNAL NodeType = iota

	// Admin/s grant read/write perms.
	OTHER

	// 1k sats per byte required for an entry
	// in a domain dataset. No one but the
	// owner of the record signs a message with
	// respect to the record. The record has a
	// user deifined unique id as the key and
	// a value associated with that key. If
	// records colide the one with the largest
	// bitcoin/byte wins and retains it's spot
	// in the db. If I rely on the date of the
	// first utxo in an address someone may use
	// that to acquiure a domain they have no
	// valid claim on. If I rely on a signature,
	// the claims on that signature are easily
	// faked. The only reliable game I can think
	// up is to rank everything by sats per byte
	// and may the best strongest signal win. A
	// nice byproduct of that is the you are able
	// to change your record as long as the update
	// takes up less space. You can always move
	// your funds to another address (or even just
	// enough to go below the required amount) and
	// that record looses its place and the domain
	// is once again unoccupied. A maliciouos actor
	// may spam the system, however, the cost to
	// do so comes at the expense of locked up funds.
	// DOMAIN
)

// MarshalText implements the encoding.TextMarshaler interface.
func (nt NodeType) MarshalText() ([]byte, error) {
	return []byte(nt.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (nt *NodeType) UnmarshalText(b []byte) error {
	aux := string(b)
	for i, name := range nt.PossibleNodeTypes() {
		if name == aux {
			*nt = NodeType(i)
			return nil
		}
	}
	return fmt.Errorf("invalid node type %q", aux)
}

func (nt NodeType) PossibleNodeTypes() []string {
	return []string{"signal", "other"}
}

func (nt NodeType) PossibleNodeTypesString() string {
	ret, _ := json.Marshal(SIGNAL.PossibleNodeTypes())
	return string(ret)
}

func (nt NodeType) String() string { return nt.PossibleNodeTypes()[nt] }

func (nt NodeType) IsPublic() bool {
	// this is another way of asking if a check
	// for read access is required
	return []bool{true, false}[nt]
}

func NodeTypeFromString(nt string) NodeType {
	switch nt {
	case SIGNAL.String():
		return SIGNAL
	case OTHER.String():
		return OTHER
	// case PUBLIC.String():
	//	return PUBLIC
	//case DOMAIN.String():
	//	return DOMAIN
	default:
		return SIGNAL
	}
}
