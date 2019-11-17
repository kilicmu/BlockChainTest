package main

func main() {
	bc := NewBlockChain("shen")
	cli := CIL{bc: bc}
	cli.Run()
}
