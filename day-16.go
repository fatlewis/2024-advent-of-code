package main

import (
	"fmt"
	"os"
	"strings"
	"slices"
	"math"
	"container/heap"
	. "advent-of-code-2024/utils"
	priorityQueue "advent-of-code-2024/utils/pq"
	"advent-of-code-2024/utils/set"
)

type MazeRoute struct {
	previousPosition *Coords
	previousDirection Direction
	currentPosition *Coords
	direction Direction
	score int
}

func (r *MazeRoute) Direction() Direction { return r.direction }
func (r *MazeRoute) Score() int { return r.score }
func (r *MazeRoute) Extend(direction Direction) *MazeRoute {
	newPosition := r.currentPosition.Translate(direction)
	newScore := r.score + 1
	if r.direction != direction { newScore += 1000 }

	return &MazeRoute{
		previousPosition: r.currentPosition,
		previousDirection: r.direction,
		currentPosition: newPosition,
		direction: direction,
		score: newScore,
	}
}

// Keeps track of the nodes in optimal paths to this location
type LocationPath struct {
	score int
	pathNodes []string
}

type LocationPaths map[string]*LocationPath

func (l LocationPaths) AddFromRoute(route *MazeRoute) {
	positionString := route.currentPosition.ToString()
	key := positionString + route.direction.ToString()
	pathNodes := []string{positionString}
	if route.previousPosition != nil {
		prevKey := route.previousPosition.ToString() + route.previousDirection.ToString()
		pathNodes = append(slices.Clone(l[prevKey].pathNodes), positionString)
	}

	l[key] = &LocationPath{score: route.score, pathNodes: pathNodes}
}

func (l LocationPaths) Contains(location *Coords, direction Direction) bool {
	_, ok := l[location.ToString() + direction.ToString()]
	return ok
}

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

func main() {
	dat, _ := os.ReadFile("./assets/day16-input.txt")
	maze, start, end := unpack(string(dat))

	// Algorithm
	pq := make(priorityQueue.PriorityQueue, 0)
	heap.Init(&pq)
	heap.Push(&pq, priorityQueue.NewPQItem(0, &MazeRoute{currentPosition: start, direction: East, score: 0}))

	optimalScore := math.MaxInt32
	visited := make(LocationPaths)
	for pq.Len() > 0 {
		nextRoute := heap.Pop(&pq).(*priorityQueue.PQItem).GetValue().(*MazeRoute)

		previousKey := ""
		if nextRoute.previousPosition != nil {
			previousKey = nextRoute.previousPosition.ToString() + nextRoute.previousDirection.ToString()
		}
		key := nextRoute.currentPosition.ToString() + nextRoute.direction.ToString()

		// If we're outside the maze or in a wall, skip
		if !maze.Contains(nextRoute.currentPosition) { continue }

		// If we've gone beyond the optimal score, skip
		if nextRoute.score > optimalScore { continue }

		// If we've visited the current location before:
		if locationPath := visited[key]; locationPath != nil {
			// If score is equal, merge our previous path into the existing one
			if locationPath.score == nextRoute.score {
				locationPath.pathNodes = set.Union(
					locationPath.pathNodes,
					visited[previousKey].pathNodes,
				)
			}
			// Next routes have aleady been added, skip
			continue
		}
		// Add current path to visited
		visited.AddFromRoute(nextRoute)

		// If we're at the end, update optimalScore
		if nextRoute.currentPosition.Equals(end) {
			optimalScore = nextRoute.score
			continue
		}

		// Generate next routes and add to pq
		direction := nextRoute.direction
		routeAhead := nextRoute.Extend(direction)
		heap.Push(&pq, priorityQueue.NewPQItem(routeAhead.score * -1, routeAhead))

		routeRight := nextRoute.Extend(direction.TurnClockwise())
		heap.Push(&pq, priorityQueue.NewPQItem(routeRight.score * -1, routeRight))

		routeLeft := nextRoute.Extend(direction.TurnCounterClockwise())
		heap.Push(&pq, priorityQueue.NewPQItem(routeLeft.score * -1, routeLeft))
	}

	endKey := end.ToString()
	endNodeKeys := []string{
		endKey + North.ToString(),
		endKey + East.ToString(),
		endKey + South.ToString(),
		endKey + West.ToString(),
	}
	combinedEndPath := &LocationPath{score: math.MaxInt32, pathNodes: []string{}}
	for _, endNodeKey := range endNodeKeys {
		endNode, ok := visited[endNodeKey]
		if ok {
			combinedEndPath.score = endNode.score
			combinedEndPath.pathNodes = set.Union(combinedEndPath.pathNodes, endNode.pathNodes)
		}
	}
	fmt.Println(combinedEndPath.score, len(combinedEndPath.pathNodes))

	fmt.Println(len(maze.ShortestPath(start, end)))
}

