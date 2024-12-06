package main

import (
	"fmt"
	"os"
	"strings"
	"strconv"
	"slices"
	"maps"
)

func difference(first, second int) int {
	if(first > second) {
		return first - second
	}
	return second - first
}

func unpack(input string) ([2]int, [][]rune) {
	lines := strings.Split(input, "\n")
	lines = lines[:len(lines)-1]
	guard := [2]int{0,0}
	labMap := [][]rune{}
	for lineNum, line := range lines {
		labMap = append(labMap, []rune(line))
		for charNum, char := range line {
			if (char == rune('^')) {
				guard = [2]int{charNum,lineNum}
				break
			}
		}
	}
	return guard, labMap
}

// Directions
const (
	UP int = 0
	DOWN int = 1
	LEFT int = 2
	RIGHT int = 3
)

// Modifying vectors, depending on direction
var modifiers = [4][2]int{[2]int{0,-1}, [2]int{0,1}, [2]int{-1,0}, [2]int{1,0}}

// Keeps track of distinct positions visited
var visitedPositions = make(map[string][]int)

// Finds the modifier, start index and maximum value when navigating the grid in a certain direction
func fetchSettings(labMap [][]rune, startPos [2]int, direction int) (modifier [2]int, startIndex int, maxValue int) {
	startIndicies := [4]int{startPos[1], startPos[1], startPos[0], startPos[0]}
	maxValues := [4]int{0, len(labMap)-1, 0, len(labMap[0])-1}
	return modifiers[direction], startIndicies[direction], maxValues[direction]
}

// Advances a position in a specified direction
func advance(position [2]int, direction int, multiple int) (newPosition [2]int) {
	modifier := modifiers[direction]
	return [2]int{position[0] + (modifier[0] * multiple), position[1] + (modifier[1] * multiple)}
}

// Keeps track of the loop creation spots
var loopCreationSpots = make(map[string]bool)

func rotate(direction int) int {
	nextDirection := [4]int{RIGHT, LEFT, UP, DOWN}
	return nextDirection[direction]
}

func toIndex(position [2]int) string {
	return strings.Join([]string{strconv.Itoa(position[0]),",",strconv.Itoa(position[1])}, "")
}

func nextObstacleLocation(labMap [][]rune, startPos [2]int, direction int) (loopFound bool, obstacleFound bool, obstacleLocation [2]int) {
	_, startIndex, maxValue := fetchSettings(labMap, startPos, direction)
	newPos := startPos
	for i := 0; i <= difference(maxValue, startIndex); i += 1 {
		newPos = advance(startPos, direction, i)
		newX := newPos[0]
		newY := newPos[1]
		positionIndex := toIndex(newPos)
		isLoop := slices.Contains(visitedPositions[positionIndex], direction)

		if (isLoop) {
			// Log position as visited
			//visitedPositions[positionIndex] = append(visitedPositions[positionIndex], direction)
			// Record loop
			return true, labMap[newY][newX] == rune('#'), newPos
		}
		if (labMap[newY][newX] == rune('#')) {
			return false, true, newPos
		}
		// once loop recorded, else here
		visitedPositions[positionIndex] = append(visitedPositions[positionIndex], direction)
	}
	return false, false, newPos
}

func valueByAxisOfDirection(location [2]int, direction int) int {
	if (direction == UP || direction == DOWN) {
		return location[1]
	}
	return location[0]
}

func cloneLabMap(labMap [][]rune) [][]rune {
	var mapClone [][]rune
	for _, mapLine := range labMap {
		mapClone = append(mapClone, slices.Clone(mapLine))
	}
	return mapClone
}

func printMap(labMap [][]rune) {
	stringMap := []string{}
	for _, mapLine := range labMap {
		stringMap =  append(stringMap, string(mapLine))
	}
	fmt.Println(strings.Join(stringMap, "\n"))
	fmt.Println("\n")
}

func main() {
	dat, _ := os.ReadFile("./assets/day06-input.txt")
	guard, labMap := unpack(string(dat))
//	printMap(labMap)

	obstacleFound := true
	var obstacleLocation [2]int
	guardLocation := guard
	direction := UP
	
	for obstacleFound {
		preMoveVisitedPositions := maps.Clone(visitedPositions)

		_, obstacleFound, obstacleLocation = nextObstacleLocation(labMap, guardLocation, direction)

		postMoveVisitedPositions := maps.Clone(visitedPositions)
		visitedPositions = maps.Clone(preMoveVisitedPositions)

		// Part 2 logic
		// for each spot in between the guard and the new location
		maxValue := valueByAxisOfDirection(obstacleLocation, direction)
		startIndex := valueByAxisOfDirection(guardLocation, direction)

		for i:= 1; i <= difference(maxValue, startIndex); i += 1 {
			// place obstacle in spot
			tempLabMap := cloneLabMap(labMap)
			fakeObstacle := advance(guardLocation, direction, i)
			tempLabMap[fakeObstacle[1]][fakeObstacle[0]] = rune('#')

			// run simulation & record loop if one is created
			tempObstacleFound := true
			var tempObstacleLocation [2]int
			tempGuardLocation := guardLocation
			tempDirection := direction
			
			isLoop := false
			for tempObstacleFound && !isLoop {
				isLoop, tempObstacleFound, tempObstacleLocation = nextObstacleLocation(tempLabMap, tempGuardLocation, tempDirection)
				tempGuardLocation = advance(tempObstacleLocation, rotate(rotate(tempDirection)), 1)
				tempDirection = rotate(tempDirection)

				if (isLoop) {
					loopCreationSpots[toIndex(fakeObstacle)] = true
					break
				}
			}

			// Reset visited positions
			visitedPositions = maps.Clone(preMoveVisitedPositions)
		}
		visitedPositions = maps.Clone(postMoveVisitedPositions)
		// End of Part 2 logic */

		guardLocation = advance(obstacleLocation, rotate(rotate(direction)), 1)
		direction = rotate(direction)
	}
	fmt.Println(len(visitedPositions))
	fmt.Println(len(loopCreationSpots))
}

