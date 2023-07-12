package main

import (
	// "fmt"
	"encoding/hex"
	"sync"
	"testing"
)

var (
	num  = uint64(1000)
	tStr = "test string"
	tBts = []byte("test string")

	once sync.Once
	db   *DB
	sc   ServerConf
)

func safeString2ByteSlice(v string) []byte { return []byte(tStr) }
func safeByteSlice2String(v []byte) string { return string(tBts) }

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
		GenRandStr(1)
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

func open() {
	once.Do(func() {
		db = MustOpenAndWrapDB("test.db")
		sc = ServerConf{&DBWithCache{DB: db, cache: make(map[string][]byte)}}
	})
}

func BenchmarkDBGetOrPut(b *testing.B) {
	open()

	for i := 0; i < b.N; i++ {
		// db.MustDo(db.GetOrPut, q)
		q := &Query{
			Bucket: ServerConf{}.ConfBucketName(),
			KV:     []Pair{NewPair(Port.Bytes(), Port.DefaultBytes())},
		}
		db.GetOrPut(q)
	}
}

func BenchmarkServerConfGetOrPutPort(b *testing.B) {
	open()
	// cleanup := func() { sc.DeleteDB() }
	// b.Cleanup(cleanup)
	for i := 0; i < b.N; i++ {
		sc.Port()
	}
}
