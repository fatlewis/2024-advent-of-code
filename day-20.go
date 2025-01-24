package main

import (
	"fmt"
	"os"
	"strings"
	. "advent-of-code-2024/utils"
)

func unpack(input string) (maze Maze, start *Coords, end *Coords) {
	maze = make(Maze)
	mazeLines := strings.Split(strings.Trim(input, "\n"), "\n")
	for lineIndex, line := range mazeLines {
		for charIndex, char := range line {
			if char == '#' { continue }
			location := &Coords{X: charIndex, Y: lineIndex}
			maze.Add(location)

			if char == 'S' { start = location }
			if char == 'E' { end = location }
		}
	}
	return maze, start, end
}

// Returns a map of coords according to the magnitude of displacement
func cheatSquares(node *Coords, maxDisplacement int) map[int][]*Coords {
	result := make(map[int][]*Coords)
	for i := 0; i <= maxDisplacement; i++ {
		// x displacement + i
		// y displacement <= 3 - i
		for j := 0; j <= (maxDisplacement - i); j++ {
			displacement := i + j
			if displacement < 2 { continue }
			squares, ok := result[displacement]
			if !ok { squares = []*Coords{} }

			squares = append(squares, &Coords{X: node.X + i, Y: node.Y + j})
			if i != 0 { squares = append(squares, &Coords{X: node.X - i, Y: node.Y + j}) }
			if j != 0 { squares = append(squares, &Coords{X: node.X + i, Y: node.Y - j}) }
			if i != 0 && j != 0 { squares = append(squares, &Coords{X: node.X - i, Y: node.Y - j}) }
			result[displacement] = squares
		}
	}
	return result
}

func getCheats(maze Maze, start *Coords, end *Coords, maxDisplacement int) map[string]int {
	shortestPath := maze.ShortestPath(end, start)
	fastestTime := len(shortestPath) - 1

	// Map of coords to distance from start
	pathWeights := make(map[string]int)
	for weight, coords := range shortestPath {
		pathWeights[coords.ToString()] = fastestTime - weight
	}

	// cheatMaps is a map of cheat start & cheat end to length of cheating route
	cheatsMap := make(map[string]int)
	for distanceFromEnd, cheatEnd := range shortestPath {
		// find cheat squares
		cheatSquaresByDisplacement := cheatSquares(cheatEnd, maxDisplacement)
		for displacement, squares := range cheatSquaresByDisplacement {
			for _, cheatStart := range squares {
				if maze.Contains(cheatStart) {
					timeToCheatStart := pathWeights[cheatStart.ToString()] 
					cheatTime := timeToCheatStart + displacement + distanceFromEnd
					cheatsMap[cheatStart.ToString() + cheatEnd.ToString()] = cheatTime
				}
			}
		}
	}
	return cheatsMap
}

func main() {
	dat, _ := os.ReadFile("./assets/day20-input.txt")
	maze, start, end := unpack(string(dat))

	fastestTime := len(maze.ShortestPath(end, start)) - 1
	shortCheats := getCheats(maze, start, end, 2)
	longCheats := getCheats(maze, start, end, 20)

	shortTotal := 0
	for _, cheatTime := range shortCheats {
		if saving := fastestTime - cheatTime; saving >= 100 { shortTotal += 1 }
	}

	longTotal := 0
	for _, cheatTime := range longCheats {
		if saving := fastestTime - cheatTime; saving >= 100 { longTotal += 1 }
	}

	fmt.Println(shortTotal)
	fmt.Println(longTotal)
}

