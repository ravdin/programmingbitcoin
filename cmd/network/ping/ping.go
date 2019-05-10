package main

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
	"strconv"

	"github.com/ravdin/programmingbitcoin/network"
	"github.com/ravdin/programmingbitcoin/util"
)

func main() {
	host := os.Args[1]
	port, _ := strconv.ParseInt(os.Args[2], 10, 32)
	var nonce [8]byte
	copy(nonce[:], util.Int64ToLittleEndian(rand.Uint64()))
	clientNode := network.NewSimpleNode(network.WithHostName(host, int(port)), true, true)
	message := network.NewPingMessage(nonce)
	receivedCh, errCh := clientNode.SendAsync(message)
	select {
	case err := <-errCh:
		fmt.Fprintf(os.Stderr, "%v\n", err)
	case pong := <-receivedCh:
		fmt.Fprintf(os.Stdout, "Received response: %s\n", hex.EncodeToString(pong.Serialize()))
	}
}
