package main

import (
	"strconv"
	"sync"
)

var (
	num         = uint64(1000)
	tStr        = "test string"
	tBts        = []byte("test string")
	tDBFileName = "test.db"

	once   sync.Once
	TestDB *DB
	TestSC ServerConf
)

func safeString2ByteSlice(v string) []byte { return []byte(tStr) }
func safeByteSlice2String(v []byte) string { return string(tBts) }

func tData(size int) []Pair {
	result := make([]Pair, size)
	for i := 0; i < size; i++ {
		iAsBytes := String2ByteSlice(strconv.Itoa(i))
		result[i] = NewPair(iAsBytes, iAsBytes) // GenRandBytes(10))
	}
	return result
}

func open() {
	once.Do(func() {
		TestDB = MustOpenAndWrapDB(tDBFileName)
		TestSC = ServerConf{
			&DBWithCache{
				DB:    TestDB,
				cache: make(map[string][]byte),
			},
		}
	})
}
