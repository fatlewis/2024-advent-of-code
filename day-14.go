package main

import (
	"fmt"
	"os"
	"strconv"
	"regexp"
	"slices"
)

const SPACE_HEIGHT = 103 // 7
const SPACE_WIDTH = 101 // 11

type robot struct {
	position [2]int
	velocity [2]int
}

func (r *robot) ToString() string {
	return fmt.Sprintf("pos(%d,%d) vel(%d,%d)", r.position[0], r.position[1], r.velocity[0], r.velocity[1])
}

func (r *robot) Move(seconds int) {
	newX := (r.position[0] + (r.velocity[0] * seconds)) % SPACE_WIDTH
	if (newX < 0) {
		newX = SPACE_WIDTH + newX
	}
	newY := (r.position[1] + (r.velocity[1] * seconds)) % SPACE_HEIGHT
	if (newY < 0) {
		newY = SPACE_HEIGHT + newY
	}
	r.position = [2]int{newX, newY}
}

func toInt(str string) int {
	res, _ := strconv.Atoi(str)
	return res
}

func unpack(input string) []*robot {
	robots := []*robot{}
	re := regexp.MustCompile(`p=(-?\d+),(-?\d+).*v=(-?\d+),(-?\d+)`)
	for _, values := range re.FindAllStringSubmatch(input, -1) {
		robots = append(robots, &robot{
			position: [2]int{toInt(values[1]), toInt(values[2])},
			velocity: [2]int{toInt(values[3]), toInt(values[4])},
		})
	}
	return robots
}

func showField(robots []*robot) {
	field := [][]rune{}
	for i := 0; i <= SPACE_HEIGHT; i++ {
		field = append(field, slices.Repeat([]rune{'.'}, SPACE_WIDTH))
	}
	for _, r := range robots {
		currentRune := field[r.position[1]][r.position[0]]
		if (currentRune == '.') {
			field[r.position[1]][r.position[0]] = '1'
			continue
		}
		// Going to assume that we never have more than 9 robots on a square...
		field[r.position[1]][r.position[0]] = rune(strconv.Itoa(toInt(string(currentRune)) + 1)[0])
	}
	fieldString := ""
	for _, line := range field {
		for _, character := range line {
			fieldString = fieldString + string(character)
		}
		fieldString = fieldString + "\n"
	}
	fmt.Println(fieldString)
}

func addRobotToQuadrants(quadrants [4]int, robotInstance *robot) [4]int {
	newQuadrants := [4]int{quadrants[0], quadrants[1], quadrants[2], quadrants[3]}
	midY := SPACE_HEIGHT / 2
	midX := SPACE_WIDTH / 2
	if (robotInstance.position[0] < midX) {
		if (robotInstance.position[1] < midY) {
			newQuadrants[0] += 1
		} else if (robotInstance.position[1] > midY) {
			newQuadrants[3] += 1
		}
	} else if robotInstance.position[0] > midX {
		if (robotInstance.position[1] < midY) {
			newQuadrants[1] += 1
		} else if (robotInstance.position[1] > midY) {
			newQuadrants[2] += 1
		}
	}
	return newQuadrants
}

func main() {
	dat, _ := os.ReadFile("./assets/day14-input.txt")
	robots := unpack(string(dat))

	// TL, TR, BR, BL
	quadrants := [4]int{0, 0, 0, 0}
	for _, robotInstance := range robots {
		// Move robot for 100 seconds
		robotInstance.Move(100)
		// Add value to relevant quadrant
		quadrants = addRobotToQuadrants(quadrants, robotInstance)
	}
	safetyFactor := quadrants[0] * quadrants[1] * quadrants[2] * quadrants[3]
	fmt.Println(safetyFactor)

	// Test each possible frame
	TEST_ITERATIONS := 103 * 101 
	minSafetyFactor := 0
	minIterationNum := 0
	robots = unpack(string(dat))
	for i := 1; i <= TEST_ITERATIONS; i++ {
		quadrants := [4]int{0, 0, 0, 0}
		for _, robotInstance := range robots {
			robotInstance.Move(1)
			quadrants = addRobotToQuadrants(quadrants, robotInstance)
		}

		// Guess that, if the easter egg is showing, the robots will be grouped together and safety is likely low
		safetyFactor := quadrants[0] * quadrants[1] * quadrants[2] * quadrants[3]
		if (minSafetyFactor == 0 || safetyFactor < minSafetyFactor) {
			minSafetyFactor = safetyFactor
			minIterationNum = i
		}

		// Found that safety min was 7286 on a previous run, check what this frame shows
		if (i == 7286) {
			showField(robots)
		}
	}
	fmt.Println(minIterationNum)
}

