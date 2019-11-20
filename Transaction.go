package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"log"
)

const reward = 12.5

type Transaction struct {
	TXID      []byte
	TXInputs  []TXInput
	TXOutputs []TXOutput
}

type TXInput struct {
	TXid      []byte //引用utxo所在交易的id
	Index     int64  //所消费utxo在output中的索引
	Signature []byte //数字签名, r, s组成的[]byte
	PubKey    []byte
}

type TXOutput struct {
	Value float64
	// PubKeyHash string
	PubKeyHash []byte //是 公钥不是哈希或地址
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

func (output *TXOutput) Lock(address string) {
	data := GetPubKeyHash(address)
	output.PubKeyHash = data
}

func GetPubKeyHash(address string) []byte {
	data := base58.Decode(address)
	len := len(data)
	data = data[1 : len-4]
	return data
}

func NewTXOutput(value float64, address string) *TXOutput {
	output := TXOutput{
		Value: value,
	}
	output.Lock(address)
	return &output
}

//挖矿交易的创建
func NewCoinBase(address string, data string) *Transaction {
	input := TXInput{TXid: []byte{}, Index: -1, Signature: nil, PubKey: []byte(data)}
	output := NewTXOutput(reward, address)
	tx := Transaction{
		[]byte{},
		[]TXInput{input},
		[]TXOutput{*output},
	}
	tx.SetHash()
	return &tx
}

//普通交易的创建
func NewTransaction(from, to string, amount float64, bc *BlockChain) *Transaction {
	ws := NewWallets()
	wallet := ws.WalletMap[from]
	if wallet == nil {
		fmt.Println("钱包不存在")
		return nil
	}
	pubKey := wallet.PublicKey
	//需要传递的是公钥的哈希
	pubKeyHash := HashPubKey(pubKey)

	var inputs []TXInput
	var outputs []TXOutput
	utxos, factAmount := bc.FindNeedUTXOs(pubKeyHash, amount)
	fmt.Println("factAmount: ", factAmount)
	if factAmount < amount {
		fmt.Println("余额不足")
		return nil
	}
	for key, outputs := range utxos {
		for _, i := range outputs {
			input := TXInput{TXid: []byte(key), Index: i, Signature: nil, PubKey: pubKey}
			inputs = append(inputs, input)
		}
	}
	outputs = append(outputs, *NewTXOutput(amount, to))
	if factAmount > amount {

		outputs = append(outputs, *NewTXOutput(factAmount-amount, from))
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
