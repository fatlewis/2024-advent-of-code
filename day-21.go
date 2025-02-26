package main

import (
	"fmt"
	"os"
	"strconv"
	"math"
	"strings"
)

type Keypad struct {
	parent *Keypad
	cache map[string]int64
	location byte
	layout [][]byte

	identifier string
}

func (keypad *Keypad) LocationCoords() []int {
	for i := 0; i < len(keypad.layout); i++ {
		for j := 0; j < len(keypad.layout[0]); j++ {
			if keypad.layout[i][j] == keypad.location { return []int{i,j} }
		}
	}
	return nil
}

func (keypad *Keypad) CostToPress(key byte) int64 {
	cacheKey := string(keypad.location) + string(key)
	//fmt.Println(keypad.identifier, string(key), keypad.cache[cacheKey])
	if result, exists := keypad.cache[cacheKey]; exists { return result }
	routes := dfs(keypad.layout, keypad.LocationCoords(), key)

	var minCost int64 = math.MaxInt64
	minRoute := ""
	for _, route := range routes {
		var routeCost int64 = 0
		for _, routeChar := range route {
			// Base case, assume parent is a human if it is absent
			if keypad.parent == nil {
				routeCost += 1
				continue
			}
			routeCost += keypad.parent.CostToPress(byte(routeChar))
			keypad.parent.Press(byte(routeChar))
		}
		if routeCost <= minCost { minRoute = route }
		minCost = min(minCost, routeCost)
	}
	if keypad.identifier == "" { // Change to identifier to log
		fmt.Println(keypad.identifier, string(key), minRoute, minCost)
	}
	keypad.cache[cacheKey] = minCost
	return keypad.cache[cacheKey]
}

func (dk *Keypad) Press(key byte) { dk.location = key }

func adjacents(grid [][]byte, location []int) map[byte][]int {
	results := map[byte][]int{}
	if location[0] > 0 {
		results['^'] = []int{location[0] - 1, location[1]}
	}
	if location[0] < len(grid)-1 {
		results['v'] = []int{location[0] + 1, location[1]}
	}
	if location[1] > 0 {
		results['<'] = []int{location[0], location[1] - 1}
	}
	if location[1] < len(grid[0])-1 {
		results['>'] = []int{location[0], location[1] + 1}
	}
	return results
}

// Find possible routes using depth-first search
func dfs(grid [][]byte, location []int, query byte) []string {
	if grid[location[0]][location[1]] == query { return []string{"A"} }
	if grid[location[0]][location[1]] == '#' { return []string{} }

	// find all routes from start to end, return keys which must be pressed
	currChar := grid[location[0]][location[1]]
	grid[location[0]][location[1]] = '#'

	results := []string{}
	nextLocations := adjacents(grid, location)
	for key, nextLocation := range nextLocations {
		routes := dfs(grid, nextLocation, query)
		for _, route := range routes {
			results = append(results, string(key) + route)
		}
	}

	grid[location[0]][location[1]] = currChar
	return results
}

func MakeKeypad(layout []string, parent *Keypad, identifier string) *Keypad {
	byteLayout := [][]byte{}
	for _, layoutLine := range layout { byteLayout = append(byteLayout, []byte(layoutLine)) }
	return &Keypad{
		location: 'A',
		cache: make(map[string]int64),
		layout: byteLayout,
		parent: parent,

		identifier: identifier,
	}
}

func unpack(input string) []string {
	return strings.Split(strings.Trim(input, "\n"), "\n")
}

func makeKeypadStack(height int) *Keypad {
	doorPad := MakeKeypad([]string{"789","456","123","#0A"}, nil, "door")
	prevPad := doorPad
	for roboNum := range height {
		robotPad := MakeKeypad(
			[]string{"#^A","<v>"},
			nil,
			fmt.Sprintf("robot-%d", roboNum),
		)
		prevPad.parent = robotPad
		prevPad = robotPad
	}
	return doorPad
}

func complexitySum(codes []string, keypad *Keypad) int64 {
	result := int64(0)
	for _, code := range codes {
		// get numeric part of the code
		numericPart, _ := strconv.Atoi(code[:3])
		// calculate complexity by pushing buttons
		var shortestSequence int64 = 0
		for _, digit := range code {
			shortestSequence += keypad.CostToPress(byte(digit))
			keypad.Press(byte(digit))
		}
		result += int64(numericPart) * shortestSequence
	}
	return result
}

func main() {
	dat, _ := os.ReadFile("./assets/day21-input.txt")
	codes := unpack(string(dat))

	fmt.Println(complexitySum(codes, makeKeypadStack(2)))
	fmt.Println(complexitySum(codes, makeKeypadStack(25)))
}

