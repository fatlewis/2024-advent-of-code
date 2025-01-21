package utils

import (
	"fmt"
	"container/heap"
	priorityQueue "advent-of-code-2024/utils/pq"
)

type Direction int

const (
	North Direction = iota
	East Direction = iota
	South Direction = iota
	West Direction = iota
)

func (d Direction) TurnClockwise() Direction { return (d + 1) % 4 }
func (d Direction) TurnCounterClockwise() Direction { return (d + 3 ) % 4 }
func (d Direction) ToString() string {
	switch d {
	case North:
		return "N"
	case East:
		return "E"
	case South:
		return "S"
	case West:
		return "W"
	}
	return ""
}

type Coords struct {
	X int
	Y int
}

func (coords *Coords) Equals(location *Coords) bool {
	return coords.X == location.X && coords.Y == location.Y
}

func (coords *Coords) Translate(direction Direction) *Coords {
	newCoords := &Coords{X: coords.X, Y: coords.Y}

	if direction == North { newCoords.Y -= 1 }
	if direction == East { newCoords.X += 1 }
	if direction == South { newCoords.Y += 1 }
	if direction == West { newCoords.X -= 1 }

	return newCoords
}

func (coords *Coords) ToString() string {
	return fmt.Sprintf("(%d,%d)", coords.X, coords.Y)
}

type Maze map[int]map[int]*Coords

func (m Maze) Contains(location *Coords) bool {
	row, yOk := m[location.Y]
	if !yOk { return false }

	_, xOk := row[location.X]
	if !xOk { return false }

	return true
}

func (m Maze) Add(location *Coords) {
	row, yOk := m[location.Y]
	if !yOk {
		row = make(map[int]*Coords)
		m[location.Y] = row
	}
	row[location.X] = location
}

func (m Maze) Remove(location *Coords) {
	row, yOk := m[location.Y]
	if !yOk { return }
	delete(row, location.X)
}

type Path struct {
	score int
	position *Coords
	nodes []*Coords
}

func (m Maze) ShortestPath(start, end *Coords) []*Coords {
        pq := make(priorityQueue.PriorityQueue, 0)
        heap.Init(&pq)
	heap.Push(&pq, priorityQueue.NewPQItem(0, &Path{position: start, nodes: []*Coords{start}}))

        visited := make(map[string]bool) 
        for pq.Len() > 0 {
                route := heap.Pop(&pq).(*priorityQueue.PQItem).GetValue().(*Path)
		p := route.position
		key := p.ToString()

                // If we're outside the maze or in a wall, skip
                if !m.Contains(p) { continue }

                // If we've visited the current location before, skip
		if visited[key] { continue }

                // Add current node to visited
                visited[key] = true

                // If we're at the end, return shortest path
                if p.Equals(end) { return route.nodes }

                // Generate next routes and add to pq
		newRoutes := []*Path{
			&Path{position: p.Translate(North), nodes: append(route.nodes, p.Translate(North))},
			&Path{position: p.Translate(East), nodes: append(route.nodes, p.Translate(East))},
			&Path{position: p.Translate(South), nodes: append(route.nodes, p.Translate(South))},
			&Path{position: p.Translate(West), nodes: append(route.nodes, p.Translate(West))},
		}
		for _, newRoute := range newRoutes {
			heap.Push(&pq, priorityQueue.NewPQItem(len(newRoute.nodes) * -1, newRoute))
		}
        }
	return []*Coords{}
}

