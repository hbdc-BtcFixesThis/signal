package main

import (
	"os"

	"path/filepath"

	bolt "go.etcd.io/bbolt"
)

type DataUpdates struct {
	Put    []*Query
	Delete []*Query
}

func (du *DataUpdates) AddDeleteQuery(q *Query)     { du.Delete = append(du.Delete, q) }
func (du *DataUpdates) AddDeleteQueries(q []*Query) { du.Delete = append(du.Delete, q...) }
func (du *DataUpdates) AddPutQuery(q *Query)        { du.Put = append(du.Put, q) }
func (du *DataUpdates) AddPutQueries(q []*Query)    { du.Put = append(du.Put, q...) }

func (du *DataUpdates) AppendUpdates(other *DataUpdates) {
	du.AddDeleteQueries(other.Delete)
	du.AddPutQueries(other.Put)
}

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
	path := db.Path()
	if err := db.Close(); err != nil {
		return err
	}
	return os.Remove(path)
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

func getBucket(q *Query, tx *bolt.Tx) (*bolt.Bucket, error) {
	var b *bolt.Bucket
	if !q.CreateBucketIfNotExists || !tx.Writable() {
		if b = tx.Bucket(q.Bucket); b == nil {
			return nil, bolt.ErrBucketNotFound
		}
		return b, nil
	}
	return tx.CreateBucketIfNotExists(q.Bucket)
}

func (db *DB) CreateBucket(q *Query) error {
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(q.Bucket)
		return err
	})
}

func (db *DB) GenID(q *Query) ([]byte, error) {
	var id uint64
	return Itob(id), db.Update(func(tx *bolt.Tx) error {
		b, err := getBucket(q, tx)
		if err != nil {
			return err
		}
		id, err = b.NextSequence()
		if err != nil {
			return err
		}
		return nil
	})
}

func (db *DB) Get(q *Query) error {
	return db.Update(func(tx *bolt.Tx) error {
		b, err := getBucket(q, tx)
		if err != nil {
			return err
		}
		for i, _ := range q.KV {
			if result := b.Get(q.KV[i].Key); result != nil {
				q.KV[i].Val = result
			}
		}
		return nil
	})
}

func (db *DB) putTx(q *Query, tx *bolt.Tx) error {
	b, err := getBucket(q, tx)
	if err != nil {
		return err
	}
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

func (db *DB) Put(q *Query) error {
	return db.Update(func(tx *bolt.Tx) error { return db.putTx(q, tx) })
}

func (db *DB) GetOrPut(q *Query) error {
	return db.Update(func(tx *bolt.Tx) error {
		b, err := getBucket(q, tx)
		if err != nil {
			return err
		}
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
	})
}

func (db *DB) deleteTx(q *Query, tx *bolt.Tx) error {
	b, err := getBucket(q, tx)
	if err != nil {
		return err
	}
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

func (db *DB) Delete(q *Query) error {
	return db.Update(func(tx *bolt.Tx) error { return db.deleteTx(q, tx) })
}

func (db *DB) MultiWrite(updates *DataUpdates) error {
	return db.Update(func(tx *bolt.Tx) error {
		for i := 0; i < len(updates.Delete); i++ {
			if err := db.deleteTx(updates.Delete[i], tx); err != nil {
				return err
			}
		}
		for i := 0; i < len(updates.Put); i++ {
			if putErr := db.putTx(updates.Put[i], tx); putErr != nil {
				return putErr
			}
		}
		return nil
	})
}

func (db *DB) GetPage(pq *PageQuery) error {
	return db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b, err := getBucket(pq.Query, tx)
		if err != nil {
			return err
		}
		c := b.Cursor()

		kv := pq.SeekFrom(c)
		next := pq.Direction(c)
		size := pq.Size()

		for i := 0; i < size; i++ {
			if kv.Key == nil {
				// if k is nil after calling next we've reached the end
				return nil
			}
			pq.KV[i] = kv
			kv = NewPair(next())
		}

		return nil
	})
}

type QueryFunc func(q *Query) error

func (db *DB) MustDo(do QueryFunc, query *Query) {
	query.CreateBucketIfNotExists = true
	if err := do(query); err != nil {
		panic(err)
	}
}
