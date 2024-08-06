package main

import (
	// "fmt"
	"encoding/hex"
	"testing"
)

func BenchmarkServerConfGetOrPutPort(b *testing.B) {
	open()
	defer TestDB.DeleteDB()

	for i := 0; i < b.N; i++ {
		TestSC.Port(nil)
	}
}

func BenchmarkServerConfMustGetOrPut(b *testing.B) {
	open()
	defer TestDB.DeleteDB()

	q := &Query{
		Bucket: tBucket,
		KV:     []Pair{NewPair(Port.Bytes(), Port.DefaultBytes())},
	}
	for i := 0; i < b.N; i++ {
		TestSC.MustDo(TestSC.DB.GetOrPut, q)
	}
}

func BenchmarkServerConfGetOrPut(b *testing.B) {
	open()
	defer TestDB.DeleteDB()

	q := &Query{
		Bucket:                  tBucket,
		KV:                      []Pair{NewPair(Port.Bytes(), Port.DefaultBytes())},
		CreateBucketIfNotExists: true,
	}
	for i := 0; i < b.N; i++ {
		TestSC.DB.GetOrPut(q)
	}
}

func BenchmarkF64tb(b *testing.B) {
	for i := 0; i < b.N; i++ {
		F64tb(fnum)
	}
}

func BenchmarkF64fb(b *testing.B) {
	for i := 0; i < b.N; i++ {
		F64fb(fnumb)
	}
}

func BenchmarkInToByteAndBack(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Btoi(Itob(num))
	}
}

func BenchmarkItob(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Itob(num)
	}
}

func BenchmarkGenRandStr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenRandBytes(1)
	}
}

func BenchmarkSHA256(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SHA256(tBts)
	}
}

func BenchmarkEncodeToString_hexPkg(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hex.EncodeToString(tBts)
	}
}

func BenchmarkEncodeToHexString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		EncodeToHexString(tBts)
	}
}

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
