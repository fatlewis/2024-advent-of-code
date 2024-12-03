package main

import(
	"fmt"
	"os"
	"strings"
	"strconv"
	"slices"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func toInt(str string) int {
	result, err := strconv.Atoi(str)
	check(err)
	return result
}

func toIntArray(strArray []string) []int {
	var intArray []int
	for _, str := range strArray {
		intArray = append(intArray, toInt(str))
	}
	return intArray
}

func difference(first, second int) int {
	if (first > second) {
		return first - second
	}
	return second - first
}

func areLevelsSafe(firstLevel, secondLevel int, isIncreasing bool) bool {
	if (firstLevel == secondLevel) {
		return false
	}
	if ((firstLevel > secondLevel) && isIncreasing) {
		return false
	}
	if ((firstLevel < secondLevel) && !isIncreasing) {
		return false
	}
	if (difference(firstLevel, secondLevel) > 3) {
		return false
	}
	return true
}

func isSafe(report []int, levelSkipped bool) bool {
	isIncreasing := report[0] < report[1]
	for index := 0; index < (len(report) - 1); index += 1 {
		levelsSafe := areLevelsSafe(report[index], report[index + 1], isIncreasing)
		if (!levelsSafe) {
			// If we've already skipped a level, return false
			if (levelSkipped) {
				return false
			}
			// if index is 1, try skipping level 0 (may flip the isIncreasing bool)
			if (isSafe(report[1:], true)) {
				return true
			}
			// Skip current level
			if (isSafe(slices.Delete(slices.Clone(report), index, index + 1), true)) {
				return true
			}
			// Skip next level
			if (isSafe(slices.Delete(slices.Clone(report), index + 1, index + 2), true)) {
				return true
			}
			return false
		}
	}
	return true
}

func main() {
	dat, err := os.ReadFile("./assets/day02-input.txt")
	check(err)

	lines := strings.Split(string(dat), "\n")
	var levels [][]int
	for _, line := range lines {
		if (len(line) == 0) {
			continue;
		}
		levels = append(levels, toIntArray(strings.Split(line, " ")))
	}
	numSafe := 0
	for _, level := range levels {
		if(isSafe(level, false)) {
			numSafe += 1
		}
	}
	fmt.Println(numSafe)
}

