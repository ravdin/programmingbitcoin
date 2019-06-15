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

// SimpleNode is a utility class for creating a TCP connection to a bitcoin network.
type SimpleNode struct {
	Connection *net.TCPConn
	Testnet    bool
	Logging    bool
}

// NodeConnectOption is an alias for a function that returns a TCP connection.
type NodeConnectOption func(*SimpleNode) *net.TCPConn

// ReceiveMessageTypeOption is an alias for a function that returns a Message type.
type ReceiveMessageTypeOption func() reflect.Type

// WithHostName returns a function for initializing a TCP connection from a host name.
// If no port is passed, use a default. If multiple port numbers are passed, use the first one.
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

// NewSimpleNode creates a SimpleNode and initializes a TCP connection.
// option: function for returning a TCPConn from a SimpleNode.
// testnet: Indicate if the node should connect to a testnet.
// logging: Set to true for more verbose messages to standard out.
func NewSimpleNode(option NodeConnectOption, testnet bool, logging bool) *SimpleNode {
	result := &SimpleNode{
		Testnet: testnet,
		Logging: logging,
	}
	result.Connection = option(result)
	return result
}

// Close the connection.
func (node *SimpleNode) Close() error {
	return node.Connection.Close()
}

// Handshake is sending a version message and getting a verack back.
func (node *SimpleNode) Handshake() (bool, error) {
	if ok, err := node.Send(NewVersionMessage(nil)); !ok {
		return ok, err
	}
	verack, err := node.WaitFor(VerackMessageOption())
	if err != nil {
		return false, err
	}
	if verack == nil {
		return false, errors.New("no response received")
	}
	return true, nil
}

// Send a message to the connected node.
func (node *SimpleNode) Send(message Message) (bool, error) {
	envelope := NewEnvelope(message.Command(), message.Serialize(), node.Testnet)
	if node.Logging {
		fmt.Fprintf(os.Stdout, "sending: %v\n", envelope)
	}
	_, err := node.Connection.Write(envelope.Serialize())
	if err != nil {
		return false, err
	}
	return true, nil
}

// Read a message from the socket.
func (node *SimpleNode) Read() (*Envelope, error) {
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
	}(node.Connection, bufCh, errCh)
	response := make([]byte, 0)
	for buf := range bufCh {
		select {
		case err := <-errCh:
			return nil, err
		default:
		}
		response = append(response, buf...)
	}
	return ParseEnvelope(bytes.NewReader(response), node.Testnet), nil
}

// WaitFor waits for one of the messages in the list
// Return a Message if successful and an error otherwise.
func (node *SimpleNode) WaitFor(messageTypes ...ReceiveMessageTypeOption) (Message, error) {
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
		envelope, err := node.Read()
		if err != nil {
			return nil, err
		}
		command := string(envelope.Command)
		if node.Logging {
			fmt.Fprintf(os.Stdout, "received: %s\n", command)
		}
		switch command {
		case "version":
			node.Send(NewVerackMessage())
		case "ping":
			node.Send(NewPongMessage(envelope.Payload))
		}
		if result, ok := commands[command]; ok {
			result.Parse(envelope.Stream())
			return result, nil
		}
	}
}
