package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

func intersect[T comparable](a []T, b []T) (result []T) {
	for _, aElem := range a {
		if slices.Contains(b, aElem) {
			result = append(result, aElem)
		}
	}
	return result
}

func setDifference[T comparable](a []T, b[]T) (result []T) {
	for _, aElem := range a {
		if !slices.Contains(b, aElem) {
			result = append(result, aElem)
		}
	}
	return result
}

type Direction int
const NORTH = 0
const SOUTH = 1
const EAST = 2
const WEST = 3

func (d Direction) ToString() string {
	switch d {
	case 0:
		return "NORTH"
	case 1:
		return "SOUTH"
	case 2:
		return "EAST"
	case 3:
		return "WEST"
	}
	return "NOT A DIRECTION"
}

type Location struct {
	x int
	y int
}

func (l *Location) Transform(vector *Location) *Location {
	return &Location{x: l.x + vector.x, y: l.y + vector.y}
}

type MazeNode struct {
	// Means by which the node has been reached from each direction
	routes map[Direction]*Route
	location *Location
}

type Route struct {
	length int
	turns int
	direction Direction
	nodes []*MazeNode
}

// Creates a new route by extending into an adjacent node
func (r *Route) Extend(node *MazeNode) *Route {
	lastNode := r.nodes[len(r.nodes)-1]
	newDirection := lastNode.Direction(node)
	newTurns := r.turns
	if newDirection != r.direction {
		newTurns += 1
	}
	return &Route{
		length: r.length + 1,
		turns: newTurns,
		direction: newDirection,
		nodes: append(slices.Clone(r.nodes), node),
	}
}

func (r *Route) Score() int {
	return r.length + (r.turns * 1000)
}

// Finds the score required to continue this route in any unvisited direction in the supplied maze
func (r *Route) MinDepartureScore(maze Maze) int {
	lastNode := r.nodes[len(r.nodes)-1]
	adjacentNodes := lastNode.AdjacentNodes(maze)

	// If there are no nodes, return score
	if len(adjacentNodes) == 0 {
		return r.Score()
	}

	// If there is a node directly ahead, return score + 1
	hasNodeAhead := slices.ContainsFunc(adjacentNodes, func (n *MazeNode) bool {
		return lastNode.Direction(n) == r.direction
	})
	if hasNodeAhead {
		return r.Score() + 1
	}

	// Else return score + 1001
	return r.Score() + 1001
}

type Maze map[int]map[int]*MazeNode

var maze string

func (n *MazeNode) ShortestRoutes() []*Route {
	// Cycles through the routes on the node and returns all with the smallest length
	shortestRoutes := []*Route{}
	for _, route := range n.routes {
		// If no routes have ben added yet, or this route has equal shortest length, append
		if (len(shortestRoutes) == 0 || route.length == shortestRoutes[0].length) {
			shortestRoutes = append(shortestRoutes, route)
			continue
		}

		// If this route is shorter, replace existing routes with this one
		if (route.length < shortestRoutes[0].length) {
			shortestRoutes = []*Route{route}
			continue
		}
	}
	return shortestRoutes
}

// Finds all adjacent nodes which exist on the mazeMap
func (n *MazeNode) AdjacentNodes(mazeMap Maze) []*MazeNode {
	adjacentLocations := []*Location{
		n.location.Transform(&Location{x: 1, y: 0}),
		n.location.Transform(&Location{x: -1, y: 0}),
		n.location.Transform(&Location{x: 0, y: 1}),
		n.location.Transform(&Location{x: 0, y: -1}),
	}

	adjacentNodes := []*MazeNode{}
	for _, adjacentLocation := range adjacentLocations {
		mazeLine, foundLine := mazeMap[adjacentLocation.y]
		if !foundLine {
			continue
		}
		adjacentNode, foundNode := mazeLine[adjacentLocation.x]
		if (foundNode) {
			adjacentNodes = append(adjacentNodes, adjacentNode)
		}
	}

	return adjacentNodes
}

// Find the direction to an adjacent node
func (n *MazeNode) Direction(to *MazeNode) Direction {
	if (n.location.x > to.location.x) {
		return WEST
	}
	if (n.location.x < to.location.x) {
		return EAST
	}
	if (n.location.y > to.location.y) {
		return NORTH
	}
	return SOUTH
}

func (n *MazeNode) MakeRoute(to *MazeNode) *Route {
	shortestRoutes := n.ShortestRoutes()

	optimalRoute := shortestRoutes[0].Extend(to)
	for _, shortestRoute := range shortestRoutes[1:] {
		extendedRoute := shortestRoute.Extend(to)
		if extendedRoute.Score() < optimalRoute.Score() {
			optimalRoute = extendedRoute
		}
		if extendedRoute.Score() == optimalRoute.Score() {
			extraNodes := setDifference(extendedRoute.nodes, optimalRoute.nodes)
			optimalRoute.nodes = slices.Concat(extraNodes, optimalRoute.nodes)
		}
	}
	return optimalRoute
}

func unpack(input string) (mazeMap Maze, unvisited []*MazeNode, start *MazeNode, end *MazeNode) {
        mazeLines := strings.Split(strings.Trim(input, "\n"), "\n")
        mazeMap = make(Maze)
	unvisited = []*MazeNode{}
        for yIndex, mazeLine := range mazeLines {
                mazeMap[yIndex] = make(map[int]*MazeNode)
                for xIndex, mazeChar := range mazeLine {
                        if (mazeChar != '#') {
				mapNode := &MazeNode{
					location: &Location{x: xIndex, y: yIndex},
					routes: make(map[Direction]*Route),
				}
                                mazeMap[yIndex][xIndex] = mapNode
                                if (mazeChar == 'S') {
                                        // Set empty route on start with starting direction
					mapNode.routes[EAST] = &Route{nodes: []*MazeNode{mapNode}, direction: EAST}
                                        start = mapNode
					unvisited = slices.Insert(unvisited, 0, start)
                                } else {
					unvisited = append(unvisited, mapNode)
				}
                                if (mazeChar == 'E') {
                                        end = mapNode
                                }
                        }
                }
        }
        return mazeMap, unvisited, start, end
}

func (route *Route) Print(inputMaze string) {
        mazeLines := strings.Split(strings.Trim(inputMaze, "\n"), "\n")
        mazeSlice := [][]string{}
        for _, mazeLine := range mazeLines {
                mazeSlice = append(mazeSlice, strings.Split(mazeLine, ""))
        }
        for _, routeNode := range route.nodes {
                mazeSlice[routeNode.location.y][routeNode.location.x] = "O"
        }
        for _, mazeSliceLine := range mazeSlice {
                fmt.Println(mazeSliceLine)
        }
        fmt.Println()
}

func minDepartureScore(routes []*Route, maze Maze) (int, *Route) {
	resultRoute := &Route{}
	result := 0
	for _, route := range routes {
		routeScore := route.MinDepartureScore(maze)
		if result == 0 || routeScore < result {
			resultRoute = route
			result = routeScore
		}
	}
	return result, resultRoute
}

func main() {
	// 442 is TOO LOW
	// Correct is 467, but this code does not output that result, and skips some nodes. Why???

	// Gather nodes in map & list
	dat, _ := os.ReadFile("./assets/day16-input.txt")
	mazeMap, unvisited, _, end := unpack(string(dat))
	maze = string(dat) 

	for len(unvisited) > 0 {
		/*for _, unvisitedNode := range unvisited {
			fmt.Print(unvisitedNode.location, " || ", len(unvisitedNode.routes), ", ")
		}
		fmt.Println("\n")*/

		// Take first node
		currentNode := unvisited[0]
		unvisited = unvisited[1:]

		// If this node is the end, break
		if currentNode == end {
			break
		}

		// Get unvisited neighbours
		unvisitedNeighbours := intersect(
			currentNode.AdjacentNodes(mazeMap),
			unvisited,
		)

		// For each neighbour
		for _, unvisitedNeighbour := range unvisitedNeighbours {
			// Find direction to neighbour
			routeDirection := currentNode.Direction(unvisitedNeighbour)
			/*fmt.Printf(
				"Direction from (%d,%d) to (%d,%d) is %s.\n",
				currentNode.location.x, currentNode.location.y,
				unvisitedNeighbour.location.x, unvisitedNeighbour.location.y,
				routeDirection.ToString(),
			)*/

			// Calculate route
			routeToNeighbour := currentNode.MakeRoute(unvisitedNeighbour)
			//routeToNeighbour.Print(string(dat))
			
			// Add to neighbour's routes
			unvisitedNeighbour.routes[routeDirection] = routeToNeighbour
		}

		// Sort unvisited by score 
		slices.SortFunc(unvisited, func (a, b *MazeNode) int {
			aRoutes := a.ShortestRoutes()
			bRoutes := b.ShortestRoutes()
			// If 1 or both of the nodes are empty
			if len(aRoutes) == 0 {
				if len(bRoutes) == 0 {
					return 0
				} else {
					return 1
				}
			} else if (len(bRoutes) == 0) {
				return -1
			}
			minADepartureScore, minADepartureRoute := minDepartureScore(aRoutes, mazeMap)
			minBDepartureScore, minBDepartureRoute := minDepartureScore(bRoutes, mazeMap)
			departureDifference := minADepartureScore - minBDepartureScore
			if (departureDifference == 0) {
				return minBDepartureRoute.turns - minADepartureRoute.turns
			}
			return departureDifference
		})
	}
	endRoutes := end.ShortestRoutes()
	optimalRoute := endRoutes[0]
	optimalRoute.Print(string(dat))
	for _, endRoute := range endRoutes[1:] {
		endRoute.Print(string(dat))
		// If another route is strictly better, replace it
		if endRoute.turns < optimalRoute.turns {
			optimalRoute = endRoute
		}
		// If another route is equal, append any missing nodes
		if endRoute.turns == optimalRoute.turns {
			extraNodes := setDifference(endRoute.nodes, optimalRoute.nodes)
			optimalRoute.nodes = slices.Concat(extraNodes, optimalRoute.nodes)
		}
	}
	fmt.Println(optimalRoute.length + (optimalRoute.turns * 1000))
	fmt.Println(len(optimalRoute.nodes))
}

