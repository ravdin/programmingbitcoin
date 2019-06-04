package merkle

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/ravdin/programmingbitcoin/util"
)

type MerkleTree struct {
	Total        int
	MaxDepth     int
	Nodes        [][][]byte
	CurrentDepth int
	CurrentIndex int
}

func NewMerkleTree(total int) *MerkleTree {
	result := &MerkleTree{Total: total}
	result.MaxDepth = int(math.Ceil(math.Log2(float64(total))))
	result.Nodes = make([][][]byte, result.MaxDepth+1)
	denom := math.Pow(2.0, float64(result.MaxDepth))
	for i, _ := range result.Nodes {
		numItems := int(math.Ceil(float64(total) / denom))
		result.Nodes[i] = make([][]byte, numItems)
		denom /= 2
	}
	result.CurrentDepth = 0
	result.CurrentIndex = 0
	return result
}

func (self *MerkleTree) String() string {
	result := make([]string, len(self.Nodes))
	for depth, level := range self.Nodes {
		items := make([]string, len(level))
		for index, h := range level {
			if len(h) == 0 {
				items[index] = "None"
			} else if depth == self.CurrentDepth && index == self.CurrentIndex {
				items[index] = fmt.Sprintf("*%x.*", h[:4])
			} else {
				items[index] = fmt.Sprintf("%x...", h[:4])
			}
		}
		result[depth] = strings.Join(items, ", ")
	}
	return strings.Join(result, "\n")
}

func (self *MerkleTree) PopulateTree(flagBits []byte, hashes [][]byte) error {
	for len(self.Root()) == 0 {
		if self.isLeaf() {
			flagBits = flagBits[1:]
			self.setCurrentNode(hashes[0])
			hashes = hashes[1:]
			self.up()
			continue
		}
		leftHash := self.getLeftNode()
		if len(leftHash) == 0 {
			flagBit := flagBits[0]
			flagBits = flagBits[1:]
			if flagBit == 0 {
				self.setCurrentNode(hashes[0])
				hashes = hashes[1:]
				self.up()
			} else {
				self.left()
			}
		} else if self.rightExists() {
			rightHash := self.getRightNode()
			if len(rightHash) == 0 {
				self.right()
			} else {
				self.setCurrentNode(util.MerkleParent(leftHash, rightHash))
				self.up()
			}
		} else {
			self.setCurrentNode(util.MerkleParent(leftHash, leftHash))
			self.up()
		}
	}
	if len(hashes) != 0 {
		return errors.New(fmt.Sprintf("hashes not all consumed %d", len(hashes)))
	}
	for _, flagBit := range flagBits {
		if flagBit != 0 {
			return errors.New("flagBits not all consumed")
		}
	}
	return nil
}

func (self *MerkleTree) up() {
	self.CurrentDepth--
	self.CurrentIndex >>= 1
}

func (self *MerkleTree) left() {
	self.CurrentDepth++
	self.CurrentIndex <<= 1
}

func (self *MerkleTree) right() {
	self.CurrentDepth++
	self.CurrentIndex <<= 1
	self.CurrentIndex++
}

func (self *MerkleTree) Root() []byte {
	return self.Nodes[0][0]
}

func (self *MerkleTree) setCurrentNode(value []byte) {
	self.Nodes[self.CurrentDepth][self.CurrentIndex] = value
}

func (self *MerkleTree) getCurrentNode() []byte {
	return self.Nodes[self.CurrentDepth][self.CurrentIndex]
}

func (self *MerkleTree) getLeftNode() []byte {
	return self.Nodes[self.CurrentDepth+1][self.CurrentIndex*2]
}

func (self *MerkleTree) getRightNode() []byte {
	return self.Nodes[self.CurrentDepth+1][self.CurrentIndex*2+1]
}

func (self *MerkleTree) isLeaf() bool {
	return self.CurrentDepth == self.MaxDepth
}

func (self *MerkleTree) rightExists() bool {
	return len(self.Nodes[self.CurrentDepth+1]) > self.CurrentIndex*2+1
}
