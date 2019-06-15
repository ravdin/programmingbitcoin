package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/ravdin/programmingbitcoin/block"
	"github.com/ravdin/programmingbitcoin/network"
	"github.com/ravdin/programmingbitcoin/util"
)

func main() {
	previous := block.Parse(bytes.NewReader(block.GenesisBlock))
	firstEpochTimestamp := previous.Timestamp
	expectedBits := block.LowestBits
	count := 1
	node := network.NewSimpleNode(network.WithHostName("mainnet.programmingbitcoin.com"), false, false)
	defer node.Close()
	if ok, err := node.Handshake(); !ok {
		panic(err)
	}
	for i := 0; i < 19; i++ {
		getheaders := network.NewGetHeadersMessage(previous.Hash())
		if ok, err := node.Send(getheaders); !ok {
			panic(err)
		}
		msg, err := node.WaitFor(network.HeadersMessageOption())
		if err != nil {
			panic(err)
		}
		headers := msg.(*network.HeadersMessage)
		for _, header := range headers.Blocks {
			if !header.CheckPow() {
				panic(fmt.Errorf("Bad PoW at block %d\n", count))
			}
			if !bytes.Equal(header.PrevBlock[:], previous.Hash()) {
				panic(fmt.Errorf("Discontinuous block at %d\n", count))
			}
			if count%2016 == 0 {
				timeDiff := previous.Timestamp - firstEpochTimestamp
				expectedBits = util.CalculateNewBits(previous.Bits[:], int(timeDiff))
				fmt.Fprintf(os.Stdout, "%x\n", expectedBits)
				firstEpochTimestamp = header.Timestamp
			}
			if !bytes.Equal(header.Bits[:], expectedBits) {
				fmt.Fprintf(os.Stdout, "expected: %x, actual: %x\n", expectedBits, header.Bits[:])
				panic(fmt.Errorf("Bad bits at block %d\n", count))
			}
			previous = header
			count++
		}
	}
}
