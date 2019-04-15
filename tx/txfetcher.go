package tx

import (
	"bytes"
	"encoding/json"
	"github.com/ravdin/programmingbitcoin/util"
	"io/ioutil"
	"sync"
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

func (self *txFetcher) fetch(txId string, testnet bool, fresh bool) *Tx {
	if tx, ok := self.cache[txId]; ok && !fresh {
		tx.Testnet = testnet
		return tx
	}
	panic("Not implemented!")
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
