package merkle

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/ravdin/programmingbitcoin/util"
)

func TestMerkleBlock(t *testing.T) {
	strMerkleBlock := `00000020df3b053dc46f162a9b00c7f0d5124e2676d47bbe7c5d0793a500000000000000ef445fef2ed495c275892206ca533e7411907971013ab83e3b47bd0d692d14d4dc7c835b67d8001ac157e670bf0d00000aba412a0d1480e370173072c9562becffe87aa661c1e4a6dbc305d38ec5dc088a7cf92e6458aca7b32edae818f9c2c98c37e06bf72ae0ce80649a38655ee1e27d34d9421d940b16732f24b94023e9d572a7f9ab8023434a4feb532d2adfc8c2c2158785d1bd04eb99df2e86c54bc13e139862897217400def5d72c280222c4cbaee7261831e1550dbb8fa82853e9fe506fc5fda3f7b919d8fe74b6282f92763cef8e625f977af7c8619c32a369b832bc2d051ecd9c73c51e76370ceabd4f25097c256597fa898d404ed53425de608ac6bfe426f6e2bb457f1c554866eb69dcb8d6bf6f880e9a59b3cd053e6c7060eeacaacf4dac6697dac20e4bd3f38a2ea2543d1ab7953e3430790a9f81e1c67f5b58c825acf46bd02848384eebe9af917274cdfbb1a28a5d58a23a17977def0de10d644258d9c54f886d47d293a411cb6226103b55635`
	reader := bytes.NewReader(util.HexStringToBytes(strMerkleBlock))
	block := new(MerkleBlock)
	block.Parse(reader)

	t.Run("test parse", func(t *testing.T) {
		var version uint32 = 0x20000000
		if version != block.Version {
			t.Errorf("Expected %d, got %d", version, block.Version)
		}
		merkleRoot := make([]byte, 32)
		copy(merkleRoot, block.MerkleRoot[:])
		expected := fmt.Sprintf("%x", util.ReverseByteArray(merkleRoot))
		actual := `ef445fef2ed495c275892206ca533e7411907971013ab83e3b47bd0d692d14d4`
		if actual != expected {
			t.Errorf("Expected %s, got %s", expected, actual)
		}
		prevBlock := make([]byte, 32)
		copy(prevBlock, block.PrevBlock[:])
		expected = fmt.Sprintf("%x", util.ReverseByteArray(prevBlock))
		actual = `df3b053dc46f162a9b00c7f0d5124e2676d47bbe7c5d0793a500000000000000`
		if actual != expected {
			t.Errorf("Expected %s, got %s", expected, actual)
		}
		timestamp := util.LittleEndianToInt32(util.HexStringToBytes(`dc7c835b`))
		if timestamp != block.Timestamp {
			t.Errorf("Expected %d, got %d", timestamp, block.Timestamp)
		}
		bits := util.HexStringToBytes(`67d8001a`)
		if !bytes.Equal(bits, block.Bits[:]) {
			t.Errorf("Expected %x, got %x", bits, block.Bits)
		}
		nonce := util.HexStringToBytes(`c157e670`)
		if !bytes.Equal(nonce, block.Nonce[:]) {
			t.Errorf("Expected %x, got %x", nonce, block.Nonce)
		}
		total := util.LittleEndianToInt32(util.HexStringToBytes(`bf0d0000`))
		if total != block.Total {
			t.Errorf("Expected %d, got %d", total, block.Total)
		}
		hashes := [][]byte{
			util.HexStringToBytes("ba412a0d1480e370173072c9562becffe87aa661c1e4a6dbc305d38ec5dc088a"),
			util.HexStringToBytes("7cf92e6458aca7b32edae818f9c2c98c37e06bf72ae0ce80649a38655ee1e27d"),
			util.HexStringToBytes("34d9421d940b16732f24b94023e9d572a7f9ab8023434a4feb532d2adfc8c2c2"),
			util.HexStringToBytes("158785d1bd04eb99df2e86c54bc13e139862897217400def5d72c280222c4cba"),
			util.HexStringToBytes("ee7261831e1550dbb8fa82853e9fe506fc5fda3f7b919d8fe74b6282f92763ce"),
			util.HexStringToBytes("f8e625f977af7c8619c32a369b832bc2d051ecd9c73c51e76370ceabd4f25097"),
			util.HexStringToBytes("c256597fa898d404ed53425de608ac6bfe426f6e2bb457f1c554866eb69dcb8d"),
			util.HexStringToBytes("6bf6f880e9a59b3cd053e6c7060eeacaacf4dac6697dac20e4bd3f38a2ea2543"),
			util.HexStringToBytes("d1ab7953e3430790a9f81e1c67f5b58c825acf46bd02848384eebe9af917274c"),
			util.HexStringToBytes("dfbb1a28a5d58a23a17977def0de10d644258d9c54f886d47d293a411cb62261"),
		}
		for i, hash := range hashes {
			util.ReverseByteArray(hash)
			if !bytes.Equal(hash, block.Hashes[i]) {
				t.Errorf("Expected %x, got %x", hash, block.Hashes[i])
			}
		}
		flags := util.HexStringToBytes(`b55635`)
		if !bytes.Equal(flags, block.Flags) {
			t.Errorf("Expected %x, got %x", flags, block.Flags)
		}
	})

	t.Run("test IsValid", func(t *testing.T) {
		if !block.IsValid() {
			t.Errorf("IsValid failed!")
		}
	})
}
