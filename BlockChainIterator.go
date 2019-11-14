package main

import (
	"github.com/bolt"
	"log"
)

type BlockChainIterator struct {
	db                 *bolt.DB
	currentHashPointer []byte
}

//创建迭代器
//思想: 返回当前区块, 指针前移直到创世区块
func (bc *BlockChain) NewBlockChainIterator() *BlockChainIterator {
	return &BlockChainIterator{
		db:                 bc.db,
		currentHashPointer: bc.tail,
	}
}

func (it *BlockChainIterator) Next() *Block {
	var block Block
	it.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))
		if bucket == nil {
			log.Panic("bucket不可以为空, 请检查")
			return bolt.ErrBucketExists
		}
		tmp := bucket.Get(it.currentHashPointer)
		block = DeSerialize(tmp)
		it.currentHashPointer = block.PrvHash
		return nil
	})
	return &block
}
