package tx

import (
	"bytes"
	"encoding/hex"
	"github.com/ravdin/programmingbitcoin/ecc"
	"github.com/ravdin/programmingbitcoin/util"
	"math/big"
	"testing"
)

func TestTx(t *testing.T) {
	var fetcher = NewTxFetcher()
	fetcher.loadCache("tx.cache")
	serialized := `0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000006b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278afeffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600`
	raw := util.HexStringToBytes(serialized)
	reader := bytes.NewReader(raw)
	testTx := ParseTx(reader, false)

	if len(fetcher.cache) == 0 {
		t.Fatalf("Failed to load cache!")
	}

	t.Run("Test Version", func(t *testing.T) {
		if testTx.Version != 1 {
			t.Errorf("Expected 1, got %d", testTx.Version)
		}
	})

	t.Run("Test parse inputs", func(t *testing.T) {
		if len(testTx.TxIns) != 1 {
			t.Errorf("Expected 1 input!")
		}
		txIn := testTx.TxIns[0]
		expectedPrevTx := `d1c789a9c60383bf715f3f6ad9d14b91fe55f3deb369fe5d9280cb1a01793f81`
		expectedScriptSig := `6b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278a`
		var expectedSequence uint32 = 0xfffffffe
		if expectedPrevTx != hex.EncodeToString(txIn.PrevTx) {
			t.Errorf("Failed to parse PrevTx!")
		}
		if txIn.PrevIndex != 0 {
			t.Errorf("Failed to parse PrevIndex!")
		}
		if expectedScriptSig != hex.EncodeToString(txIn.ScriptSig.Serialize()) {
			t.Errorf("Failed to parse script sig!")
		}
		if expectedSequence != txIn.Sequence {
			t.Errorf("Failed to parse sequence!")
		}
	})

	t.Run("Test parse outputs", func(t *testing.T) {
		if len(testTx.TxOuts) != 2 {
			t.Errorf("Expected 2 outputs!")
		}
		expectedAmounts := []uint64{32454049, 10011545}
		expectedScriptSigs := []string{
			`1976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac`,
			`1976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac`,
		}
		for i, amount := range expectedAmounts {
			if testTx.TxOuts[i].Amount != amount {
				t.Errorf("Failed to parse amount!")
			}
		}
		for i, scriptSig := range expectedScriptSigs {
			if hex.EncodeToString(testTx.TxOuts[i].ScriptPubKey.Serialize()) != scriptSig {
				t.Errorf("Failed to parse script sig!")
			}
		}
	})

	t.Run("Test parse locktime", func(t *testing.T) {
		var expected uint32 = 410393
		if testTx.Locktime != expected {
			t.Errorf("Expected %v, got %v", expected, testTx.Locktime)
		}
	})

	t.Run("Test fee", func(t *testing.T) {
		var expected uint64 = 40000
		if testTx.Fee() != expected {
			t.Errorf("Expected %v, got %v", expected, testTx.Fee())
		}

		serialized2 := `010000000456919960ac691763688d3d3bcea9ad6ecaf875df5339e148a1fc61c6ed7a069e010000006a47304402204585bcdef85e6b1c6af5c2669d4830ff86e42dd205c0e089bc2a821657e951c002201024a10366077f87d6bce1f7100ad8cfa8a064b39d4e8fe4ea13a7b71aa8180f012102f0da57e85eec2934a82a585ea337ce2f4998b50ae699dd79f5880e253dafafb7feffffffeb8f51f4038dc17e6313cf831d4f02281c2a468bde0fafd37f1bf882729e7fd3000000006a47304402207899531a52d59a6de200179928ca900254a36b8dff8bb75f5f5d71b1cdc26125022008b422690b8461cb52c3cc30330b23d574351872b7c361e9aae3649071c1a7160121035d5c93d9ac96881f19ba1f686f15f009ded7c62efe85a872e6a19b43c15a2937feffffff567bf40595119d1bb8a3037c356efd56170b64cbcc160fb028fa10704b45d775000000006a47304402204c7c7818424c7f7911da6cddc59655a70af1cb5eaf17c69dadbfc74ffa0b662f02207599e08bc8023693ad4e9527dc42c34210f7a7d1d1ddfc8492b654a11e7620a0012102158b46fbdff65d0172b7989aec8850aa0dae49abfb84c81ae6e5b251a58ace5cfeffffffd63a5e6c16e620f86f375925b21cabaf736c779f88fd04dcad51d26690f7f345010000006a47304402200633ea0d3314bea0d95b3cd8dadb2ef79ea8331ffe1e61f762c0f6daea0fabde022029f23b3e9c30f080446150b23852028751635dcee2be669c2a1686a4b5edf304012103ffd6f4a67e94aba353a00882e563ff2722eb4cff0ad6006e86ee20dfe7520d55feffffff0251430f00000000001976a914ab0c0b2e98b1ab6dbf67d4750b0a56244948a87988ac005a6202000000001976a9143c82d7df364eb6c75be8c80df2b3eda8db57397088ac46430600`
		raw2 := util.HexStringToBytes(serialized2)
		reader = bytes.NewReader(raw2)
		testTx2 := ParseTx(reader, false)
		expected = 140500
		if expected != testTx2.Fee() {
			t.Errorf("Expected %v, got %v", expected, testTx2.Fee())
		}
	})

	t.Run("Test sig hash", func(t *testing.T) {
		tx := fetcher.fetch("452c629d67e41baec3ac6f04fe744b4b9617f8f859c63b3002f8684e7a4fee03", false, false)
		expected := util.HexStringToBytes("27e0c5994dec7824e56dec6b2fcb342eb7cdb0d0957c2fce9882f715e85d81a6")
		actual := tx.SigHash(0)
		if !bytes.Equal(actual, expected) {
			t.Errorf("Expected %x, got %x", expected, actual)
		}
	})

	t.Run("Test Verify p2pkh", func(t *testing.T) {
		txIds := []string{
			"452c629d67e41baec3ac6f04fe744b4b9617f8f859c63b3002f8684e7a4fee03",
			"5418099cc755cb9dd3ebc6cf1a7888ad53a1a3beb5a025bce89eb1bf7f1650a2",
		}
		for i, txId := range txIds {
			testnet := i == 1
			tx := fetcher.fetch(txId, testnet, false)
			if !tx.Verify() {
				t.Errorf("Verify failed!")
			}
		}
	})

	t.Run("Test Verify p2sh", func(t *testing.T) {
		tx := fetcher.fetch("46df1a9484d0a81d03ce0ee543ab6e1a23ed06175c104a178268fad381216c2b", false, false)
		if !tx.Verify() {
			t.Errorf("Verify failed!")
		}
	})

	t.Run("Test private key", func(t *testing.T) {
		pk := ecc.NewPrivateKey(big.NewInt(8675309))
		data := util.HexStringToBytes("010000000199a24308080ab26e6fb65c4eccfadf76749bb5bfa8cb08f291320b3c21e56f0d0d00000000ffffffff02408af701000000001976a914d52ad7ca9b3d096a38e752c2018e6fbc40cdf26f88ac80969800000000001976a914507b27411ccf7f16f10297de6cef3f291623eddf88ac00000000")
		reader := bytes.NewReader(data)
		tx_obj := ParseTx(reader, true)
		if !tx_obj.SignInput(0, pk) {
			t.Errorf("Private key sign failed!")
		}
		expected := `010000000199a24308080ab26e6fb65c4eccfadf76749bb5bfa8cb08f291320b3c21e56f0d0d0000006b4830450221008ed46aa2cf12d6d81065bfabe903670165b538f65ee9a3385e6327d80c66d3b502203124f804410527497329ec4715e18558082d489b218677bd029e7fa306a72236012103935581e52c354cd2f484fe8ed83af7a3097005b2f9c60bff71d35bd795f54b67ffffffff02408af701000000001976a914d52ad7ca9b3d096a38e752c2018e6fbc40cdf26f88ac80969800000000001976a914507b27411ccf7f16f10297de6cef3f291623eddf88ac00000000`
		actual := hex.EncodeToString(tx_obj.Serialize())
		if actual != expected {
			t.Errorf("Expected %v, got %v", expected, actual)
		}
	})
}
