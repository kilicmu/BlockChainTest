package main



type BlockChain struct {
	//定义区块链数组
	blocks []*Block
}

func NewGenesisBlock() *Block {
	return NewBlock("这是创世块", []byte{})
}

func NewBlockChain() *BlockChain {
	genesisBlock := NewGenesisBlock()
	return &BlockChain{
		blocks: []*Block{genesisBlock},
	}
}

func(bc *BlockChain) AddBlock(data string){
	b := NewBlock(data, bc.blocks[len(bc.blocks)-1].Hash)
	bc.blocks = append(bc.blocks, b)

}
