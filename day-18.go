package main

import (
	"fmt"
	"os"
	"strings"
	"strconv"
	. "advent-of-code-2024/utils"
)

func makeMaze(size int) Maze {
	maze := make(Maze)
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			maze.Add(&Coords{X: j, Y: i})
		}
	}
	return maze
}

func printMaze(m Maze, size int) {
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if m.Contains(&Coords{X: j, Y: i}) {
				fmt.Print(".")
			} else {
				fmt.Print("#")
			}
		}
		fmt.Print("\n")
	}
	fmt.Print("\n")
}

func unpack(input string) []*Coords {
	lines := strings.Split(strings.Trim(input, "\n"), "\n")
	bytes := []*Coords{}
	for _, line := range lines {
		positions := strings.Split(line, ",")
		xPos, _ := strconv.Atoi(positions[0])
		yPos, _ := strconv.Atoi(positions[1])
		bytes = append(bytes, &Coords{X: xPos, Y: yPos}) 
	}
	return bytes
}

func main() {
	dat, _ := os.ReadFile("./assets/day18-input.txt")
	bytes := unpack(string(dat))

	// Part 1
	// 71x71 maze
	maze := makeMaze(71)
	// 1024 bytes
	for i := 0; i < 1024; i++ {
		maze.Remove(bytes[i])
	}
	shortestPath := maze.ShortestPath(&Coords{X: 0, Y: 0}, &Coords{X: 70, Y: 70})
	// -1 from path to count steps
	fmt.Println(len(shortestPath) - 1)

	// Part 2
	maze = makeMaze(71)
	minByte := 0
	maxByte := len(bytes) - 1
	for minByte + 1 < maxByte {
		maze = makeMaze(71)
		testByte := (minByte + maxByte) / 2
		for byteIndex, byteValue := range bytes {
			if byteIndex > testByte { break }
			maze.Remove(byteValue)
		}
		shortestPath := maze.ShortestPath(&Coords{X: 0, Y: 0}, &Coords{X: 70, Y: 70})
		if len(shortestPath) == 0 {
			maxByte = testByte
		} else {
			minByte = testByte
		}
	}
	fmt.Println(bytes[maxByte])
}

