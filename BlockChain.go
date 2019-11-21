package main

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
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
	for _, tx := range txs {
		if !bc.VerifyTransaction(tx) {
			fmt.Println("矿工发现无效交易")
			return
		}
	}
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

//找到所有当前地址的UTXO
func (bc *BlockChain) FindUTXOs(pubKeyHash []byte) []TXOutput {
	var UTXOs []TXOutput
	txs := bc.FindAllTxs(pubKeyHash)
	for _, tx := range txs {
		for _, output := range tx.TXOutputs {
			if bytes.Equal(output.PubKeyHash, pubKeyHash) {
				UTXOs = append(UTXOs, output)
			}
		}
	}
	return UTXOs
}

//传入支付地址与需要支付的价格, 返回一个需要使用的output的map[交易id]交易output标志,与实际找到的金额
func (bc *BlockChain) FindNeedUTXOs(senderPubKeyHash []byte, need float64) (map[string][]int64, float64) {
	UTXOs := make(map[string][]int64)
	var amount float64 = 0.00
	txs := bc.FindAllTxs(senderPubKeyHash)
	for _, tx := range txs {
		for i, output := range tx.TXOutputs {
			//两个byte数组比较方式bytes.Equal
			if bytes.Equal(output.PubKeyHash, senderPubKeyHash) {
				if amount < need {
					UTXOs[string(tx.TXID)] = append(UTXOs[string(tx.TXID)], int64(i))
					amount += output.Value
					if amount > need {
						return UTXOs, amount
					}
				}

			}
		}
	}
	return UTXOs, -1
}

//返回包含一个地址交易的所有交易
func (bc *BlockChain) FindAllTxs(senderPubKeyHash []byte) []*Transaction {
	//TODO  优化查询算法
	var txs []*Transaction
	it := bc.NewBlockChainIterator()
	spentOutputs := make(map[string][]int64)
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
							continue OutTag
						}
					}
				}
				if bytes.Equal(output.PubKeyHash, senderPubKeyHash) {
					txs = append(txs, tx)
				}
			}
			//遍历input, 找到花费的utxo合集
			if !tx.IsCoinbase() {
				for _, input := range tx.TXInputs {
					hashPubKey := HashPubKey(input.PubKey)
					if bytes.Equal(hashPubKey, senderPubKeyHash) {
						spentOutputs[string(input.TXid)] = append(spentOutputs[string(input.TXid)], input.Index)
					}
				}
			}

		}

		if len(block.PrvHash) == 0 {
			break
		}
	}

	return txs
}

func (bc *BlockChain) FindTxByID(id []byte) (Transaction, error) {
	it := bc.NewBlockChainIterator()
	for {
		block := it.Next()
		for _, tx := range block.Transactions {
			if bytes.Equal(tx.TXID, id) {
				return *tx, nil
			}
		}

		if block.PrvHash == nil {
			break
		}
	}
	return Transaction{}, errors.New("无效交易ID")
}

func (bc *BlockChain) SignTransaction(privateKey ecdsa.PrivateKey, tx *Transaction) {
	prevTXs := make(map[string]Transaction)
	for _, input := range tx.TXInputs {
		tx, err := bc.FindTxByID(input.TXid)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[(string(input.TXid))] = tx
	}
	tx.Sign(&privateKey, prevTXs)
}

func (bc *BlockChain) VerifyTransaction(tx *Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}
	prevTX := make(map[string]Transaction)
	for _, input := range tx.TXInputs {
		tx, err := bc.FindTxByID(input.TXid)
		if err != nil {
			log.Panic(err)
		}
		prevTX[string(input.TXid)] = tx
	}
	return tx.Verify(prevTX)
}
