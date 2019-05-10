package main

import (
	"fmt"
	"net"
	"os"

	"github.com/ravdin/programmingbitcoin/network"
)

func main() {
	addr := os.Args[1]
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(os.Stdout, "Receiving messages...\n")
	defer listener.Close()
	for {
		server, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go func() {
			serverNode := network.NewSimpleNode(network.WithConnection(server), true, true)
			serverNode.Receive(network.PingMessageOption())
		}()
	}
}
