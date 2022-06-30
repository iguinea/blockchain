package main

import (
	"fmt"
	"log"
	"os"
)

func init() {

	fmt.Println("\033[2J") // clear
	log.SetPrefix(fmt.Sprintf("%.16s ", "Wallet           "))
	log.SetOutput(os.Stdout)
}

func main() {
	/*
		walletM := wallet.NewWallet() // Wallet del minero
		walletA := wallet.NewWallet() // Wallet person A
		walletB := wallet.NewWallet() // Wallet person B

		//	Transaction from A to B
		t := transaction.NewTransaction(
			walletA.PrivateKey(),
			walletA.PublicKey(),
			walletA.BlockChainAddress(),
			walletB.BlockChainAddress(),
			1.0,
		)

		// Blockchain

			blockch := blockchain.NewBlockchain(walletM.BlockChainAddress())
			isAdded := blockch.AddTransaction(
				walletA.BlockChainAddress(),
				walletB.BlockChainAddress(),
				t.GetValue(),
				walletA.PublicKey(),
				t.GenerateSignature(),
			)
			log.Printf("Transaction added? %t", isAdded)
			blockch.Mining()
			blockch.Print()

			log.Printf("A: %.3f", blockch.CalculateTotalAmount(walletA.BlockChainAddress()))
			log.Printf("B: %.3f", blockch.CalculateTotalAmount(walletB.BlockChainAddress()))
			log.Printf("M: %.3f", blockch.CalculateTotalAmount(walletM.BlockChainAddress()))
	*/
}
