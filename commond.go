package main

import "fmt"

import "time"

func (cil *CIL) ShowChain() {
	for it := cil.bc.NewBlockChainIterator(); len(it.currentHashPointer) != 0; {
		block := it.Next()
		fmt.Print("==================\n")
		fmt.Printf("当前hash值%x\n", block.Hash)
		fmt.Printf("前hash值: %x\n", block.PrvHash)
		fmt.Printf("时间戳: %v\n", time.Unix(int64(block.TimeStamp), 0).Format("2006-01-02 15:04:05"))
		fmt.Printf("难度值: %v\n", block.Difficulty)
		fmt.Printf("MerkelRoot: %x\n", block.MerkelRoot)
		fmt.Println("tx: ")
		for _, tx := range block.Transactions {
			fmt.Println("------------------------")
			fmt.Println("inputs: ")
			for _, input := range tx.TXInputs {
				fmt.Println(input.Sig)
			}
			fmt.Println("------------------------")
			fmt.Println("outputs: ")
			for _, output := range tx.TXOutputs {
				fmt.Println("value: ", output.Value)
				fmt.Println("to", output.PubKeyHash)
			}

		}
		fmt.Printf("工作量证明: %d\n", block.Nonce)
	}
}

//添加区块接口
func (cil *CIL) AddBlock(txs []*Transaction) {
	cil.bc.AddBlock(txs)
}

//用于查询地址的余额接口
func (cil *CIL) GetBlance(address string) {
	utxos := cil.bc.FindUTXOs(address)
	var value float64
	for _, utxo := range utxos {
		value += utxo.Value
	}
	fmt.Printf("address: %v has %v btc", address, value)
}

//用于调用转账功能的接口
func (cil *CIL) Send(from, to string, amount float64, miner, data string) {
	tx := NewTransaction(from, to, amount, cil.bc)
	if tx == nil {
		fmt.Println("")
		fmt.Println("转账失败")
		return
	}
	coinbase := NewCoinBase(miner, data)
	cil.AddBlock([]*Transaction{coinbase, tx})
	fmt.Println("转账成功")
}
