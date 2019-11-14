package main

import "fmt"

func main() {
	blockchain := NewBlockChain()
	// fmt.Printf("%x\n", blockchain.tail)
	blockchain.AddBlock("小王给小红转100")
	blockchain.AddBlock("小红转给小王50")
	for it := blockchain.NewBlockChainIterator(); len(it.currentHashPointer) != 0; {
		block := it.Next()
		fmt.Print("==================\n")
		fmt.Printf("当前hash值%x\n", block.Hash)
		fmt.Printf("前hash值: %x\n", block.PrvHash)
		fmt.Printf("当前的区块数据: %s\n", block.Data)
		fmt.Printf("工作量证明: %d\n", block.Nonce)
	}
}
