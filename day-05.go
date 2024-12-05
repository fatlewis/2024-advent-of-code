package main

import (
	"fmt"
	"os"
	"strings"
	"strconv"
	"slices"
)

func strToInt(str string) int {
	result, _ := strconv.Atoi(str)
	return result
}

func unpack(input string) (map[int][]int, [][]string) {
	sections := strings.Split(input, "\n\n")
	ruleSection, manualsSection := sections[0], sections[1]
	
	rules := make(map[int][]int)
	rulesSlice := strings.Split(ruleSection, "\n")
	for _, rule := range rulesSlice {
		splitRule := strings.Split(rule, "|")
		key := strToInt(splitRule[0])
		rules[key] = append(rules[key], strToInt(splitRule[1]))
	}

	manuals := make([][]string, 0)
	manualsSlice := strings.Split(manualsSection, "\n")
	manualsSlice = manualsSlice[:len(manualsSlice)-1]
	for _, manual := range manualsSlice {
		manuals = append(manuals, strings.Split(manual, ","))
	}
	return rules, manuals
}

func earliestViolation(manual map[string]int, earlyPage string, laterPages []int) int {
	earliestViolation := -1
	for _, laterPage := range laterPages {
		_, inManual := manual[strconv.Itoa(laterPage)]
		if (inManual && manual[strconv.Itoa(laterPage)] <= manual[earlyPage]) {
			if (earliestViolation < 0 || manual[strconv.Itoa(laterPage)] < earliestViolation) {
				earliestViolation = manual[strconv.Itoa(laterPage)]
			}
		}
	}
	return earliestViolation
}

func toMap(manual []string) map[string]int {
	manualMap := make(map[string]int)
	manualMap["middle"] = strToInt(manual[len(manual)/2])
	for pageNum, page := range manual {
		manualMap[page] = pageNum
	}
	return manualMap
}

func swap(manual []string, first int, second int) []string {
	return slices.Concat(
		manual[0:first],
		[]string{manual[second]},
		manual[first:second],
		manual[second+1:],
	)
}

func main() {
	dat, _ := os.ReadFile("./assets/day05-input.txt")
	rules, manuals := unpack(string(dat))
	
	validSum := 0
	for _, manual := range manuals {
		manualMap := toMap(manual)
		isValid := true
		for _, page := range manual {
			if (earliestViolation(manualMap, page, rules[strToInt(page)]) >= 0) {
				isValid = false
				break
			}
		}
		if (isValid) {
			validSum += manualMap["middle"]
		}
	}
	fmt.Println(validSum)

	validSum = 0
	for _, manual := range manuals {
		manualMap := toMap(manual)
		isValid := true
		for pageNum, page := range manual {
			swapIndex := earliestViolation(manualMap, page, rules[strToInt(page)])
			if (swapIndex >= 0) {
				isValid = false
				manual = swap(manual, swapIndex, pageNum)
				manualMap = toMap(manual)
			}
		}
		if (!isValid) {
			validSum += manualMap["middle"]
		}
	}
	fmt.Println(validSum)
}

