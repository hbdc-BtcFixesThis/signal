package main

import (
	"fmt"

	"crypto/rand"
	"crypto/sha256"
	"math/big"
)

func GenRandStr(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}

func SHA256(str string) string {
	h := sha256.New()
	h.Write([]byte(str))
	hash := fmt.Sprintf("%x", h.Sum(nil))
	return hash
}
