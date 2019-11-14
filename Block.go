package main

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"log"
	"time"
)

type Block struct {
	Version    uint64
	MerkelRoot []byte
	Data       []byte
	TimeStamp  uint64
	Difficulty uint64
	Nonce      uint64
	PrvHash    []byte
	Hash       []byte
}

func NewBlock(Data string, PrvHash []byte) *Block {
	block := Block{
		Version:    00,
		MerkelRoot: []byte{},
		TimeStamp:  uint64(time.Now().Unix()),
		Difficulty: 0,
		Nonce:      0,
		PrvHash:    PrvHash,
		Hash:       []byte{},
		Data:       []byte(Data),
	}
	pow := NewProofOfWork(&block)
	nonce, hash := pow.Run()
	block.Nonce = nonce
	block.Hash = hash
	return &block
}

func NewGenesisBlock() *Block {
	return NewBlock("这是创世块", []byte{})
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

//TODO
func Uint64ToByte(num uint64) []byte {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buf.Bytes()
}

//func (b *Block) SetHash() {
//	var BlockInfo []byte
//	//BlockInfo = append(BlockInfo, Uint64ToByte(b.Version)...)
//	//BlockInfo = append(BlockInfo, Uint64ToByte(b.Difficulty)...)
//	//BlockInfo = append(BlockInfo, Uint64ToByte(b.TimeStamp)...)
//	//BlockInfo = append(BlockInfo, Uint64ToByte(b.Nonce)...)
//	//BlockInfo = append(BlockInfo, b.MerkelRoot...)
//	//BlockInfo = append(BlockInfo, b.Hash...)
//	//BlockInfo = append(BlockInfo, b.Data...)
//	//BlockInfo = append(BlockInfo, b.PrvHash...)
//	tmp := [][]byte{
//		Uint64ToByte(b.Version),
//		Uint64ToByte(b.Difficulty),
//		Uint64ToByte(b.TimeStamp),
//		Uint64ToByte(b.Nonce),
//		b.MerkelRoot,
//		b.Hash,
//		b.Data,
//		b.PrvHash,
//	}
//	BlockInfo = bytes.Join(tmp, []byte{})
//
//	hash := sha256.Sum256(BlockInfo)
//	b.Hash = hash[:]
//}
