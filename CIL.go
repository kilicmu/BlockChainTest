package main

import (
	"fmt"
	"os"
	"strconv"
)

var useinfo string = `
list				打印链
search --address ADDRESS	查询address余额
send FROM to TO PRICE SOMEONE Data	发送货币从地址1到地址2, 其中最后指定矿工打包交易
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
	case "list":
		cil.ShowChain()
	case "search":
		if len(args) == 4 && args[2] == "--address" {
			cil.GetBlance(args[3])
		}
	case "send":
		if len(args) == 8 && args[3] == "to" {
			//send FROM to TO PRICE SOMEONE
			from := args[2]
			to := args[4]
			amount, _ := strconv.ParseFloat(args[5], 64)
			miner := args[6]
			data := args[7]
			cil.Send(from, to, amount, miner, data)
		}
	default:
		fmt.Println("参数错误")
		fmt.Println(useinfo)
	}
	return nil
}
