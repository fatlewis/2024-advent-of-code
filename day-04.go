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

func boolToInt(boolean bool) int {
	if (boolean) {
		return 1
	}
	return 0
}

func isWordPresent(word string, grid []string, location [2]int, vector [2]int) bool {
	fitsHorizontal := vector[0] <= 0 || canFitHorizontalLR(word, len(grid[location[1]]), location[0])
	fitsVerticalDown := vector[1] <= 0 || canFitVerticalUD(word, len(grid), location[1])
	fitsVerticalUp := vector[1] >= 0 || location[1] >= (len(word) -1)
	if(!fitsHorizontal || !fitsVerticalDown || !fitsVerticalUp) {
		return false
	}
	for i:= 0; i < len(word); i += 1 {
		if (word[i] != grid[location[1] + (i * vector[1])][location[0] + (i * vector[0])]) {
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
	vectors := [4][2]int{[2]int{1,0}, [2]int{0,1}, [2]int{1,-1}, [2]int{1,1}}
	for yIndex, line := range grid {
		for xIndex, character := range line {
			location := [2]int{xIndex,yIndex}
			if (character == rune(word[0])) {
				for _, vector := range vectors {
					occurrences += boolToInt(isWordPresent(word, grid, location, vector))
				}
			}
			if (character == rune(reverseWord[0])) {
				for _, vector := range vectors {
					occurrences += boolToInt(isWordPresent(reverseWord, grid, location, vector))
				}
			}
		}
	}
	fmt.Println(occurrences)

	// Finding X-MASs
	word = "MAS"
	occurrences = 0
	for yIndex, line := range grid {
		for xIndex, character := range line {
			location := [2]int{xIndex,yIndex}
			var testWord string
			if (character == rune(word[0])) {
				testWord = word
			}
			if (character == rune(reverse(word)[0])) {
				testWord = reverse(word)
			}
			if (testWord != "") {
				if (isWordPresent(testWord, grid, location, [2]int{1,1})) {
					newLocation := [2]int{location[0],location[1]+2}
					wordCrosses := isWordPresent(testWord, grid, newLocation, [2]int{1,-1})
					reverseWordCrosses := isWordPresent(reverse(testWord), grid, newLocation, [2]int{1,-1})
					occurrences += boolToInt(wordCrosses || reverseWordCrosses)
				}
			}
		}
	}
	fmt.Println(occurrences)
}

