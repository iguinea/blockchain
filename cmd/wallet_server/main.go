package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/iguinea/cryptodemo/pkgs/walletserver"
)

func init() {
	fmt.Println("\033[2J") // clear
	log.SetPrefix(fmt.Sprintf("%.16s ", "WalletServer        "))
	log.SetOutput(os.Stdout)
}

func main() {
	port := flag.Uint("port", 8080, "TCP port number for wallet server")
	gateway := flag.String("gateway", "http://0.0.0.0:5000", "Blockchain gateway")
	flag.Parse()

	log.Printf("PortNumber to run WalletServer: %d", *port)
	log.Printf("Gateway to connect to : %s", *gateway)

	gw := walletserver.NewWalletServer(uint(*port), *gateway)
	gw.Run()

}
