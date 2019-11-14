package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

type ProofOfWork struct {
	b      *Block
	target *big.Int
}

func NewProofOfWork(b *Block) *ProofOfWork {

	pow := ProofOfWork{b: b}
	temp := big.Int{}
	temp.SetString("0000000000000100000001000001000000000000000000000000000000000000000000000", 16)
	pow.target = &temp
	return &pow
}

func (pow *ProofOfWork) Run() (uint64, []byte) {
	b := pow.b
	nonce := uint64(00)
	for {
		tmp := [][]byte{
			Uint64ToByte(b.Version),
			Uint64ToByte(b.Difficulty),
			Uint64ToByte(b.TimeStamp),
			Uint64ToByte(nonce),
			b.MerkelRoot,
			b.Hash,
			b.Data,
			b.PrvHash,
		}
		BlockInfo := bytes.Join(tmp, []byte{})
		hash := sha256.Sum256(BlockInfo)
		tmpBigInt := big.Int{}
		tmpBigInt.SetBytes(hash[:])

		if tmpBigInt.Cmp(pow.target) == -1 {
			fmt.Printf("get hash: %x nonce: %d\n", hash, nonce)
			return nonce, hash[:]
		} else {
			nonce++
			//fmt.Print(nonce)
		}
	}

}
