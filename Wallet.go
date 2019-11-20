package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
	"log"
)

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  []byte
}

func NewWallet() *Wallet {
	var w Wallet
	curve := elliptic.P256()
	privatekey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic()
	}

	w.PrivateKey = privatekey
	w.PublicKey = append(w.PublicKey, privatekey.PublicKey.X.Bytes()...)
	w.PublicKey = append(w.PublicKey, privatekey.PublicKey.Y.Bytes()...)
	return &w
}

//返回需要的地址
func (w *Wallet) NewAddress() string {
	pubKey := w.PublicKey
	rip160hasherValue := HashPubKey(pubKey)
	version := []byte{0}
	fmt.Println(rip160hasherValue)
	payload := append(version, rip160hasherValue...)
	fmt.Println(payload)
	checkcode := GetCheckCode(payload)
	payload = append(payload, checkcode...)

	address := base58.Encode(payload)

	return address

}

//对公钥进行hash运算, 将运算切片返回
func HashPubKey(data []byte) []byte {
	hash := sha256.Sum256(data)
	//新建一个ripemd160编码器
	rip160hasher := ripemd160.New()
	_, err := rip160hasher.Write(hash[:])
	if err != nil {
		log.Panic()
	}

	rip160hasherValue := rip160hasher.Sum(nil)
	return rip160hasherValue
}

//获取校验值
func GetCheckCode(payload []byte) []byte {
	hash1 := sha256.Sum256(payload)
	hash2 := sha256.Sum256(hash1[:])
	checkcode := hash2[:4]
	return checkcode
}

//传入一个地址, 判断这个地址是否有效
func IsValidAddress(address string) bool {
	addressByte := base58.Decode(address)
	payload := addressByte[:len(addressByte)-4]
	checksum1 := addressByte[len(addressByte)-4:]
	checksum2 := GetCheckCode(payload)
	return bytes.Equal(checksum1, checksum2)
}
