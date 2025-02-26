package main

import (
	"fmt"
	"os"
	"strings"
	"strconv"
)

func mix(a, b int) int {
	// bitwise xor a & b
	return a ^ b
}

func prune(a int) int {
	// a mod 16777216
	return a % 16777216
}

func next(number int) int {
	// Multiply by 64
	// "mix" result into number
	// "prune" number
	result := prune(mix(number, number << 6))

	// Divide by 32 (rounded down to int)
	// "mix" result into number
	// "prune" number
	result = prune(mix(result, result >> 5))

	// Multiply by 2048
	// "mix" result into number
	// "prune" number
	result = prune(mix(result, result << 11))
	return result
}

func unpack(input string) []int {
	stringNumbers := strings.Split(strings.Trim(input, "\n"), "\n")
	result := []int{}
	for _, stringNumber := range stringNumbers {
		number, _ := strconv.Atoi(stringNumber)
		result = append(result, number)
	}
	return result
}

func generate(number, times int) int {
	result := number
	for range times { result = next(result) }
	return result
}

func main() {
	dat, _ := os.ReadFile("./assets/day22-input.txt")
	secretNumbers := unpack(string(dat))

	result := 0
	for _, secretNumber := range secretNumbers {
		result += generate(secretNumber, 2000)
	}
	fmt.Println(result)

	maxCount := 0
	totalCountMap := make(map[string]int)
	for _, secretNumber := range secretNumbers {
		countMap := make(map[string]int)
		lastChanges := []int{}
		secret := secretNumber
		for range 2000 {
			prevPrice := secret % 10
			secret = next(secret)
			currPrice := secret % 10
			
			lastChanges = append(lastChanges, currPrice-prevPrice)
			if len(lastChanges) > 4 { lastChanges = lastChanges[1:] }
			if len(lastChanges) == 4 {
				countKey := fmt.Sprintf("%v", lastChanges)
				if _, exists := countMap[countKey]; !exists {
					countMap[countKey] = currPrice
				}
			}
		}
		// copy countmap to totalcountmap and update max
		for key, value := range countMap {
			totalCountMap[key] = totalCountMap[key] + value
			maxCount = max(maxCount, totalCountMap[key])
		}
	}
	fmt.Println(maxCount)
}
