package main

import "crypto/ecdsa"

import "crypto/elliptic"

import "crypto/rand"

import "log"

import "crypto/sha256"

import "math/big"

import "fmt"

//
func main() {
	//创建曲线
	curve := elliptic.P256()
	//根据一个随机数和曲线创建私钥
	privatekey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic()
	}
	//根据私钥找公钥
	pubKey := privatekey.PublicKey
	data := "hello"
	hash := sha256.Sum256([]byte(data))
	//数据的hash和私钥进行签名
	r, s, err := ecdsa.Sign(rand.Reader, privatekey, hash[:])
	fmt.Println(r)
	fmt.Println(s)
	if err != nil {
		return
	}

	//进行序列化传输
	signature := append(r.Bytes(), s.Bytes()...)

	//校验需要三个东西: 数据, 签名, 公钥
	//在本地对数据流进行拆分校验
	//对r,s进行还原
	r1 := big.Int{}
	s1 := big.Int{}
	r1.SetBytes(signature[:len(signature)/2])
	s1.SetBytes(signature[len(signature)/2:])
	//对还原后内容进行校验
	res := ecdsa.Verify(&pubKey, hash[:], &r1, &s1)
	//return true
	fmt.Println(res)
}
