package main

import (
	"fmt"
	"log"

	"github.com/bolt"
)

const DBname = "BlockDB.db"
const BucketName = "BlockBucket"
const LastHashKey = "LastHashKey"

//此处为区块链结构体定义
type BlockChain struct {
	//定义区块链数组
	db   *bolt.DB
	tail []byte //保存最后的hash值
}

func NewBlockChain(address string) *BlockChain {
	var LastHash []byte
	//初始化DB
	db, err := bolt.Open(DBname, 0600, nil)
	if err != nil {
		log.Panic("open db fail")
	}
	//对DB进行值更新
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))

		if bucket == nil {
			bucket, err = tx.CreateBucket([]byte(BucketName))
			if err != nil {
				log.Panic("create bucket err")
				return err
			}
			//加入创世区块的HASH(key)与他的序列化(value)
			genesisBlock := NewGenesisBlock(address)
			err := bucket.Put([]byte(genesisBlock.Hash), genesisBlock.Serialize())
			if err != nil {
				log.Panic("Put GenesisBlock err")
				return err
			}
			err = bucket.Put([]byte(LastHashKey), []byte(genesisBlock.Hash))
			LastHash = genesisBlock.Hash
			if err != nil {
				log.Panic("Put LastHash err")
				return err
			}
		} else {
			LastHash = bucket.Get([]byte(LastHashKey))

		}

		return nil
	})

	return &BlockChain{
		db:   db,
		tail: []byte(LastHash),
	}
}

func (bc *BlockChain) AddBlock(txs []*Transaction) {

	bc.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))
		if bucket == nil {
			log.Panic("add block no bucket")
			return nil
		} else {
			fmt.Println("挖矿ing...")
			b := NewBlock(txs, bc.tail)
			err := bucket.Put([]byte(b.Hash), []byte(b.Serialize()))
			if err != nil {
				log.Panic("put new block err")
				return err
			}
			bc.tail = b.Hash
			bucket.Put([]byte(LastHashKey), b.Hash)
		}
		return nil
	})
}

func (bc *BlockChain) FindUTXOs(address string) []TXOutput {
	var UTXOs []TXOutput
	spentOutputs := make(map[string][]int64)
	//遍历区块
	//遍历交易
	//遍历output找到与自己相关的utxo(添加output前检查一下是否消耗过)
	//遍历input(找到花过的utxo集合)
	it := bc.NewBlockChainIterator()
	for {
		block := it.Next()
		//此处开始遍历交易
		for _, tx := range block.Transactions {
			//遍历output, 找到与自己相关的utxo
		OutTag:
			for i, output := range tx.TXOutputs {

				if spentOutputs[string(tx.TXID)] != nil {
					for _, j := range spentOutputs[string(tx.TXID)] {
						if int64(i) == j {
							//说明这个utox已经被消耗过了
							goto OutTag
						}
					}
				}
				if output.PubKeyHash == address {
					UTXOs = append(UTXOs, output)
				}
			}
			//遍历input, 找到花费的utxo合集
			for _, input := range tx.TXInputs {
				if input.Sig == address {
					spentOutputs[string(input.TXid)] = append(spentOutputs[string(input.TXid)], input.Index)
				}
			}

		}

		if len(block.PrvHash) == 0 {
			break
		}
	}

	return UTXOs
}
