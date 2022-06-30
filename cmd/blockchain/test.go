package main

import (
	"fmt"
	"log"
	"os"
)

func init() {
	fmt.Println("\033[2J") // clear
	log.SetPrefix(fmt.Sprintf("%.16s ", "CryptoDemo  "))
	log.SetOutput(os.Stdout)
}

func main() {
	/*
		myBlockChainAddress := "my_block_chain_address"

		blockChain := blockchain.NewBlockchain(myBlockChainAddress)
		blockChain.Print()

		blockChain.AddTransaction("A", "B", 1.0)
		blockChain.Mining()
		blockChain.Print()

		blockChain.AddTransaction("C", "B", 3.0)
		blockChain.AddTransaction("D", "B", 44.0)
		blockChain.AddTransaction("Y", "X", 1.0)
		blockChain.Mining()
		blockChain.Print()

		log.Printf("My B                 = %.3f", blockChain.CalculateTotalAmount("B"))
		log.Printf("My BlockChainAddress = %.3f", blockChain.CalculateTotalAmount(myBlockChainAddress))
	*/
}
