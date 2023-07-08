package main

import (
	"crypto/rand"
	"crypto/sha256"
	"math/big"
)

func GenRandStr(n int) ([]byte, error) {
	const (
		choices    = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
		numChoices = int64(len(choices))
	)

	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(numChoices))
		if err != nil {
			return []byte{}, err
		}
		ret[i] = choices[num.Int64()]
	}

	return ret, nil
}

func SHA256(v []byte) string {
	h := sha256.New()
	h.Write(v)
	return EncodeToHexString(h.Sum(nil))
}
