package main

import (
	"fmt"
)

func main() {
	blockchain := NewBlockChain()
	blockchain.AddBlock("小王给小红转100")
	blockchain.AddBlock("小红转给小王50")
	for i, val := range blockchain.blocks {
		fmt.Printf("========== 区块高度:%d ===========\n", i)
		fmt.Printf("区块的Hash: %x\n", val.Hash)
		fmt.Printf("区块的Data: %s\n", val.Data)
		fmt.Printf("区块前hash: %x\n", val.PrvHash)
	}
	// fmt.Println(block.Data...)
}
