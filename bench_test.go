package main

import (
	// "fmt"
	"encoding/hex"
	"testing"
)

var num = 1000

func BenchmarkByteToInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Btoi(Itob(uint64(i)))
	}
}

func BenchmarkItob(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Itob(uint64(i))
	}
}

func BenchmarkGenRandStr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenRandStr(1)
	}
}

func BenchmarkSHA256(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SHA256([]byte("v"))
	}
}

func BenchmarkEncodeToString_hexPkg(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// compare the two
		hex.EncodeToString([]byte("v"))
	}
}

func BenchmarkEncodeToHexString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		EncodeToHexString([]byte("v"))
	}
}

var tStr = "test string"
var tBts = []byte("test string")

func safeString2ByteSlice(v string) []byte { return []byte(tStr) }
func safeByteSlice2String(v []byte) string { return string(tBts) }

func BenchmarkString2ByteSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		String2ByteSlice(tStr)
	}
}

func BenchmarkSafeString2ByteSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		safeString2ByteSlice(tStr)
	}
}

func BenchmarkByteSlice2String(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ByteSlice2String(tBts)
	}
}

func BenchmarkSafeByteSlice2String(b *testing.B) {
	for i := 0; i < b.N; i++ {
		safeByteSlice2String(tBts)
	}
}
