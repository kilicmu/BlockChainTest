package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"log"
	"math/big"
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
	Signature []byte //共识: 数字签名, r, s组成的[]byte,实际不是这样,这麽做为了简化流程
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
	privateKey := wallet.PrivateKey

	var inputs []TXInput
	var outputs []TXOutput
	utxos, factAmount := bc.FindNeedUTXOs(pubKeyHash, amount)
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

	bc.SignTransaction(*privateKey, &tx)

	return &tx

}

//判断此交易是否是挖矿交易
func (tx *Transaction) IsCoinbase() bool {
	if len(tx.TXInputs) == 1 && len(tx.TXInputs[0].TXid) == 0 && tx.TXInputs[0].Index == -1 {
		return true
	}
	return false
}

func (tx *Transaction) Sign(privateKey *ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	if tx.IsCoinbase() {
		return
	}
	txCopy := tx.TrimmedCopy()
	for i, input := range txCopy.TXInputs {
		//循环遍历txCopy的inputs,得到对应的output的公钥hash
		prevTX := prevTXs[string(input.TXid)]
		if len(prevTX.TXID) == 0 {
			log.Panic("交易不存在啊")
		}
		//要对txCopy中相应input的putkey进行赋值(为上一交易的output的PubkeyHash)
		txCopy.TXInputs[i].PubKey = prevTX.TXOutputs[input.Index].PubKeyHash
		//生成要签名的对象(TXID)
		txCopy.SetHash()
		txCopy.TXInputs[i].PubKey = nil
		signDataHash := txCopy.TXID
		//引用椭圆曲线算法对数据进行签名, 得到rs字节slice,存入相应的TXinput的Signature位置
		r, s, err := ecdsa.Sign(rand.Reader, privateKey, signDataHash)
		if err != nil {
			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)
		tx.TXInputs[i].Signature = signature
	}
}

//进行修建的拷贝,去掉签名位和公钥位的值
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TXInput
	var outputs []TXOutput
	for _, input := range tx.TXInputs {
		inputs = append(inputs, TXInput{input.TXid, input.Index, nil, nil})
	}
	for _, output := range tx.TXOutputs {
		outputs = append(outputs, output)
	}
	return Transaction{tx.TXID, inputs, outputs}
}

//需要数据(TXID)和公钥和签名
func (tx *Transaction) Verify(prevTX map[string]Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}
	txCopy := tx.TrimmedCopy()
	for i, input := range tx.TXInputs {
		//循环遍历tx的inputs,得到对应的output的公钥hash
		prevTX := prevTX[string(input.TXid)]
		if len(prevTX.TXID) == 0 {
			log.Panic("交易不存在啊")
		}
		//要对txCopy中相应input的putkey进行赋值(为上一交易的output的PubkeyHash)
		txCopy.TXInputs[i].PubKey = prevTX.TXOutputs[input.Index].PubKeyHash
		//生成要签名的对象(TXID)
		txCopy.SetHash()
		txCopy.TXInputs[i].PubKey = nil
		signDataHash := txCopy.TXID
		signature := input.Signature
		pubKey := input.PubKey
		r := big.Int{}
		s := big.Int{}
		r.SetBytes(signature[:len(signature)/2])
		s.SetBytes(signature[len(signature)/2:])
		X := big.Int{}
		Y := big.Int{}
		X.SetBytes(pubKey[:len(pubKey)/2])
		Y.SetBytes(pubKey[len(pubKey)/2:])
		pubKeyOrigin := ecdsa.PublicKey{elliptic.P256(), &X, &Y}
		if !ecdsa.Verify(&pubKeyOrigin, signDataHash, &r, &s) {
			return false
		}

	}
	return true
}
