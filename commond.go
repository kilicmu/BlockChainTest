package main

import "fmt"

func (cil *CIL) ShowChain() {
	for it := cil.bc.NewBlockChainIterator(); len(it.currentHashPointer) != 0; {
		block := it.Next()
		fmt.Print("==================\n")
		fmt.Printf("当前hash值%x\n", block.Hash)
		fmt.Printf("前hash值: %x\n", block.PrvHash)
		fmt.Printf("当前的区块数据: %s\n", block.Transactions[0].TXInputs[0].Sig)
		fmt.Printf("工作量证明: %d\n", block.Nonce)
	}
}

func (cil *CIL) AddBlock(txs []*Transaction) {
	cil.bc.AddBlock(txs)
}

func (cil *CIL) GetBlance(address string) {
	utxos := cil.bc.FindUTXOs(address)
	var value float64
	for _, utxo := range utxos {
		value += utxo.Value
	}
	fmt.Printf("address: %v has %v btc", address, value)
}
