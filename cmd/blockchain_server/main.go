package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/iguinea/cryptodemo/pkgs/blockchainserver"
	"github.com/iguinea/cryptodemo/pkgs/utils"
)

var bcs *blockchainserver.BlockChainServer

func init() {
	fmt.Println("\033[2J") // clear
	log.SetPrefix(fmt.Sprintf("%.16s ", "BlockChainServer        "))
	log.SetOutput(os.Stdout)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGABRT, syscall.SIGALRM)
	go func() {
		for sig := range c {
			// sig is a ^C, handle it
			log.Printf("Received the following signal: %v", sig)
			if sig == os.Interrupt {
				log.Printf("Ctrl-C catch, terminating the process")
				bcs.Finish()
				os.Exit(1)
			}
		}
	}()
}

func main() {
	port := flag.Uint("port", 5000, "TCP port number to run server")
	flag.Parse()

	log.Printf("Host       to run BlockChainServer: %s", utils.GetHost())
	log.Printf("PortNumber to run BlockChainServer: %d", *port)

	bcs = blockchainserver.NewBlockChainServer(uint16(*port))
	bcs.Run()
}
