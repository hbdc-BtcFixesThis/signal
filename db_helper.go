package main

import (
	"errors"
	"os"
	"time"

	"path/filepath"

	bolt "go.etcd.io/bbolt"
)

type DBActions uint8

const (
	Get DBActions = iota
	Put
	GetOrPut
	GetPage

	DbTimeout = 2 * time.Second
)

func (dba DBActions) String() string { return [...]string{"Get", "Put", "GetOrPut"}[dba] }

type DB struct{ *bolt.DB }

func OpenDB(fp string) (*bolt.DB, error) {
	if len(fp) == 0 {
		return nil, bolt.ErrInvalid
	}
	if err := os.MkdirAll(filepath.Dir(fp), os.ModePerm); err != nil {
		return nil, err
	}
	return bolt.Open(fp, 0600, &bolt.Options{Timeout: DbTimeout})
}

func MustOpenDB(fp string) *bolt.DB {
	db, err := OpenDB(fp)
	if err != nil {
		panic(err)
	}
	return db
}

func MustOpenAndWrapDB(fp string) *DB {
	return &DB{MustOpenDB(fp)}
}

func (db *DB) DeletDB() error {
	return os.Remove(db.Path())
}

func (db *DB) CreateBucket(name []byte) error {
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(name)
		return err
	})
}

func (db *DB) Buckets() ([][]byte, error) {
	var buckets [][]byte
	return buckets, db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, _ *bolt.Bucket) error {
			buckets = append(buckets, name)
			return nil
		})
	})
}

// Delete a key from target bucket
func (db *DB) Delete(bucket []byte, key []byte) error {
	err := db.Update(func(tx *bolt.Tx) error {
		if b := tx.Bucket(bucket); b != nil {
			return b.Delete(key)
		}
		return bolt.ErrBucketNotFound
	})

	return err
}

type nextFunc func() ([]byte, []byte)

func (db DB) parsePageStart(kv map[string][]byte, c *bolt.Cursor) ([]byte, []byte, uint64, nextFunc) {
	start, startProvided := kv["start"]
	k, v := c.Last()
	if startProvided {
		delete(kv, "start")
		k, v = c.Seek(start)
	}

	size, sizeProvided := kv["size"]
	sizeInt := uint64(10)
	if sizeProvided {
		delete(kv, "size")
		sizeInt = Btoi(size)
	}

	direction, directionProvided := kv["direction"]
	next := c.Prev
	if directionProvided {
		if ByteSlice2String(direction) == "asc" {
			next = c.Next
			if !startProvided {
				k, v = c.First()
			}
		}
		delete(kv, "direction")
	}

	return k, v, sizeInt, next
}

func (db *DB) GetPage(bucket []byte, kv map[string][]byte) error {
	return db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket(bucket)
		if b == nil {
			return bolt.ErrBucketNotFound
		}

		c := b.Cursor()
		i := uint64(0)
		k, v, pageSize, next := db.parsePageStart(kv, c)
		for i < pageSize {
			kv[ByteSlice2String(k)] = v
			k, v = next()
			i += 1
		}

		return nil
	})
}

func (db *DB) Get(bucket []byte, kv map[string][]byte) error {
	return db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return bolt.ErrBucketNotFound
		}
		for k, _ := range kv {
			if result := b.Get(String2ByteSlice(k)); result != nil {
				kv[k] = result
			}
		}
		return nil
	})
}

func (db *DB) Put(bucket []byte, kv map[string][]byte) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return bolt.ErrBucketNotFound
		}
		for k, v := range kv {
			if err := b.Put(String2ByteSlice(k), v); err != nil {
				return err
			}
		}
		return nil
	})
}

func (db *DB) GetOrPut(bucket []byte, kv map[string][]byte) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return bolt.ErrBucketNotFound
		}

		for k, v := range kv {
			bk := String2ByteSlice(k)
			if dbVal := b.Get(bk); dbVal != nil {
				kv[k] = v
			} else {
				if err := b.Put(bk, v); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (db *DB) Do(action DBActions, bucket []byte, kv map[string][]byte) error {
	switch action {
	case Get:
		return db.Get(bucket, kv)
	case Put:
		return db.Put(bucket, kv)
	case GetOrPut:
		return db.GetOrPut(bucket, kv)
	case GetPage:
		return db.GetPage(bucket, kv)
	default:
		return errors.New("The db action specified is not implemented!")
	}
}

func (db *DB) MustDo(action DBActions, bucket []byte, result map[string][]byte) {
	err := db.Do(action, bucket, result)
	if err != nil {
		if err == bolt.ErrBucketNotFound {
			if err := db.CreateBucket(bucket); err != nil {
				panic(err)
			}
			if err = db.Do(action, bucket, result); err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
}
