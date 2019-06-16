package merkle

import (
	"fmt"
	"testing"

	"github.com/ravdin/programmingbitcoin/util"
)

func TestMerkleTree(t *testing.T) {
	t.Run("test init", func(t *testing.T) {
		tree := NewTree(9)
		values := []int{1, 2, 3, 5, 9}
		for i, value := range values {
			actual := len(tree.Nodes[i])
			expected := value
			if actual != expected {
				t.Errorf("Expected %d, got %d", expected, actual)
			}
		}
	})

	t.Run("test populate tree 1", func(t *testing.T) {
		hashes := [][]byte{
			util.HexStringToBytes("9745f7173ef14ee4155722d1cbf13304339fd00d900b759c6f9d58579b5765fb"),
			util.HexStringToBytes("5573c8ede34936c29cdfdfe743f7f5fdfbd4f54ba0705259e62f39917065cb9b"),
			util.HexStringToBytes("82a02ecbb6623b4274dfcab82b336dc017a27136e08521091e443e62582e8f05"),
			util.HexStringToBytes("507ccae5ed9b340363a0e6d765af148be9cb1c8766ccc922f83e4ae681658308"),
			util.HexStringToBytes("a7a4aec28e7162e1e9ef33dfa30f0bc0526e6cf4b11a576f6c5de58593898330"),
			util.HexStringToBytes("bb6267664bd833fd9fc82582853ab144fece26b7a8a5bf328f8a059445b59add"),
			util.HexStringToBytes("ea6d7ac1ee77fbacee58fc717b990c4fcccf1b19af43103c090f601677fd8836"),
			util.HexStringToBytes("457743861de496c429912558a106b810b0507975a49773228aa788df40730d41"),
			util.HexStringToBytes("7688029288efc9e9a0011c960a6ed9e5466581abf3e3a6c26ee317461add619a"),
			util.HexStringToBytes("b1ae7f15836cb2286cdd4e2c37bf9bb7da0a2846d06867a429f654b2e7f383c9"),
			util.HexStringToBytes("9b74f89fa3f93e71ff2c241f32945d877281a6a50a6bf94adac002980aafe5ab"),
			util.HexStringToBytes("b3a92b5b255019bdaf754875633c2de9fec2ab03e6b8ce669d07cb5b18804638"),
			util.HexStringToBytes("b5c0b915312b9bdaedd2b86aa2d0f8feffc73a2d37668fd9010179261e25e263"),
			util.HexStringToBytes("c9d52c5cb1e557b92c84c52e7c4bfbce859408bedffc8a5560fd6e35e10b8800"),
			util.HexStringToBytes("c555bc5fc3bc096df0a0c9532f07640bfb76bfe4fc1ace214b8b228a1297a4c2"),
			util.HexStringToBytes("f9dbfafc3af3400954975da24eb325e326960a25b87fffe23eef3e7ed2fb610e"),
		}
		flagBits := make([]byte, 31)
		for i := range flagBits {
			flagBits[i] = 1
		}
		tree := NewTree(len(hashes))
		tree.PopulateTree(flagBits, hashes)
		expected := `597c4bafe3832b17cbbabe56f878f4fc2ad0f6a402cee7fa851a9cb205f87ed1`
		actual := fmt.Sprintf("%x", tree.Root())
		if actual != expected {
			t.Errorf("Expected %s, got %s", expected, actual)
		}
	})

	t.Run("test populate tree 2", func(t *testing.T) {
		hashes := [][]byte{
			util.HexStringToBytes("42f6f52f17620653dcc909e58bb352e0bd4bd1381e2955d19c00959a22122b2e"),
			util.HexStringToBytes("94c3af34b9667bf787e1c6a0a009201589755d01d02fe2877cc69b929d2418d4"),
			util.HexStringToBytes("959428d7c48113cb9149d0566bde3d46e98cf028053c522b8fa8f735241aa953"),
			util.HexStringToBytes("a9f27b99d5d108dede755710d4a1ffa2c74af70b4ca71726fa57d68454e609a2"),
			util.HexStringToBytes("62af110031e29de1efcad103b3ad4bec7bdcf6cb9c9f4afdd586981795516577"),
		}
		flagBits := make([]byte, 11)
		for i := range flagBits {
			flagBits[i] = 1
		}
		tree := NewTree(len(hashes))
		tree.PopulateTree(flagBits, hashes)
		expected := `a8e8bd023169b81bc56854137a135b97ef47a6a7237f4c6e037baed16285a5ab`
		actual := fmt.Sprintf("%x", tree.Root())
		if actual != expected {
			t.Errorf("Expected %s, got %s", expected, actual)
		}
	})
}
