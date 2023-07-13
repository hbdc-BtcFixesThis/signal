package main

import (
	"os"

	"path/filepath"

	bolt "go.etcd.io/bbolt"
)

// You can rollback the transaction at any point by returning an error.
// All database operations are allowed inside a read-write transaction.
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

func (db *DB) DeleteDB() error {
	if err := db.Close(); err != nil {
		return err
	}
	return os.Remove(db.Path())
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

func (db *DB) CreateBucket(q *Query) error {
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(q.Bucket)
		return err
	})
}

func (db *DB) Get(q *Query) error {
	return db.View(func(tx *bolt.Tx) error {
		if b := tx.Bucket(q.Bucket); b != nil {
			for i, _ := range q.KV {
				if result := b.Get(q.KV[i].Key); result != nil {
					q.KV[i].Val = result
				}
			}
			return nil
		}
		return bolt.ErrBucketNotFound
	})
}

func (db *DB) Put(q *Query) error {
	return db.Update(func(tx *bolt.Tx) error {
		if b := tx.Bucket(q.Bucket); b != nil {
			for _, kv := range q.KV {
				if err := b.Put(kv.Key, kv.Val); err != nil {
					// Returns an error if the bucket was created from
					// a read-only transaction, if the key is blank, if
					// the key is too large, or if the value is too large.
					return err
				}
			}
			return nil
		}
		return bolt.ErrBucketNotFound
	})
}

func (db *DB) GetOrPut(q *Query) error {
	return db.Update(func(tx *bolt.Tx) error {
		if b := tx.Bucket(q.Bucket); b != nil {
			for i, _ := range q.KV {
				if v := b.Get(q.KV[i].Key); v != nil {
					q.KV[i].Val = v
				} else {
					if err := b.Put(q.KV[i].Key, q.KV[i].Val); err != nil {
						return err
					}
				}
			}
			return nil
		}
		return bolt.ErrBucketNotFound
	})
}

// Delete a key from target bucket
func (db *DB) Delete(q *Query) error {
	return db.Update(func(tx *bolt.Tx) error {
		if b := tx.Bucket(q.Bucket); b != nil {
			for _, kv := range q.KV {
				// :: From bbolt source ::
				// Delete removes a key from the bucket.
				// If the key does not exist then nothing is done and a nil error is returned.
				// Returns an error if the bucket was created from a read-only transaction.
				if err := b.Delete(kv.Key); err != nil {
					return err
				}
			}
			// You commit the transaction by returning nil at the end.
			return nil
		}
		return bolt.ErrBucketNotFound
	})
}

func (db *DB) GetPage(q *PageQuery) error {
	return db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		if b := tx.Bucket(q.Bucket); b != nil {
			c := b.Cursor()

			kv := q.SeekFrom(c)
			next := q.Direction(c)
			size := q.Size()

			for i := 0; i < size; i++ {
				if kv.Key == nil {
					// if k is nil after calling next we've reached the end
					return nil
				}
				q.KV[i] = kv
				kv = NewPair(next())
			}

			return nil
		}
		return bolt.ErrBucketNotFound
	})
}

type QueryFunc func(q *Query) error

func (db *DB) MustDo(do QueryFunc, query *Query) {
	if err := do(query); err != nil {
		if err == bolt.ErrBucketNotFound {
			if err := db.CreateBucket(query); err != nil {
				panic(err)
			} else {
				if err = do(query); err != nil {
					panic(err)
				}
			}
		} else {
			panic(err)
		}
	}
}
