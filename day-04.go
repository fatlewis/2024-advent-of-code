package main

import (
	"fmt"
	"os"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

/* Find XMAS algorithm
**
** - For each line. Find X or S
** - If S, reverse word, try word L-R, try word U-D, try word L-R diagonal (up and down)
*/

func canFitHorizontalLR(word string, width int, xIndex int) bool {
	return len(word) - 1 + xIndex < width
}

func canFitVerticalUD(word string, height int, yIndex int) bool {
	return len(word) - 1 + yIndex < height
}

func reverse(s string) string {
    runes := []rune(s)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}

func isWordPresentLR(word string, grid []string, x int, y int) bool {
	if (!canFitHorizontalLR(word, len(grid[y]), x)) {
		return false
	}
	for i := 0; i < len(word); i += 1 {
		if (word[i] != grid[y][x + i]) {
			return false
		}
	}
	return true
}

func isWordPresentUD(word string, grid []string, x int, y int) bool {
	if (!canFitVerticalUD(word, len(grid), y)) {
		return false
	}
	for i := 0; i < len(word); i += 1 {
		if (word[i] != grid[y + i][x]) {
			return false
		}
	}
	return true
}

func isWordPresentDiagonalUp(word string, grid []string, x int, y int) bool {
	if (y < (len(word) - 1) || !canFitHorizontalLR(word, len(grid[y]), x)) {
		return false
	}
	for i := 0; i < len(word); i += 1 {
		if (word[i] != grid[y - i][x + i]) {
			return false
		}
	}
	return true
}

func isWordPresentDiagonalDown(word string, grid []string, x int, y int) bool {
	if (!canFitVerticalUD(word, len(grid), y) || !canFitHorizontalLR(word, len(grid[y]), x)) {
		return false
	}
	for i:= 0; i < len(word); i += 1 {
		if (word[i] != grid[y + i][x + i]) {
			return false
		}
	}
	return true
}

func main() {
	dat, err := os.ReadFile("./assets/day04-input.txt")
	check(err)

	grid := strings.Split(string(dat), "\n")
	grid = grid[:len(grid)-1]
	word := "XMAS"
	reverseWord := reverse(word)
	occurrences := 0
	for yIndex, line := range grid {
		for xIndex, character := range line {
			if (character == rune(word[0])) {
				if (isWordPresentLR(word, grid, xIndex, yIndex)) {
					occurrences += 1
				}
				if (isWordPresentUD(word, grid, xIndex, yIndex)) {
					occurrences += 1
				}
				if (isWordPresentDiagonalUp(word, grid, xIndex, yIndex)) {
					occurrences += 1
				}
				if (isWordPresentDiagonalDown(word, grid, xIndex, yIndex)) {
					occurrences += 1
				}
			}
			if (character == rune(reverseWord[0])) {
				if (isWordPresentLR(reverseWord, grid, xIndex, yIndex)) {
					occurrences += 1
				}
				if (isWordPresentUD(reverseWord, grid, xIndex, yIndex)) {
					occurrences += 1
				}
				if (isWordPresentDiagonalUp(reverseWord, grid, xIndex, yIndex)) {
					occurrences += 1
				}
				if (isWordPresentDiagonalDown(reverseWord, grid, xIndex, yIndex)) {
					occurrences += 1
				}
			}
		}
	}
	fmt.Println(occurrences)

	// Finding X-MASs
	word = "MAS"
	reverseWord = reverse(word)
	occurrences = 0
	for yIndex, line := range grid {
		for xIndex, character := range line {
			if (character == rune(word[0])) {
				if (isWordPresentDiagonalDown(word, grid, xIndex, yIndex)) {
					if (isWordPresentDiagonalUp(word, grid, xIndex, yIndex + 2) || isWordPresentDiagonalUp(reverseWord, grid, xIndex, yIndex + 2)) {
						occurrences += 1
					}
				}
			}
			if (character == rune(reverseWord[0])) {
				if (isWordPresentDiagonalDown(reverseWord, grid, xIndex, yIndex)) {
					if (isWordPresentDiagonalUp(word, grid, xIndex, yIndex + 2) || isWordPresentDiagonalUp(reverseWord, grid, xIndex, yIndex + 2)) {
						occurrences += 1
					}
				}
			}
		}
	}
	fmt.Println(occurrences)
}

