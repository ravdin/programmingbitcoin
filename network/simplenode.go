package network

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"reflect"

	"github.com/ravdin/programmingbitcoin/util"
)

type SimpleNode struct {
	Connection *net.TCPConn
	Testnet    bool
	Logging    bool
}

type NodeConnectOption func(*SimpleNode) *net.TCPConn
type ReceiveMessageTypeOption func() reflect.Type

func WithHostName(host string, ports ...int) NodeConnectOption {
	return func(node *SimpleNode) *net.TCPConn {
		var port int
		if len(ports) == 0 {
			port = 8333
			if node.Testnet {
				port = 18333
			}
		} else {
			port = ports[0]
		}
		result, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
		if err != nil {
			panic(err)
		}
		return result.(*net.TCPConn)
	}
}

func NewSimpleNode(option NodeConnectOption, testnet bool, logging bool) *SimpleNode {
	result := &SimpleNode{
		Testnet: testnet,
		Logging: logging,
	}
	result.Connection = option(result)
	return result
}

func (self *SimpleNode) Close() error {
	return self.Connection.Close()
}

// Do a handshake with the other node.
// Handshake is sending a version message and getting a verack back.
func (self *SimpleNode) Handshake() (bool, error) {
	if ok, err := self.Send(NewVersionMessage(nil)); !ok {
		return ok, err
	}
	verack, err := self.WaitFor(VerackMessageOption())
	if err != nil {
		return false, err
	}
	if verack == nil {
		return false, errors.New("No response received!")
	}
	return true, nil
}

// Send a message to the connected node.
func (self *SimpleNode) Send(message Message) (bool, error) {
	envelope := NewEnvelope(message.Command(), message.Serialize(), self.Testnet)
	if self.Logging {
		fmt.Fprintf(os.Stdout, "sending: %v\n", envelope)
	}
	_, err := self.Connection.Write(envelope.Serialize())
	if err != nil {
		return false, err
	}
	return true, nil
}

// Read a message from the socket.
func (self *SimpleNode) Read() (*Envelope, error) {
	bufCh := make(chan []byte)
	errCh := make(chan error)
	go func(conn *net.TCPConn, bufCh chan []byte, errCh chan error) {
		header := make([]byte, 24)
		_, err := io.ReadFull(conn, header)
		if err != nil {
			errCh <- err
			return
		}
		bufCh <- header
		payloadLength := int(util.LittleEndianToInt32(header[16:20]))
		payload := make([]byte, payloadLength)
		_, err = io.ReadFull(conn, payload)
		if err != nil {
			errCh <- err
			return
		}
		bufCh <- payload
		close(bufCh)
	}(self.Connection, bufCh, errCh)
	response := make([]byte, 0)
	for buf := range bufCh {
		select {
		case err := <-errCh:
			return nil, err
		default:
		}
		response = append(response, buf...)
	}
	return ParseEnvelope(bytes.NewReader(response), self.Testnet), nil
}

// Wait for one of the messages in the list
func (self *SimpleNode) WaitFor(messageTypes ...ReceiveMessageTypeOption) (Message, error) {
	commands := make(map[string]Message)
	for _, option := range messageTypes {
		messageType := option()
		message, ok := reflect.New(messageType.Elem()).Interface().(Message)
		if !ok {
			panic("Failed to cast to Message type!")
		}
		commands[string(message.Command())] = message
	}
	for {
		envelope, err := self.Read()
		if err != nil {
			return nil, err
		}
		command := string(envelope.Command)
		if self.Logging {
			fmt.Fprintf(os.Stdout, "received: %s\n", command)
		}
		switch command {
		case "version":
			self.Send(NewVerackMessage())
		case "ping":
			self.Send(NewPongMessage(envelope.Payload))
		}
		if result, ok := commands[command]; ok {
			result.Parse(envelope.Stream())
			return result, nil
		}
	}
}
