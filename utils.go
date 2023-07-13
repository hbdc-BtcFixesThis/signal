package main

import (
	"unsafe"

	"encoding/binary"
	"encoding/hex"
	"os/user"
)

func ConcatSlice[T any](first []T, second []T) []T {
	n := len(first)
	return append(first[:n:n], second...)
}

func Btoi(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

func Itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

func Contains(sl []string, name string) bool {
	for _, value := range sl {
		if value == name {
			return true
		}
	}
	return false
}

func GetCurrentUserHomeDir() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}
	return currentUser.HomeDir, nil
}

func EncodeToHexString(src []byte) string {
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)
	return ByteSlice2String(dst)
}

func String2ByteSlice(str string) []byte {
	if str == "" {
		return nil
	}
	return unsafe.Slice(unsafe.StringData(str), len(str))
}

func ByteSlice2String(bs []byte) string {
	if len(bs) == 0 {
		return ""
	}
	return unsafe.String(unsafe.SliceData(bs), len(bs))
}

/*
// ^^ careful when using unsafe. Efficiency on this front comes with some gotchyas
// example of failing test:
func TestByteSlice2String(t *testing.T) {
	// given
	sampleBytes1 := []byte{0x30, 0x31, 0x32}
	// when
	sampleStr1 := ByteSlice2String(sampleBytes1)
	// then
	assert.Equal(t, "012", sampleStr1)

	// when
	sampleBytes1[1] = 0x39
	// then
	assert.Equal(t, "012", sampleStr1)
}

--- FAIL: TestByteSlice2String (0.00s)
Error:       Not equal:
             expected: "012"
             actual  : "092"
*/
