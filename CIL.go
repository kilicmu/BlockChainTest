package main

import (
	"fmt"
	"os"
)

var useinfo string = `
add --data datamsg     添加新区块
list			打印链
search --address ADDRESS	查询address余额
`

type CIL struct {
	bc *BlockChain
}

func (cil *CIL) Run() error {
	args := os.Args
	if len(args) < 2 {
		fmt.Println(useinfo)
		return nil
	}
	cmd := args[1]
	switch cmd {

	case "add":
		//添加区块
		if len(args) == 4 && args[2] == "--data" {
			// data := args[3]
			cil.AddBlock([]*Transaction{})
		} else {
			fmt.Println("参数错误")
			fmt.Println(useinfo)
		}
	case "list":
		cil.ShowChain()
	case "search":
		if len(args) == 4 && args[2] == "--address" {
			cil.GetBlance(args[3])
		}
	default:
		fmt.Println("...")
	}
	return nil
}
