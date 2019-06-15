package network

import (
	"bytes"
	"math/rand"
	"time"

	"github.com/ravdin/programmingbitcoin/util"
)

// VersionMessage represents a "version" message.
type VersionMessage struct {
	Version          uint32
	Services         uint64
	Timestamp        uint64
	ReceiverServices uint64
	ReceiverIP       [4]byte
	ReceiverPort     uint16
	SenderServices   uint64
	SenderIP         [4]byte
	SenderPort       uint16
	Nonce            [8]byte
	UserAgent        string
	LatestBlock      uint32
	Relay            bool
}

// Constants that represent the fields in a VersionMessage.
const (
	VersionArg = iota
	ServicesArg
	TimestampArg
	ReceiverServicesArg
	ReceiverIPArg
	ReceiverPortArg
	SenderServicesArg
	SenderIPArg
	SenderPortArg
	NonceArg
	UserAgentArg
	LatestBlockArg
	RelayArg
)

var defaultValues = map[int]interface{}{
	VersionArg:          uint32(70015),
	ServicesArg:         uint64(0),
	TimestampArg:        nil,
	ReceiverServicesArg: uint64(0),
	ReceiverIPArg:       [4]byte{0, 0, 0, 0},
	ReceiverPortArg:     uint16(8333),
	SenderServicesArg:   uint64(0),
	SenderIPArg:         [4]byte{0, 0, 0, 0},
	SenderPortArg:       uint16(8333),
	NonceArg:            nil,
	UserAgentArg:        "/programmingbitcoin:0.1/",
	LatestBlockArg:      uint32(0),
	RelayArg:            false,
}

// NewVersionMessage initializes a VersionMessage object.
// args: map of initialization values. Use this to override the default values.
func NewVersionMessage(args map[int]interface{}) *VersionMessage {
	values := make(map[int]interface{})
	for k, v := range defaultValues {
		values[k] = v
	}
	if args != nil {
		for k, v := range args {
			values[k] = v
		}
	}
	if values[TimestampArg] == nil {
		values[TimestampArg] = uint64(time.Now().Unix())
	}
	if values[NonceArg] == nil {
		var nonce [8]byte
		copy(nonce[:], util.Int64ToLittleEndian(rand.Uint64()))
		values[NonceArg] = nonce
	}
	return &VersionMessage{
		Version:          values[VersionArg].(uint32),
		Services:         values[ServicesArg].(uint64),
		Timestamp:        values[TimestampArg].(uint64),
		ReceiverServices: values[ReceiverServicesArg].(uint64),
		ReceiverIP:       values[ReceiverIPArg].([4]byte),
		ReceiverPort:     values[ReceiverPortArg].(uint16),
		SenderServices:   values[SenderServicesArg].(uint64),
		SenderIP:         values[SenderIPArg].([4]byte),
		SenderPort:       values[SenderPortArg].(uint16),
		Nonce:            values[NonceArg].([8]byte),
		UserAgent:        values[UserAgentArg].(string),
		LatestBlock:      values[LatestBlockArg].(uint32),
		Relay:            values[RelayArg].(bool),
	}
}

// Command sequence that identifies this type of message.
func (*VersionMessage) Command() []byte {
	return []byte("version")
}

// Serialize this message to send over the network
func (msg *VersionMessage) Serialize() []byte {
	version := util.Int32ToLittleEndian(msg.Version)
	services := util.Int64ToLittleEndian(msg.Services)
	timestamp := util.Int64ToLittleEndian(msg.Timestamp)
	receiverServices := util.Int64ToLittleEndian(msg.ReceiverServices)
	receiverIP := make([]byte, 16)
	copy(receiverIP[10:12], []byte{0xff, 0xff})
	copy(receiverIP[12:], msg.ReceiverIP[:])
	receiverPort := util.Int16ToLittleEndian(msg.ReceiverPort)
	senderServices := util.Int64ToLittleEndian(msg.SenderServices)
	senderIP := make([]byte, 16)
	copy(senderIP[10:12], []byte{0xff, 0xff})
	copy(senderIP[12:], msg.SenderIP[:])
	senderPort := util.Int16ToLittleEndian(msg.SenderPort)
	userAgentLength := util.EncodeVarInt(len(msg.UserAgent))
	latestBlock := util.Int32ToLittleEndian(msg.LatestBlock)
	var relay byte
	if msg.Relay {
		relay = 1
	}
	serializedLength := 85 + len(userAgentLength) + len(msg.UserAgent)
	result := make([]byte, serializedLength)
	copy(result[:4], version)
	copy(result[4:12], services)
	copy(result[12:20], timestamp)
	copy(result[20:28], receiverServices)
	copy(result[28:44], receiverIP)
	copy(result[44:46], receiverPort)
	copy(result[46:54], senderServices)
	copy(result[54:70], senderIP)
	copy(result[70:72], senderPort)
	copy(result[72:80], msg.Nonce[:])
	copy(result[80:80+len(userAgentLength)], userAgentLength)
	pos := 80 + len(userAgentLength)
	copy(result[pos:pos+len(msg.UserAgent)], []byte(msg.UserAgent))
	pos += len(msg.UserAgent)
	copy(result[pos:pos+4], latestBlock)
	pos += 4
	result[pos] = relay
	return result
}

// Parse a message from a byte steam.
func (msg *VersionMessage) Parse(reader *bytes.Reader) Message {
	version := make([]byte, 4)
	reader.Read(version)
	msg.Version = util.LittleEndianToInt32(version)
	services := make([]byte, 8)
	reader.Read(services)
	msg.Services = util.LittleEndianToInt64(services)
	timestamp := make([]byte, 8)
	reader.Read(timestamp)
	msg.Timestamp = util.LittleEndianToInt64(timestamp)
	receiverServices := make([]byte, 8)
	reader.Read(receiverServices)
	msg.ReceiverServices = util.LittleEndianToInt64(receiverServices)
	receiverIP := make([]byte, 16)
	reader.Read(receiverIP)
	copy(msg.ReceiverIP[:], receiverIP[12:])
	receiverPort := make([]byte, 2)
	reader.Read(receiverPort)
	msg.ReceiverPort = util.LittleEndianToInt16(receiverPort)
	senderServices := make([]byte, 8)
	reader.Read(senderServices)
	msg.ReceiverServices = util.LittleEndianToInt64(senderServices)
	senderIP := make([]byte, 16)
	reader.Read(senderIP)
	copy(msg.SenderIP[:], senderIP[12:])
	senderPort := make([]byte, 2)
	reader.Read(senderPort)
	msg.SenderPort = util.LittleEndianToInt16(senderPort)
	nonce := make([]byte, 8)
	reader.Read(nonce)
	copy(msg.Nonce[:], nonce)
	userAgentLength := util.ReadVarInt(reader)
	userAgent := make([]byte, userAgentLength)
	reader.Read(userAgent)
	msg.UserAgent = string(userAgent)
	latestBlock := make([]byte, 4)
	reader.Read(latestBlock)
	msg.LatestBlock = util.LittleEndianToInt32(latestBlock)
	relay, _ := reader.ReadByte()
	msg.Relay = relay != 0
	return msg
}
