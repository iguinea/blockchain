package main

import (
	"bufio"
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
	//gateway := flag.String("gateway", "http://172.17.0.3:5000", "Blockchain gateway")
	flag.Parse()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the blockchain server gateway (http://IP:port): ")
	gatewayStr, _ := reader.ReadString('\n')
	gatewayStr = gatewayStr[:len(gatewayStr)-1]

	log.Printf("PortNumber to run WalletServer: %d", *port)
	log.Printf("Gateway to connect to : %s", gatewayStr)
	//log.Printf("Gateway to connect to : %s", *gateway)

	gw := walletserver.NewWalletServer(uint(*port), gatewayStr)
	//gw := walletserver.NewWalletServer(uint(*port), *gateway)

	fmt.Printf("%+v", gw)
	gw.Run()

}
