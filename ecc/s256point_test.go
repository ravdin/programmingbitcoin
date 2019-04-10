package ecc

import (
	"bytes"
	"github.com/ravdin/programmingbitcoin/util"
	"math/big"
	"testing"
)

func TestS256Point(t *testing.T) {
	t.Run("Test Order", func(t *testing.T) {
		point := G.Rmul(N)
		if point.X != nil || point.Y != nil {
			t.Errorf("Expected point at infinity, got %v", point)
		}
	})

	t.Run("Test Public Point", func(t *testing.T) {
		points := [][]string{
			{"7", "5cbdf0646e5db4eaa398f365f2ea7a0e3d419b7e0330e39ce92bddedcac4f9bc", "6aebca40ba255960a3178d6d861a54dba813d0b813fde7b5a5082628087264da"},
			{"5cd", "c982196a7466fbbbb0e27a940b6af926c1a74d5ad07128c82824a11b5398afda", "7a91f9eae64438afb9ce6448a1c133db2d8fb9254e4546b6f001637d50901f55"},
			{"100000000000000000000000000000000", "8f68b9d2f63b5f339239c1ad981f162ee88c5678723ea3351b7b444c9ec4c0da", "662a9f2dba063986de1d90c2b6be215dbbea2cfe95510bfdf23cbf79501fff82"},
			{"1000000000000000000000000000000000000000000000000000080000000", "9577ff57c8234558f293df502ca4f09cbc65a6572c842b39b366f21717945116", "10b49c67fa9365ad7b90dab070be339a1daf9052373ec30ffae4f72d5e66d053"},
		}
		for _, item := range points {
			secret := util.HexStringToBigInt(item[0])
			x := util.HexStringToBigInt(item[1])
			y := util.HexStringToBigInt(item[2])
			expected, _ := NewS256Point(x, y)
			actual := G.Rmul(secret)
			if !actual.Eq(expected) {
				t.Errorf("Expected %v, got %v", expected, actual)
			}
		}
	})

	t.Run("Test Verify", func(t *testing.T) {
		px := util.HexStringToBigInt("887387e452b8eacc4acfde10d9aaf7f6d9a0f975aabb10d006e4da568744d06c")
		py := util.HexStringToBigInt("61de6d95231cd89026e286df3b6ae4a894a3378e393e93a0f45b666329a0ae34")
		point, _ := NewS256Point(px, py)
		tests := [][]string{
			{
				"ec208baa0fc1c19f708a9ca96fdeff3ac3f230bb4a7ba4aede4942ad003c0f60",
				"ac8d1c87e51d0d441be8b3dd5b05c8795b48875dffe00b7ffcfac23010d3a395",
				"68342ceff8935ededd102dd876ffd6ba72d6a427a3edb13d26eb0781cb423c4",
			},
			{
				"7c076ff316692a3d7eb3c3bb0f8b1488cf72e1afcd929e29307032997a838a3d",
				"eff69ef2b1bd93a66ed5219add4fb51e11a840f404876325a1e8ffe0529a2c",
				"c7207fee197d27c618aea621406f6bf5ef6fca38681d82b2f06fddbdce6feab6",
			},
		}
		for _, test := range tests {
			z := util.HexStringToBigInt(test[0])
			r := util.HexStringToBigInt(test[1])
			s := util.HexStringToBigInt(test[2])
			if !point.Verify(z, NewSignature(r, s)) {
				t.Errorf("Verify failed!")
			}
		}
	})

	t.Run("Test SEC", func(t *testing.T) {
		tests := map[int64][]string{
			997002999: {
				"049d5ca49670cbe4c3bfa84c96a8c87df086c6ea6a24ba6b809c9de234496808d56fa15cc7f3d38cda98dee2419f415b7513dde1301f8643cd9245aea7f3f911f9",
				"039d5ca49670cbe4c3bfa84c96a8c87df086c6ea6a24ba6b809c9de234496808d5",
			},
			123: {
				"04a598a8030da6d86c6bc7f2f5144ea549d28211ea58faa70ebf4c1e665c1fe9b5204b5d6f84822c307e4b4a7140737aec23fc63b65b35f86a10026dbd2d864e6b",
				"03a598a8030da6d86c6bc7f2f5144ea549d28211ea58faa70ebf4c1e665c1fe9b5",
			},
			42424242: {
				"04aee2e7d843f7430097859e2bc603abcc3274ff8169c1a469fee0f20614066f8e21ec53f40efac47ac1c5211b2123527e0e9b57ede790c4da1e72c91fb7da54a3",
				"03aee2e7d843f7430097859e2bc603abcc3274ff8169c1a469fee0f20614066f8e",
			},
		}
		for coefficient, sec := range tests {
			point := G.Rmul(big.NewInt(coefficient))
			uncompressed := sec[0]
			compressed := sec[1]
			if !bytes.Equal(point.Sec(false), util.HexStringToBytes(uncompressed)) {
				t.Errorf("Uncompressed SEC failed!")
			}
			if !bytes.Equal(point.Sec(true), util.HexStringToBytes(compressed)) {
				t.Errorf("Uncompressed SEC failed!")
			}
		}
	})

	t.Run("Test Address", func(t *testing.T) {
		tests := map[int64][]string{
			700227072: {
				"148dY81A9BmdpMhvYEVznrM45kWN32vSCN",
				"mieaqB68xDCtbUBYFoUNcmZNwk74xcBfTP",
			},
			321: {
				"1S6g2xBJSED7Qr9CYZib5f4PYVhHZiVfj",
				"mfx3y63A7TfTtXKkv7Y6QzsPFY6QCBCXiP",
			},
			4242424242: {
				"1226JSptcStqn4Yq9aAmNXdwdc2ixuH9nb",
				"mgY3bVusRUL6ZB2Ss999CSrGVbdRwVpM8s",
			},
		}
		compressed := true
		for secret, addresses := range tests {
			point := G.Rmul(big.NewInt(secret))
			mainnetExpected := addresses[0]
			testnetExpected := addresses[1]
			mainnetActual := point.Address(compressed, false)
			testnetActual := point.Address(compressed, true)
			if mainnetActual != mainnetExpected {
				t.Errorf("mainnet address failed, expected '%v', got '%v'", mainnetExpected, mainnetActual)
			}
			if testnetActual != testnetExpected {
				t.Errorf("testnet address failed, expected '%v', got '%v'", testnetExpected, testnetActual)
			}
			// The first test is compressed and the rest are uncompressed.
			compressed = false
		}
	})
}
