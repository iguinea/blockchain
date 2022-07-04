package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/iguinea/cryptodemo/pkgs/utils"
)

func init() {
	fmt.Println("\033[2J") // clear
	log.SetPrefix(fmt.Sprintf("%.16s ", "Playground         "))
	log.SetOutput(os.Stdout)
}

func main() {
	localAddresses()
	//log.Println(utils.IsFoundHost("172.17.0.1", 5000))

	log.Printf("%v", utils.GetHost())
	//log.Println(utils.FindNeighbors("172.17.0.4", 5000, 1, 20, 5000, 5001))

}

func localAddresses() map[string]net.Addr {
	returnedValues := make(map[string]net.Addr)
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Print(fmt.Errorf("localAddresses: %v\n", err.Error()))
		return nil
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			log.Print(fmt.Errorf("localAddresses: %v\n", err.Error()))
			continue
		}
		for _, a := range addrs {
			returnedValues[i.Name] = a
			log.Printf("%v %v\n", i.Name, a)
		}
	}
	return returnedValues
}
