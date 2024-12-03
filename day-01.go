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

func difference(first, second int) int {
	if (first > second) {
		return first - second
	}
	return second - first
}

func main() {
	dat, err := os.ReadFile("./assets/day01-input.txt")
	check(err)

	var groupOneLocations, groupTwoLocations []int

	lines := strings.Split(string(dat), "\n")
	for _, line := range lines {
		if (len(line) == 0) {
			continue;
		}
		before, after, found := strings.Cut(line, "   ")
		if (!found) {
			panic(fmt.Sprintf("Could not split line: %s", line))
		}

		groupOneLocations = append(groupOneLocations, toInt(before))
		groupTwoLocations = append(groupTwoLocations, toInt(after))
	}

	slices.Sort(groupOneLocations)
	slices.Sort(groupTwoLocations)
	
	var distanceSum int
	for index, groupOneLocation := range groupOneLocations {
		distanceSum = distanceSum + difference(groupOneLocation, groupTwoLocations[index])
	}
	fmt.Println(distanceSum)

	// convert second locations to weights table
	weights := make(map[int]int)
	for _, groupTwoLocation := range groupTwoLocations {
		_, ok := weights[groupTwoLocation]
		if (!ok) {
			weights[groupTwoLocation] = 0
		}
		weights[groupTwoLocation] = weights[groupTwoLocation] + 1
	}
	var weightSum int
	for _, groupOneLocation := range groupOneLocations {
		_, ok := weights[groupOneLocation]
		if (ok) {
			weightSum = weightSum + (groupOneLocation * weights[groupOneLocation])
		}
	}
	fmt.Println(weightSum)
}

