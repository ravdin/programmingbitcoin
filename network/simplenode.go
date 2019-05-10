package network

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"reflect"
)

type SimpleNode struct {
	Socket  net.Conn
	Testnet bool
	Logging bool
}

type NodeConnectOption func(*SimpleNode) net.Conn
type ReceiveMessageTypeOption func() reflect.Type

func WithHostName(host string, ports ...int) NodeConnectOption {
	return func(node *SimpleNode) net.Conn {
		var port int
		if len(ports) == 0 {
			port = 8333
			if node.Testnet {
				port = 18333
			}
		} else {
			port = ports[0]
		}
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
		if err != nil {
			panic(err)
		}
		return conn
	}
}

func WithConnection(conn net.Conn) NodeConnectOption {
	return func(node *SimpleNode) net.Conn {
		return conn
	}
}

func NewSimpleNode(option NodeConnectOption, testnet bool, logging bool) *SimpleNode {
	result := &SimpleNode{Testnet: testnet, Logging: logging}
	result.Socket = option(result)
	return result
}

// Send a message to the connected node.
func (self *SimpleNode) Send(message Message) {
	envelope := NewEnvelope(message.Command(), message.Serialize(), self.Testnet)
	if self.Logging {
		fmt.Fprintf(os.Stdout, "sending: %s\n", hex.EncodeToString(envelope.Serialize()))
	}
	self.Socket.Write(envelope.Serialize())
}

// Send a message and return a response in a channel.
func (self *SimpleNode) SendAsync(message Message) (chan Message, chan error) {
	receivedCh := make(chan Message)
	errCh := make(chan error)
	self.Send(message)
	if tcpconn, ok := self.Socket.(*net.TCPConn); ok {
		tcpconn.CloseWrite()
	} else {
		panic("SendAsync requires a TCP connection!")
	}
	go func() {
		envelope, err := self.Read()
		if err != nil {
			fmt.Fprintf(os.Stdout, "Error reading message: %v\n", err)
			errCh <- err
			return
		}
		received := message.AckMessage()
		received.Parse(envelope.Stream())
		receivedCh <- received
	}()
	return receivedCh, errCh
}

func (self *SimpleNode) Read() (*Envelope, error) {
	data, err := ioutil.ReadAll(self.Socket)
	if err != nil {
		fmt.Fprintf(os.Stdout, "Error reading: %v\n", err)
		return nil, err
	}
	if self.Logging {
		fmt.Fprintf(os.Stdout, "received: %s\n", hex.EncodeToString(data))
	}
	reader := bytes.NewReader(data)
	envelope := ParseEnvelope(reader, self.Testnet)
	return envelope, nil
}

func (self *SimpleNode) Receive(messageTypes ...ReceiveMessageTypeOption) Message {
	envelope, err := self.Read()
	if err != nil {
		panic(err)
	}
	reader := envelope.Stream()
	for _, option := range messageTypes {
		messageType := option()
		message, ok := reflect.New(messageType.Elem()).Interface().(Message)
		if !ok {
			panic("Failed to cast to Message type!")
		}
		if bytes.Equal(message.Command(), envelope.Command) {
			message.Parse(reader)
			response := message.AckMessage()
			if response != nil {
				if self.Logging {
					fmt.Fprintf(os.Stdout, "Received %s, sending %s...\n", message.Command(), response.Command())
				}
				self.Send(response)
				if tcpconn, ok := self.Socket.(*net.TCPConn); ok {
					tcpconn.CloseWrite()
				}
			}
			return message.Parse(reader)
		}
	}
	return nil
}
