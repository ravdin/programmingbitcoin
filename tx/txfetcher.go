package tx

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/ravdin/programmingbitcoin/util"
)

type txFetcher struct {
	cache map[string]*Tx
}

var (
	// Singleton instance.
	instance *txFetcher
	once     sync.Once
)

func NewTxFetcher() *txFetcher {
	once.Do(func() {
		instance = &txFetcher{cache: make(map[string]*Tx)}
	})

	return instance
}

func getUrl(testnet bool) string {
	if testnet {
		return "http://testnet.programmingbitcoin.com"
	}
	return "http://mainnet.programmingbitcoin.com"
}

func (self *txFetcher) fetch(txId string, testnet bool, fresh bool) *Tx {
	if tx, ok := self.cache[txId]; ok && !fresh {
		tx.Testnet = testnet
		return tx
	}
	url := fmt.Sprintf("%s/tx/%s.hex", getUrl(testnet), txId)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	raw := make([]byte, hex.DecodedLen(len(body)))
	hex.Decode(raw, body)
	var tx *Tx
	if raw[4] == 0 {
		length := len(raw)
		locktime := util.LittleEndianToInt32(raw[length-4:])
		copy(raw[4:], raw[6:])
		reader := bytes.NewReader(raw[:length-2])
		tx = ParseTx(reader, testnet)
		tx.Locktime = locktime
	} else {
		reader := bytes.NewReader(raw)
		tx = ParseTx(reader, testnet)
	}
	if txId != tx.Id() {
		panic(fmt.Sprintf("Not the same id: %s vs %s", txId, tx.Id()))
	}
	self.cache[txId] = tx
	tx.Testnet = testnet
	return tx
}

func (self *txFetcher) loadCache(filename string) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	var v map[string]string
	err = json.Unmarshal(data, &v)
	if err != nil {
		panic(err)
	}
	for k, rawHex := range v {
		raw := util.HexStringToBytes(rawHex)
		var t *Tx
		if raw[4] == 0 {
			tmp := make([]byte, len(raw)-2)
			copy(tmp[:4], raw[:4])
			copy(tmp[4:], raw[6:])
			reader := bytes.NewReader(tmp)
			t = ParseTx(reader, false)
			t.Locktime = util.LittleEndianToInt32(tmp[len(tmp)-4:])
		} else {
			reader := bytes.NewReader(raw)
			t = ParseTx(reader, false)
		}
		self.cache[k] = t
	}
}
