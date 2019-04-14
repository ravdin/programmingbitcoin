package tx

import (
	"bytes"
	"encoding/json"
	"github.com/ravdin/programmingbitcoin/util"
	"io/ioutil"
)

type TxFetcher struct {
	cache map[string]*Tx
}

func (self *TxFetcher) loadCache(filename string) {
	data, _ := ioutil.ReadFile(filename)
	var v map[string]string
	json.Unmarshal(data, v)
	for k, rawHex := range v {
		raw := util.HexStringToBytes(rawHex)
		var t *Tx
		if raw[4] == 0 {
			tmp := make([]byte, len(raw)-1)
			copy(tmp[:4], raw[:4])
			copy(tmp[5:], raw[6:])
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
