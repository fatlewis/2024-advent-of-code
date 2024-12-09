package main

import (
	"os"
	"fmt"
	"strings"
	"strconv"
	"slices"
)

type location struct {
	x int
	y int
}

func (loc *location) toString() string {
	return strings.Join([]string{"[",strconv.Itoa(loc.x),",",strconv.Itoa(loc.y),"]"}, "")
}

type cityMap struct {
	grid []string
	antennae map[rune][]*location
}

func (city *cityMap) contains(l *location) bool {
	if (l.x < 0 || l.y <0) {
		return false
	}
	if (l.x >= len(city.grid[0]) || l.y >= len(city.grid)) {
		return false
	}
	return true
}

func difference(first int, second int) int {
	if (first > second) {
		return first - second
	}
	return second - first
}

func getAntinodes(first *location, second *location) [2]*location {
	xDiff := difference(first.x, second.x)
	yDiff := difference(first.y, second.y)
	if (first.x < second.x) {
		xDiff = -xDiff
	}
	if (first.y < second.y) {
		yDiff = -yDiff
	}
	return [2]*location{
		&location{x: first.x + xDiff, y: first.y + yDiff},
		&location{x: second.x - xDiff, y: second.y - yDiff},
	}
}

func getAntinodesUsingVector(vector [2]int, startLocation *location, antennaeMap *cityMap) []*location {
	antinodes := []*location{}
	currentAntinode := startLocation
	for antennaeMap.contains(currentAntinode) {
		antinodes = append(antinodes, currentAntinode)
		currentAntinode = &location{currentAntinode.x + vector[0], currentAntinode.y + vector[1]}
	}
	return antinodes
}

func reverse(vector [2]int) [2]int {
	return [2]int{vector[0]*-1, vector[1]*-1}
}

func getExpandedAntinodes(first *location, second *location, antennaeMap *cityMap) []*location {
	// keep getting antinodes until the new one generated isn't contained
	antinodeVector := [2]int{first.x - second.x, first.y - second.y}
	// antinote generation function:
	expandedAntinodes := []*location{}
	expandedAntinodes = slices.Concat(expandedAntinodes, getAntinodesUsingVector(antinodeVector, first, antennaeMap))
	expandedAntinodes = slices.Concat(expandedAntinodes, getAntinodesUsingVector(reverse(antinodeVector), first, antennaeMap))
	return expandedAntinodes
}

func unpack(input string) *cityMap {
	grid := strings.Split(input, "\n")
	grid = grid[:len(grid)-1]

	antennae := make(map[rune][]*location)
	for lineNum, line := range grid {
		for charNum, char := range line {
			if (char != '.') {
				antennae[char] = append(antennae[char], &location{x: charNum, y: lineNum})
			}
		}
	}
	return &cityMap{grid, antennae}
}

func main() {
	dat, _ := os.ReadFile("./assets/day08-input.txt")
	antennaeMap := unpack(string(dat))
//	fmt.Println(antennaeMap)

	validAntinodeLocations := make(map[string]bool)
	for _, antennaeSet := range antennaeMap.antennae {
		for index, antenna := range antennaeSet {
			for i := index + 1; i < len(antennaeSet); i += 1 {
				antinodes := getAntinodes(antenna, antennaeSet[i])
				if (antennaeMap.contains(antinodes[0])) {
					validAntinodeLocations[antinodes[0].toString()] = true
				}
				if (antennaeMap.contains(antinodes[1])) {
					validAntinodeLocations[antinodes[1].toString()] = true
				}
			}
		}
	}
	fmt.Println(len(validAntinodeLocations))

	newValidAntinodes := make(map[string]bool)
	for _, antennaeSet := range antennaeMap.antennae {
		for index, antenna := range antennaeSet {
			for i := index + 1; i < len(antennaeSet); i += 1 {
				antinodes := getExpandedAntinodes(antenna, antennaeSet[i], antennaeMap)
				for _, antinode := range antinodes {
					newValidAntinodes[antinode.toString()] = true
				}
			}
		}
	}
	fmt.Println(len(newValidAntinodes))
}
