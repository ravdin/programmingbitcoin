package ecc

import (
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
			secret := new(big.Int)
			x := new(big.Int)
			y := new(big.Int)
			secret.SetString(item[0], 16)
			x.SetString(item[1], 16)
			y.SetString(item[2], 16)
			expected, _ := NewS256Point(x, y)
			actual := G.Rmul(secret)
			if !actual.Eq(expected) {
				t.Errorf("Expected %v, got %v", expected, actual)
			}
		}
	})

	// TODO: test_verify
	// TODO: test_sec
	// TODO: test_address
}
