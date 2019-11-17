package main

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"log"
	"time"
)

const GenesisMessage = "创世区块信息"

type Block struct {
	Version    uint64
	MerkelRoot []byte
	// Data       []byte
	Transactions []*Transaction
	TimeStamp    uint64
	Difficulty   uint64
	Nonce        uint64
	PrvHash      []byte
	Hash         []byte
}

func NewBlock(txs []*Transaction, PrvHash []byte) *Block {
	block := Block{
		Version:      00,
		MerkelRoot:   []byte{},
		TimeStamp:    uint64(time.Now().Unix()),
		Difficulty:   0,
		Nonce:        0,
		PrvHash:      PrvHash,
		Hash:         []byte{},
		Transactions: txs,
	}
	//只对区块头做hash区块体通过影响MerkelRoot决定区块的最终hash结果
	block.MakeMelRoot()
	pow := NewProofOfWork(&block)
	nonce, hash := pow.Run()
	block.Nonce = nonce
	block.Hash = hash
	return &block
}

func NewGenesisBlock(address string) *Block {
	tx := NewCoinBase(address, GenesisMessage)
	return NewBlock([]*Transaction{tx}, []byte{})
}

//区块的序列化
func (b *Block) Serialize() []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(&b)
	if err != nil {
		log.Panic("encode err")
		return []byte{}
	}
	return buffer.Bytes()
}

//区块的反序列化
func DeSerialize(buffer []byte) Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(buffer))
	err := decoder.Decode(&block)
	if err != nil {

		return block
	}
	return block
}

func Uint64ToByte(num uint64) []byte {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buf.Bytes()
}

func (b *Block) MakeMelRoot() {
	//TODO
	b.MerkelRoot = []byte{}
}
