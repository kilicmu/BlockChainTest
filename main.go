package main

func main() {
	bc := NewBlockChain("18S2hXh3dGQAwuM2fRjPeuxWvhzNcWSt7h")
	cli := CIL{bc: bc}
	cli.Run()
}
