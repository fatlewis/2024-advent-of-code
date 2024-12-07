package main

import (
	"os"
	"fmt"
	"strings"
	"strconv"
	"time"
)

func strToInt(str string) int {
	result, _ := strconv.Atoi(str)
	return result
}

func unpack(input string) (testValues []int, numbers [][]int) {
	lines := strings.Split(input, "\n")
	lines = lines[:len(lines)-1]

	for index, line := range lines {
		splitLine := strings.Split(line, ": ")
		testValues = append(testValues, strToInt(splitLine[0]))

		numbers = append(numbers, []int{})
		for _, lineNumber := range strings.Split(splitLine[1], " ") {
			numbers[index] = append(numbers[index], strToInt(lineNumber))
		}
	}
	return
}

var operators []string = []string{"+", "*", "||"}

var isPartTwo bool = false

func canSolve(result int, current int, remaining []int) bool {
	if (len(remaining) == 0) {
		return result == current
	}
	if (current > result) {
		return false
	}

	nextValue := remaining[0]
	for _, operator := range operators {
		switch operator {
			case "+":
				if(canSolve(result, current+nextValue, remaining[1:])) {
					return true
				}
			case "*":
				if(canSolve(result, current*nextValue, remaining[1:])) {
					return true
				}
			case "||":
				if(isPartTwo && canSolve(result, strToInt(fmt.Sprintf("%d%d", current, nextValue)), remaining[1:])) {
					return true
				}
		}
	}
	return false
}

func main() {
	dat, _ := os.ReadFile("./assets/day07-input.txt")
	testValues, numbers := unpack(string(dat))

	startTime := time.Now()
	solveableSum := 0
	for index, testValue := range testValues {
		if(canSolve(testValue, numbers[index][0], numbers[index][1:])) {
			solveableSum += testValue
		}
	}
	fmt.Println(time.Since(startTime))
	fmt.Println(solveableSum)

	startTime = time.Now()
	solveableSum = 0
	isPartTwo = true
	for index, testValue := range testValues {
		if(canSolve(testValue, numbers[index][0], numbers[index][1:])) {
			solveableSum += testValue
		}
	}
	fmt.Println(time.Since(startTime))
	fmt.Println(solveableSum)
}

