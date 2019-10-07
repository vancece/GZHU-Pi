package bolt

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/boltdb/bolt"
	"time"
)

type BlotClient struct {
	DB     *bolt.DB
	Bucket string
}

func newBlotClient(path, bucket string) (b *BlotClient) {
	if path == "" || bucket == "" {
		return nil
	}
	b = &BlotClient{Bucket: bucket}
	var err error
	b.DB, err = bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		logs.Error(err)
		return nil
	}
	return
}

func (bdb *BlotClient) Put(key, value []byte) (err error) {
	if key == nil || value == nil {
		return fmt.Errorf("the key or value is null")
	}
	err = bdb.DB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bdb.Bucket))
		if err != nil {
			return fmt.Errorf("create bucket failed: %s", err)
		}
		b = tx.Bucket([]byte(bdb.Bucket))
		err = b.Put(key, value)
		if err != nil {
			logs.Error("put data failed: ", err)
			return err
		}
		return nil
	})
	return
}

func (bdb *BlotClient) Get(key string) (value []byte, err error) {
	if bdb.Bucket == "" || key == "" {
		err = fmt.Errorf("invalid bucket or key")
		return
	}
	err = bdb.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bdb.Bucket))
		value = b.Get([]byte(key))
		return nil
	})
	return
}
