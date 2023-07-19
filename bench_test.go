package main

import (
	// "fmt"
	"encoding/hex"
	"testing"
)

func BenchmarkDBGetOrPut(b *testing.B) {
	open()

	for i := 0; i < b.N; i++ {
		// TestDB.MustDo(TestDB.GetOrPut, q)
		q := &Query{
			Bucket: ServerConf{}.ConfBucketName(),
			KV:     []Pair{NewPair(Port.Bytes(), Port.DefaultBytes())},
		}
		TestDB.GetOrPut(q)
	}
}

func BenchmarkServerConfGetOrPutPort(b *testing.B) {
	open()
	// cleanup := func() { TestSC.DeleteDB() }
	// b.Cleanup(cleanup)
	for i := 0; i < b.N; i++ {
		TestSC.Port()
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
		// compare the two
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
