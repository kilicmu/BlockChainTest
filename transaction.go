package main

import "bytes"

import "encoding/gob"

import "log"

import "crypto/sha256"

const amount = 12.5

type Transaction struct {
	TXID      []byte
	TXInputs  []TXInput
	TXOutputs []TXOutput
}

type TXInput struct {
	TXid  []byte 	//引用utxo所在交易的id
	Index int64		//所消费utxo在output中的索引
	Sig   string
}

type TXOutput struct {
	Value      float64
	PubKeyHash string
}

func (tx *Transaction) SetHash() {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic("transaction set hash err!")
		return
	}
	data := buffer.Bytes()
	hash := sha256.Sum256(data)
	tx.TXID = hash[:]
}

func NewCoinBase(address string, data string) *Transaction {
	input := TXInput{TXid: []byte{}, Index: -1, Sig: data}
	output := TXOutput{Value: amount, PubKeyHash: address}
	tx := Transaction{
		[]byte{},
		[]TXInput{input},
		[]TXOutput{output},
	}
	tx.SetHash()
	return &tx

}
