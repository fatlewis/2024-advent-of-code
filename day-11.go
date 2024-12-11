package main

import (
	"fmt"
	"os"
	"strings"
	"strconv"
	"time"
	"maps"
)

func toInt(str string) int {
	result, _ := strconv.Atoi(str)
	return result
}

func updateStone(stone int) []int {
	if (stone == 0) {
		return []int{1}
	}
	stoneString := strconv.Itoa(stone)
	if (len(stoneString) % 2 == 0) {
		return []int{
			toInt(stoneString[0:len(stoneString)/2]),
			toInt(stoneString[len(stoneString)/2:]),
		}
	}
	// check for overflow? maybe use float?
	return []int{stone * 2024}
}

func updateStones(stoneMap map[int]int) map[int]int {
	newMap := make(map[int]int)
	for stone := range maps.Keys(stoneMap) {
		newStones := updateStone(stone)
		for _, newStone := range newStones {
			newMap[newStone] += stoneMap[stone]
		}
	}
	return newMap
}

func sumValues(stoneMap map[int]int) int {
	result := 0
	for stoneValue := range maps.Values(stoneMap) {
		result += stoneValue
	}
	return result
}

func main() {
	dat, _ := os.ReadFile("./assets/day11-input.txt")
	stoneStrings := strings.Split(strings.Trim(string(dat), "\n"), " ")

	stoneMap := make(map[int]int)
	for _, stoneString := range stoneStrings {
		stoneInt := toInt(stoneString)
		stoneMap[stoneInt] += 1
	}

	start := time.Now()
	for i := 0; i < 25; i += 1 {
		stoneMap = updateStones(stoneMap)
	}
	fmt.Println(sumValues(stoneMap))
	fmt.Println(time.Since(start))

	clear(stoneMap)
	for _, stoneString := range stoneStrings {
		stoneInt := toInt(stoneString)
		stoneMap[stoneInt] += 1
	}

	start = time.Now()
	for i := 0; i < 75; i += 1 {
		stoneMap = updateStones(stoneMap)
	}
	fmt.Println(sumValues(stoneMap))
	fmt.Println(time.Since(start))
}

