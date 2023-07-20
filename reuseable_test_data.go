package main

import (
	"strconv"
)

var (
	num         = uint64(1000)
	fnum        = float64(10.00)
	fnumb       = F64tb(fnum)
	tStr        = "test string"
	tBts        = []byte("test string")
	tDBFileName = "test.db"

	TestDB *DB
	TestSC ServerConf
)

func safeString2ByteSlice(v string) []byte { return []byte(tStr) }
func safeByteSlice2String(v []byte) string { return string(tBts) }

func tDataKV(size int, randKey bool, randVal bool) []Pair {
	result := make([]Pair, size)
	var key, val []byte
	for i := 0; i < size; i++ {
		if randKey {
			key = MustGenRandBytes(10)
		} else {
			key = F64tb(float64(i))
		}

		if randVal {
			val = MustGenRandBytes(10)
		} else {
			val = String2ByteSlice(strconv.Itoa(i))
		}
		result[i] = NewPair(key, val) // GenRandBytes(10))
	}
	return result
}

func open() {
	TestDB = MustOpenAndWrapDB(tDBFileName)
	TestSC = ServerConf{
		&DBWithCache{
			DB:    TestDB,
			cache: make(map[string][]byte),
		},
	}
}

func tDataQ(size int, put bool, randKey bool, randVal bool) *Query {
	q := &Query{
		Bucket: []byte("bucket"),
		KV:     tDataKV(size, randKey, randVal),
	}
	if put {
		TestDB.MustDo(TestDB.Put, q)
	}
	return q
}

func tDataQueryRandKeys(size int, put bool) *Query {
	return tDataQ(size, put, true, false)
}

func tDataQueryRandVals(size int, put bool) *Query {
	return tDataQ(size, put, false, true)
}

func tDataQueryRandKV(size int, put bool) *Query {
	return tDataQ(size, put, true, true)
}

func tDataQuery(size int, put bool) *Query {
	return tDataQ(size, put, false, false)
}
