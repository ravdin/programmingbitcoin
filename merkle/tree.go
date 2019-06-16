package merkle

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/ravdin/programmingbitcoin/util"
)

// Tree represents a merkle tree.
type Tree struct {
	Total        int
	MaxDepth     int
	Nodes        [][][]byte
	CurrentDepth int
	CurrentIndex int
}

// NewTree initializes a new merkle tree.
// total: number of nodes in the tree.
func NewTree(total int) *Tree {
	result := &Tree{Total: total}
	result.MaxDepth = int(math.Ceil(math.Log2(float64(total))))
	result.Nodes = make([][][]byte, result.MaxDepth+1)
	denom := math.Pow(2.0, float64(result.MaxDepth))
	for i := range result.Nodes {
		numItems := int(math.Ceil(float64(total) / denom))
		result.Nodes[i] = make([][]byte, numItems)
		denom /= 2
	}
	result.CurrentDepth = 0
	result.CurrentIndex = 0
	return result
}

func (tree *Tree) String() string {
	result := make([]string, len(tree.Nodes))
	for depth, level := range tree.Nodes {
		items := make([]string, len(level))
		for index, h := range level {
			if len(h) == 0 {
				items[index] = "None"
			} else if depth == tree.CurrentDepth && index == tree.CurrentIndex {
				items[index] = fmt.Sprintf("*%x.*", h[:4])
			} else {
				items[index] = fmt.Sprintf("%x...", h[:4])
			}
		}
		result[depth] = strings.Join(items, ", ")
	}
	return strings.Join(result, "\n")
}

// PopulateTree populates the tree.
func (tree *Tree) PopulateTree(flagBits []byte, hashes [][]byte) error {
	for len(tree.Root()) == 0 {
		if tree.isLeaf() {
			flagBits = flagBits[1:]
			tree.setCurrentNode(hashes[0])
			hashes = hashes[1:]
			tree.up()
			continue
		}
		leftHash := tree.getLeftNode()
		if len(leftHash) == 0 {
			flagBit := flagBits[0]
			flagBits = flagBits[1:]
			if flagBit == 0 {
				tree.setCurrentNode(hashes[0])
				hashes = hashes[1:]
				tree.up()
			} else {
				tree.left()
			}
		} else if tree.rightExists() {
			rightHash := tree.getRightNode()
			if len(rightHash) == 0 {
				tree.right()
			} else {
				tree.setCurrentNode(util.MerkleParent(leftHash, rightHash))
				tree.up()
			}
		} else {
			tree.setCurrentNode(util.MerkleParent(leftHash, leftHash))
			tree.up()
		}
	}
	if len(hashes) != 0 {
		return fmt.Errorf("hashes not all consumed %d", len(hashes))
	}
	for _, flagBit := range flagBits {
		if flagBit != 0 {
			return errors.New("flagBits not all consumed")
		}
	}
	return nil
}

func (tree *Tree) up() {
	tree.CurrentDepth--
	tree.CurrentIndex >>= 1
}

func (tree *Tree) left() {
	tree.CurrentDepth++
	tree.CurrentIndex <<= 1
}

func (tree *Tree) right() {
	tree.CurrentDepth++
	tree.CurrentIndex <<= 1
	tree.CurrentIndex++
}

// Root returns the root node.
func (tree *Tree) Root() []byte {
	return tree.Nodes[0][0]
}

func (tree *Tree) setCurrentNode(value []byte) {
	tree.Nodes[tree.CurrentDepth][tree.CurrentIndex] = value
}

func (tree *Tree) getCurrentNode() []byte {
	return tree.Nodes[tree.CurrentDepth][tree.CurrentIndex]
}

func (tree *Tree) getLeftNode() []byte {
	return tree.Nodes[tree.CurrentDepth+1][tree.CurrentIndex*2]
}

func (tree *Tree) getRightNode() []byte {
	return tree.Nodes[tree.CurrentDepth+1][tree.CurrentIndex*2+1]
}

func (tree *Tree) isLeaf() bool {
	return tree.CurrentDepth == tree.MaxDepth
}

func (tree *Tree) rightExists() bool {
	return len(tree.Nodes[tree.CurrentDepth+1]) > tree.CurrentIndex*2+1
}
