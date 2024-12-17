package main

import (
	"fmt"
	"os"
	"strings"
	"maps"
	"slices"
)

const BOX = 0
const WALL = 1
const ROBOT = 2

const BOX_LEFT = 3
const BOX_RIGHT = 4

type location struct {
	x int
	y int
}

func (l *location) ToString() string {
	return fmt.Sprintf("(%d,%d)", l.x, l.y)
}

func (l *location) Transform(vector *location) *location {
	return &location{x: l.x + vector.x, y: l.y + vector.y}
}

func (l *location) Invert() *location {
	return &location{l.x * -1, l.y * -1}
}

type warehouse map[string]*object

func toRune(objectType int) rune {
	switch objectType {
	case BOX:
		return 'O'
	case WALL:
		return '#'
	case ROBOT:
		return '@'
	case BOX_LEFT:
		return '['
	case BOX_RIGHT:
		return ']'
	}
	return 'X'
}

func printWarehouse(w warehouse, width int, height int) {
	warehouseStr := ""
	for yIndex := range height {
		for xIndex := range width {
			warehouseLocation := &location{x: xIndex, y: yIndex}
			obj, found := w[warehouseLocation.ToString()]
			if (found) {
				warehouseStr = warehouseStr + string(toRune(obj.objectType))
			} else {
				warehouseStr = warehouseStr + "."
			}
		}
		warehouseStr = warehouseStr + "\n"
	}
	fmt.Println(warehouseStr)
}

type object struct {
	location *location
	objectType int // BOX || WALL || ROBOT
}

func moveObject(w warehouse, from *location, to *location) {
		// Move object
		objectToMove := w[from.ToString()]
		w[to.ToString()] = objectToMove

		// Update location on object
		objectToMove.location = &location{x: to.x, y: to.y}

		// Remove old location from warehouse
		delete(w, from.ToString())
}

func move(robot *object, w warehouse, vector *location) {
	// Search squares in the vector direction until a wall or empty square is reached
	currentLocation := robot.location.Transform(vector)
	searchObject, objectFound := w[currentLocation.ToString()]
	for objectFound && searchObject.objectType != WALL {
		currentLocation = currentLocation.Transform(vector)
		searchObject, objectFound = w[currentLocation.ToString()]
	}

	// If WALL, do nothing
	if (objectFound && searchObject.objectType == WALL) {
		return
	}

	// Reverse vector and shift objects until robot square is reached.
	reverseVector := vector.Invert()
	robotLocation := robot.location
	shiftLocation := currentLocation.Transform(reverseVector)
	for currentLocation.x != robotLocation.x || currentLocation.y != robotLocation.y {
		// Move object
		moveObject(w, shiftLocation, currentLocation)

		// Update locations
		currentLocation = shiftLocation
		shiftLocation = currentLocation.Transform(reverseVector)
	}
	// Empty final location
	delete(w, currentLocation.ToString())
}

func fillWarehouse(w warehouse, items []string, expand bool) *object {
	// For each warehouse character, add to map
	var robot *object
	for yValue, warehouseLine := range items {
		for xValue, character := range warehouseLine {
			if (expand) {
				xValue *= 2
			}

			item := &object{location: &location{x: xValue, y: yValue}}
			switch character {
			case '@':
				item.objectType = ROBOT
				robot = item
			case 'O':
				item.objectType = BOX
			case '#':
				item.objectType = WALL
			default:
				continue
			}

			if (expand && item.objectType != ROBOT) {
				// Create duplicate
				duplicateItem := &object{
					location: item.location.Transform(&location{1, 0}),
					objectType: item.objectType,
				}
				// Update object types if it's a box
				if (item.objectType == BOX) {
					item.objectType = BOX_LEFT
					duplicateItem.objectType = BOX_RIGHT
				}
				// Add duplicate
				w[duplicateItem.location.ToString()] = duplicateItem
			}

			// Add item
			w[item.location.ToString()] = item
		}
	}
	return robot
}

func unpack(input string, expand bool) (warehouse, *object, []*location, *location) {
	// Split string to warehouse and moves
	parts := strings.Split(strings.Trim(input, "\n"), "\n\n")

	// Generate warehouse
	w := make(warehouse)
	warehouseLines := strings.Split(parts[0], "\n")
	mapSize := &location{x: len(warehouseLines[0]), y: len(warehouseLines)}
	if (expand) {
		mapSize.x *= 2
	}
	robot := fillWarehouse(w, warehouseLines, expand)

	// For each move, transform to vector and add to slice
	moves := []*location{}
	for _, vectorValue := range parts[1] {
		switch vectorValue {
		case '<':
			moves = append(moves, &location{-1, 0})
		case 'v':
			moves = append(moves, &location{0, 1})
		case '>':
			moves = append(moves, &location{1, 0})
		case '^':
			moves = append(moves, &location{0, -1})
		}
	}
	
	return w, robot, moves, mapSize
}

func movePart2(robot *object, w warehouse, vector *location) {
	// Same as part1 for x-axis
	if (vector.x != 0) {
		move(robot, w, vector)
		return
	}

	currentLineObjs := []*object{robot}
	objectsToMove := [][]*object{}
	nextLineObjs := []*object{}
	// While the current line has objects in it
	for len(currentLineObjs) > 0 {
		// Add currentLineObjs to objectsToMove
		objectsToMove = append(objectsToMove, currentLineObjs)
		// Gather nextLineObjs
		for _, currentLineObj := range currentLineObjs {
			dropLocation := currentLineObj.location.Transform(vector)
			dropObj, dropObjExists := w[dropLocation.ToString()]

			// If square is empty, continue
			if (!dropObjExists) {
				continue
			}

			// If square is a wall, stop, no movement
			if (dropObj.objectType == WALL) {
				return
			}

			// If square is BOX_LEFT and it's not recorded already, add the box
			if (dropObj.objectType == BOX_LEFT && !slices.Contains(nextLineObjs, dropObj)) {
				boxLeft := dropObj
				boxRight := w[dropLocation.Transform(&location{x: 1, y: 0}).ToString()]
				nextLineObjs = append(nextLineObjs, boxLeft, boxRight)
			}

			// If square is BOX_RIGHT and it's not recorded already, add the box
			if (dropObj.objectType == BOX_RIGHT && !slices.Contains(nextLineObjs, dropObj)) {
				boxRight := dropObj
				boxLeft := w[dropLocation.Transform(&location{x: -1, y: 0}).ToString()]
				nextLineObjs = append(nextLineObjs, boxLeft, boxRight)
			}
		}
		// Set currentLineObjs = nextLineObjs, reset nextLineObjs
		currentLineObjs = nextLineObjs
		nextLineObjs = []*object{}
	}
	
	// Move objects, from last line to first line
	slices.Reverse(objectsToMove)
	for _, objLine := range objectsToMove {
		for _, obj := range objLine {
			// moveObject to square in direction of movement
			moveObject(w, obj.location, obj.location.Transform(vector))
		}
	}
}

func main() {
	dat, _ := os.ReadFile("./assets/day15-input.txt")
	warehouseMap, robot, moves, mapSize := unpack(string(dat), false)

	for _, robotMove := range moves {
		move(robot, warehouseMap, robotMove)
	}

	// This is horribly inefficient but if it works... ¯\\_(ツ)_/¯
	gpsSum := 0
	for warehouseKey := range maps.Keys(warehouseMap) {
		warehouseObject := warehouseMap[warehouseKey]
		if (warehouseObject.objectType == BOX) {
			gpsSum += 100*warehouseObject.location.y + warehouseObject.location.x
		}
	}
	fmt.Println(gpsSum)

	warehouseMap, robot, moves, mapSize = unpack(string(dat), true)

	for _, robotMove := range moves {
		movePart2(robot, warehouseMap, robotMove)
	}

	// Calculate in a new way
	gpsSum = 0
	for warehouseKey := range maps.Keys(warehouseMap) {
		warehouseObject := warehouseMap[warehouseKey]
		if (warehouseObject.objectType == BOX_LEFT) {
			gpsSum += 100*warehouseObject.location.y + warehouseObject.location.x
		}
	}
	fmt.Println(gpsSum)
}
