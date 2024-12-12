package main

import (
	"fmt"
	"os"
	"time"
	"strings"
	"strconv"
	"maps"
	"slices"
)

type location struct {
	x int
	y int
}

func (l *location) ToString() string {
	return strconv.Itoa(l.x) + "|" + strconv.Itoa(l.y)
}

func (l *location) Translate(vector [2]int) *location {
	return &location{x: l.x+vector[0], y: l.y+vector[1]}
}

type region struct {
	symbol rune
	locations []*location
	locationMap map[string]struct{}
}

func NewRegion(symbol rune) *region {
	return &region{
		symbol: symbol,
		locations: []*location{},
		locationMap: make(map[string]struct{}),
	}
}

func (r *region) AddLocation(l *location) {
	_, isSet := r.locationMap[l.ToString()]
	if !isSet {
		r.locationMap[l.ToString()] = struct{}{}
		r.locations = append(r.locations, l)
	}
}

func (r *region) AddLocations(locations []*location) {
	for _, l := range locations {
		r.AddLocation(l)
	}
}

func (r *region) NumAdjacentLocations(l *location) int {
	number := 0
	adjacentLocations := []*location{
		l.Translate([2]int{-1, 0}),
		l.Translate([2]int{1, 0}),
		l.Translate([2]int{0, -1}),
		l.Translate([2]int{0, 1}),
	}
	for _, adjacentLocation := range adjacentLocations {
		_, isSet := r.locationMap[adjacentLocation.ToString()]
		if isSet {
			number += 1
		}
	}
	return number
}

func (r *region) GetArea() int {
	return len(r.locations)
}

func (r *region) GetPerimeter() int {
	perimeter := 0
	for _, location := range r.locations {
		perimeter += 4 - r.NumAdjacentLocations(location)
	}
	return perimeter
}

func rotate90(vector [2]int) [2]int {
	// rotates a movement vector 90 degrees
	return [2]int{vector[1]*-1, vector[0]}
}

func (r *region) GetSides() int {
	number := 0
	testVectors := [][2]int{
		[2]int{-1, 0},
		[2]int{1, 0},
		[2]int{0, -1},
		[2]int{0, 1},
	}
	for _, l := range r.locations {
		for _, testVector := range testVectors {
			// If side continues in this direction, it isn't new
			adjacentLocation := l.Translate(testVector)
			_, adjacentLocationExists := r.locationMap[adjacentLocation.ToString()]
			if adjacentLocationExists {
				continue
			}
			// _| shape corner
			// if no square in direction of rotate90(vector), it's a corner
			_, locationAt90DegreesExists := r.locationMap[l.Translate(rotate90(testVector)).ToString()]
			if !locationAt90DegreesExists {
				number += 1
				continue
			}
			// |_ shape corner
			// if square exists in direction of rotate90(vector) from the adjacent square, it's a corner
			_, locationAt90DegreesFromAdjacentLocationExists := r.locationMap[adjacentLocation.Translate(rotate90(testVector)).ToString()]
			if locationAt90DegreesFromAdjacentLocationExists {
				number += 1
			}
		}
	}
	return number
}

func main() {
	dat, _ := os.ReadFile("./assets/day12-input.txt")
	lines := strings.Split(strings.Trim(string(dat), "\n"), "\n")

	start := time.Now()

	regionsMap := make(map[rune][]*region)
	for lineIndex, line := range lines {
		for charIndex, character := range line {
			charLocation := &location{x: charIndex, y: lineIndex}
			testRegions := regionsMap[character]
			var locationRegion *region
			for regionIndex, testRegion := range testRegions {
				if testRegion.NumAdjacentLocations(charLocation) > 0 {
					// if it fits in multiple regions, join the regions together
					if (locationRegion != nil) {
						// we've already attached the location to a region, so merge the regions
						locationRegion.AddLocations(testRegion.locations)
						regionsMap[character] = slices.Delete(regionsMap[character], regionIndex, regionIndex+1)
						break
					}

					testRegion.AddLocation(charLocation)
					locationRegion = testRegion
				}
			}
			if (locationRegion == nil) {
				newRegion := NewRegion(character)
				newRegion.AddLocation(charLocation)
				regionsMap[character] = append(regionsMap[character], newRegion)
			}
		}
	}

	fmt.Println("Setup time:", time.Since(start))

	start = time.Now()

	totalPrice := 0
	for key := range maps.Keys(regionsMap) {
		for _, mapRegion := range regionsMap[key] {
			regionPrice := mapRegion.GetArea() * mapRegion.GetPerimeter()
			totalPrice += regionPrice
		}
	}
	fmt.Println(totalPrice)
	fmt.Println("Part 1 time:", time.Since(start))

	totalPrice = 0
	start = time.Now()
	for key := range maps.Keys(regionsMap) {
		for _, mapRegion := range regionsMap[key] {
			regionPrice := mapRegion.GetArea() * mapRegion.GetSides()
			totalPrice += regionPrice
		}
	}
	fmt.Println(totalPrice)
	fmt.Println("Part 2 time:", time.Since(start))
}

