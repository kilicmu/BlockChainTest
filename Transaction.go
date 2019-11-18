package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

const amount = 12.5

type Transaction struct {
	TXID      []byte
	TXInputs  []TXInput
	TXOutputs []TXOutput
}

type TXInput struct {
	TXid  []byte //引用utxo所在交易的id
	Index int64  //所消费utxo在output中的索引
	Sig   string
}

type TXOutput struct {
	Value      float64
	PubKeyHash string
}

//为交易设置TXID(交易hash)
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

//挖矿交易的创建
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

//普通交易的创建
func NewTransaction(from, to string, amount float64, bc *BlockChain) *Transaction {

	var inputs []TXInput
	var outputs []TXOutput
	utxos, factAmount := bc.FindNeedUTXOs(from, amount)
	if factAmount < amount {
		fmt.Println(factAmount)
		return nil
	}
	for key, outputs := range utxos {
		for _, i := range outputs {
			input := TXInput{TXid: []byte(key), Index: i, Sig: from}
			inputs = append(inputs, input)
		}
	}
	outputs = append(outputs, TXOutput{Value: amount, PubKeyHash: to})
	if factAmount > amount {
		outputs = append(outputs, TXOutput{Value: factAmount - amount, PubKeyHash: from})
	}
	tx := Transaction{[]byte{}, inputs, outputs}
	return &tx

}

//判断此交易是否是挖矿交易
func (tx *Transaction) IsCoinbase() bool {

	if len(tx.TXInputs) == 1 && len(tx.TXInputs[0].TXid) == 0 && tx.TXInputs[0].Index != -1 {
		return true
	}
	return false
}
