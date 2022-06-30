package utils

import (
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"time"
)

func IsFoundHost(host string, port uint) bool {
	target := fmt.Sprintf("%s:%d", host, port)

	_, err := net.DialTimeout("tcp", target, 1*time.Second)
	if err != nil {
		log.Printf("%s %v\n", target, err)
		return false
	}

	return true

}

var PATTERN = regexp.MustCompile(`((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?\.){3})(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`)

func FindNeighbors(myHost string, myPort uint16, startIp uint8, endIp uint8, startPort uint16, endPort uint16) []string {
	address := fmt.Sprintf("%s:%d", myHost, myPort)
	//fmt.Printf("Address: %s\n", address)
	m := PATTERN.FindStringSubmatch(myHost)
	if m == nil {
		return nil
	}
	prefixHost := m[1]
	//	lastIp, _ := strconv.Atoi(m[len(m)-1])
	neighbors := make([]string, 0)

	for port := startPort; port <= endPort; port += 1 {
		for ip := startIp; ip <= endIp; ip += 1 {
			//guessHost := fmt.Sprintf("%s%d", prefixHost, lastIp+int(ip))
			guessHost := fmt.Sprintf("%s%d", prefixHost, int(ip))
			guessTarget := fmt.Sprintf("%s:%d", guessHost, port)
			//fmt.Printf("guessHost: %s   --- guessTarget: %s\n", guessHost, guessTarget)
			if guessTarget != address && IsFoundHost(guessHost, uint(port)) {
				fmt.Printf("found host: %s\n", guessTarget)
				neighbors = append(neighbors, guessTarget)
			}
		}
	}
	return neighbors
}

func GetHost() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "127.0.0.1"
	}
	address, err := net.Lo
}
