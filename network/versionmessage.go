package network

import (
	"math/rand"
	"time"

	"github.com/ravdin/programmingbitcoin/util"
)

type VersionMessage struct {
	Version          uint32
	Services         uint64
	Timestamp        uint64
	ReceiverServices uint64
	ReceiverIp       [4]byte
	ReceiverPort     uint16
	SenderServices   uint64
	SenderIp         [4]byte
	SenderPort       uint16
	Nonce            [8]byte
	UserAgent        string
	LatestBlock      uint32
	Relay            bool
}

const (
	VERSION = iota
	SERVICES
	TIMESTAMP
	RECEIVER_SERVICES
	RECEIVER_IP
	RECEIVER_PORT
	SENDER_SERVICES
	SENDER_IP
	SENDER_PORT
	NONCE
	USER_AGENT
	LATEST_BLOCK
	RELAY
)

var defaultValues = map[int]interface{}{
	VERSION:           uint32(70015),
	SERVICES:          uint64(0),
	TIMESTAMP:         nil,
	RECEIVER_SERVICES: uint64(0),
	RECEIVER_IP:       [4]byte{0, 0, 0, 0},
	RECEIVER_PORT:     uint16(8333),
	SENDER_SERVICES:   uint64(0),
	SENDER_IP:         [4]byte{0, 0, 0, 0},
	SENDER_PORT:       uint16(8333),
	NONCE:             nil,
	USER_AGENT:        "/programmingbitcoin:0.1/",
	LATEST_BLOCK:      uint32(0),
	RELAY:             false,
}

func NewVersionMessage(args map[int]interface{}) *VersionMessage {
	values := make(map[int]interface{})
	for k, v := range defaultValues {
		values[k] = v
	}
	for k, v := range args {
		values[k] = v
	}
	if values[TIMESTAMP] == nil {
		values[TIMESTAMP] = time.Now().Unix()
	}
	if values[NONCE] == nil {
		var nonce [8]byte
		copy(nonce[:], util.Int64ToLittleEndian(rand.Uint64()))
		values[NONCE] = nonce
	}
	return &VersionMessage{
		Version:          values[VERSION].(uint32),
		Services:         values[SERVICES].(uint64),
		Timestamp:        values[TIMESTAMP].(uint64),
		ReceiverServices: values[RECEIVER_SERVICES].(uint64),
		ReceiverIp:       values[RECEIVER_IP].([4]byte),
		ReceiverPort:     values[RECEIVER_PORT].(uint16),
		SenderServices:   values[SENDER_SERVICES].(uint64),
		SenderIp:         values[SENDER_IP].([4]byte),
		SenderPort:       values[SENDER_PORT].(uint16),
		Nonce:            values[NONCE].([8]byte),
		UserAgent:        values[USER_AGENT].(string),
		LatestBlock:      values[LATEST_BLOCK].(uint32),
		Relay:            values[RELAY].(bool),
	}
}

func (*VersionMessage) Command() []byte {
	return []byte("version")
}

// Serialize this message to send over the network
func (self *VersionMessage) Serialize() []byte {
	version := util.Int32ToLittleEndian(self.Version)
	services := util.Int64ToLittleEndian(self.Services)
	timestamp := util.Int64ToLittleEndian(self.Timestamp)
	receiverServices := util.Int64ToLittleEndian(self.ReceiverServices)
	receiverIp := make([]byte, 16)
	copy(receiverIp[10:12], []byte{0xff, 0xff})
	copy(receiverIp[12:], self.ReceiverIp[:])
	receiverPort := util.Int16ToLittleEndian(self.ReceiverPort)
	senderServices := util.Int64ToLittleEndian(self.SenderServices)
	senderIp := make([]byte, 16)
	copy(senderIp[10:12], []byte{0xff, 0xff})
	copy(senderIp[12:], self.SenderIp[:])
	senderPort := util.Int16ToLittleEndian(self.SenderPort)
	userAgentLength := util.EncodeVarInt(len(self.UserAgent))
	latestBlock := util.Int32ToLittleEndian(self.LatestBlock)
	var relay byte = 0
	if self.Relay {
		relay = 1
	}
	serializedLength := 85 + len(userAgentLength) + len(self.UserAgent)
	result := make([]byte, serializedLength)
	copy(result[:4], version)
	copy(result[4:12], services)
	copy(result[12:20], timestamp)
	copy(result[20:28], receiverServices)
	copy(result[28:44], receiverIp)
	copy(result[44:46], receiverPort)
	copy(result[46:54], senderServices)
	copy(result[54:70], senderIp)
	copy(result[70:72], senderPort)
	copy(result[72:80], self.Nonce[:])
	copy(result[80:80+len(userAgentLength)], userAgentLength)
	pos := 80 + len(userAgentLength)
	copy(result[pos:pos+len(self.UserAgent)], []byte(self.UserAgent))
	pos += len(self.UserAgent)
	copy(result[pos:pos+4], latestBlock)
	pos += 4
	result[pos] = relay
	return result
}
