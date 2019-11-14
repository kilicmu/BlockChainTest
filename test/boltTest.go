package main

import (
	"fmt"
	"github.com/bolt"
	"log"
)

func main() {
	//open db
	db, err := bolt.Open("test.db", 0600, nil)
	if err != nil {
		log.Panic("open failed")
		return
	}
	// update db & do sth
	db.Update(func(tx *bolt.Tx) error {
		// init bucket
		bucket := tx.Bucket([]byte("b1"))
		if bucket == nil {
			bucket, err = tx.CreateBucket([]byte("b1"))
			if err != nil {
				log.Panic(bucket, "create db err")
			}
		}
		bucket.Put([]byte("test"), []byte("hahah"))
		//bucket.Put([]byte("test2"), []byte("2"))

		return nil
	})

	db.View(func(fn *bolt.Tx) error {
		bucket := fn.Bucket([]byte("b1"))
		resp := bucket.Get([]byte("test"))
		fmt.Print(string(resp))
		return nil
	})
}
