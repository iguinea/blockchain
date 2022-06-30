package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/iguinea/cryptodemo/pkgs/blockchainserver"
)

func init() {
	fmt.Println("\033[2J") // clear
	log.SetPrefix(fmt.Sprintf("%.16s ", "BlockChainServer        "))
	log.SetOutput(os.Stdout)
}

func main() {
	port := flag.Uint("port", 5000, "TCP port number to run server")
	flag.Parse()
	log.Printf("PortNumber to run BlockChainServer: %d", *port)

	bcs := blockchainserver.NewBlockChainServer(uint16(*port))
	bcs.Run()
}
