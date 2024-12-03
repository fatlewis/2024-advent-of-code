package main

import (
	"os"
	"regexp"
	"strconv"
	"fmt"
	"strings"
)

func check(e error) {
	if (e != nil) {
		panic(e)
	}
}

func toInt(str string) int {
	result, err := strconv.Atoi(str)
	check(err)
	return result
}

func main() {
	dat, err := os.ReadFile("./assets/day03-input.txt")
	check(err)

	instructions := regexp.MustCompile(`mul\(([0-9]{1,3}),([0-9]{1,3})\)`).FindAllStringSubmatch(string(dat), -1)
	var result int
	for _, instruction := range instructions {
		result += toInt(instruction[1]) * toInt(instruction[2])
	}
	fmt.Println(result)

	fullInstructions := regexp.MustCompile(`(?:mul\((?:([0-9]{1,3}),([0-9]{1,3})))\)|do\(\)|don't\(\)`).FindAllStringSubmatch(string(dat), -1)
	var fullResult int
	enabled := true
	for _, instruction := range fullInstructions {
		command := strings.Split(instruction[0], "(")[0]
		switch command {
		case "do":
			enabled = true
		case "don't":
			enabled = false
		case "mul":
			if (enabled == true) {
				fullResult += toInt(instruction[1]) * toInt(instruction[2])
			}
		}
	}
	fmt.Println(fullResult)
}

