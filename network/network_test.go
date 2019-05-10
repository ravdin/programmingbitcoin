package network

import (
	"fmt"
	"net"
	"os"
	"testing"
)

func TestNetwork(t *testing.T) {
	t.Run("test ping", func(t *testing.T) {
		client, server := net.Pipe()
		nonce := [8]byte{0, 0, 0, 0, 0, 0, 0, 1}
		serverCommand := make(chan string)
		//responseMessage := make(chan Message)
		clientNode := NewSimpleNode(WithConnection(client), true, true)
		go func() {
			serverNode := NewSimpleNode(WithConnection(server), true, true)
			fmt.Fprintf(os.Stdout, "Accepting message...\n")
			message := serverNode.Receive(PingMessageOption())
			serverCommand <- string(message.Command())
			fmt.Fprintf(os.Stdout, "Closing server connection...\n")
			server.Close()
		}()
		ping := NewPingMessage(nonce)
		clientNode.Send(ping)
		fmt.Fprintf(os.Stdout, "Send called...\n")
		fmt.Fprintf(os.Stdout, "Closing client connection...\n")
		clientNode.Socket.Close()
	})
}
