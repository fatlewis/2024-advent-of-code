package main

import (
	"fmt"
	"strings"
	"os"
)

type Lock [5]int
type Key [5]int

// col: the lock column this node represents
// height: the column height this tree represents
// count: the number of locks which match this node
// next: nodes recording the details of locks for the next column, depending on height
// locks: complete locks which this node holds
type LockTreeNode struct {
	Col int
	Height int
	Count int
	Next [6]*LockTreeNode
	Locks []Lock
}

func NewLockTreeNode(col int, height int) *LockTreeNode {
	newNode := &LockTreeNode{
		Col: col,
		Height: height,
		Count: 0,
		Next: [6]*LockTreeNode{},
		Locks: []Lock{},
	}
	return newNode
}

// Inserts a lock into the tree, starting at column col
func (n *LockTreeNode) Insert(lock Lock, col int) {
	n.Count++
	// If we're on the last column, add the lock
	if (len(lock)-1 == col) {
		n.Locks = append(n.Locks, lock)
	} else {
		// Else increment the columns and pass down
		if (n.Next[lock[col+1]] == nil) { n.Next[lock[col+1]] = NewLockTreeNode(col+1, lock[col+1]) }
		n.Next[lock[col+1]].Insert(lock, col+1)
	}
}

// Checks how many locks a key would fit, starting at column col
func (n *LockTreeNode) GetCount(key Key, col int) int {
	if (col == len(key)-1) { return n.Count }

	lockSum := 0
	maxLockHeight := 5-key[col+1]
	for lockIndex := range maxLockHeight+1 {
		if (n.Next[lockIndex] != nil) {
			lockSum += n.Next[lockIndex].GetCount(key, col+1)
		}
	}
	return lockSum
}

func unpack(input string) ([]Key, []Lock, [6]*LockTreeNode) {
	blocks := strings.Split(strings.Trim(input, "\n"), "\n\n")
	keys := make([]Key, 0)
	locks := make([]Lock, 0)
	lockTree := [6]*LockTreeNode{}
	for lockIndex, _ := range lockTree { lockTree[lockIndex] = NewLockTreeNode(0, lockIndex) }
	// For each block
	for _, block := range blocks {
		splitBlock := strings.Split(block, "\n")
		// Find out if it's a key or lock
		isKey := splitBlock[0] != "#####"

		// Generate object
		if (isKey) {
			key := Key{}
			for lineIndex, blockLine := range splitBlock {
				if (lineIndex == 0 || lineIndex == 6) { continue }
				splitLine := strings.Split(blockLine, "")
				for charIndex, char := range splitLine {
					if (char == "#") { key[charIndex]++ }
				}
			}
			keys = append(keys, key)
		} else {
			lock := Lock{}
			for lineIndex, blockLine := range splitBlock {
				if (lineIndex == 0 || lineIndex == 6) { continue }
				splitLine := strings.Split(blockLine, "")
				for charIndex, char := range splitLine {
					if (char == "#") { lock[charIndex]++ }
				}
			}
			locks = append(locks, lock)

			// If lock, add to lockTree
			lockTree[lock[0]].Insert(lock, 0)
		}
	}
	return keys, locks, lockTree
}

func main() {
	// unpack data
	dat, _ := os.ReadFile("./assets/day25-input.txt")
	keys, locks, lockTree := unpack(string(dat))

	fmt.Println("Keys:")
	fmt.Println(len(keys))
	// for _, key := range keys { fmt.Println(key) }
	fmt.Println("Locks:")
	fmt.Println(len(locks))
	// for _, lock := range locks { fmt.Println(lock) }

	// for each key
	combinationSum := 0
	for _, key := range keys {
		// get lock count
		maxLockHeight := 5-key[0]
		for lockHeight := range maxLockHeight+1 {
			if (lockTree[lockHeight] != nil) {
				combinationSum += lockTree[lockHeight].GetCount(key, 0)
			}
		}
	}
	// print sum
	fmt.Print("Combination Sum: ")
	fmt.Println(combinationSum)
}
