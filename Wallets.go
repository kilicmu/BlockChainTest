package main

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"io/ioutil"
	"log"
	"os"
)

//定义一个钱包结构体

type Wallets struct {
	WalletMap map[string]*Wallet
}

//定义创建函数

func NewWallets() *Wallets {
	var ws Wallets
	ws.loadFile()
	return &ws
}

func (ws *Wallets) CreateWallet() string {
	w := NewWallet()
	address := w.NewAddress()
	ws.WalletMap[address] = w
	ws.saveToFile()
	return address
}

func (ws *Wallets) loadFile() {
	//TODO 添加反序列化
	_, err := os.Stat("./wallets.dc")
	var tmp Wallets
	if err == nil {
		data, error := ioutil.ReadFile("wallets.dc")
		if error != nil {
			log.Panic(err)
		}
		gob.Register(elliptic.P256())
		decoder := gob.NewDecoder(bytes.NewReader(data))
		//解析出对象

		err = decoder.Decode(&tmp)
		if err != nil {
			log.Panic(err)
		}
	} else {
		os, _ := os.Create("wallets.dc")
		os.Close()
		tmp.WalletMap = make(map[string]*Wallet)
	}
	ws.WalletMap = tmp.WalletMap

}

func (ws *Wallets) saveToFile() {
	var buffer bytes.Buffer
	//gob编码的对象中存在接口需要在gob.register中将接口注册然后使用
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(ws)
	if err != nil {
		log.Panic(err)
	}
	err = ioutil.WriteFile("wallets.dc", buffer.Bytes(), 0600)
	if err != nil {
		log.Panic(err)
	}
}

func (ws *Wallets) GetAllAddresses() []string {
	var addresses []string
	for k, _ := range ws.WalletMap {
		addresses = append(addresses, k)
	}
	return addresses
}
