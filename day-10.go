package main

import (
	"fmt"
	"os"
	"strings"
	"strconv"
	"slices"
	"time"
)

type location struct {
	x int
	y int
}

func (l *location) Equals(other *location) bool {
	return l.x == other.x && l.y == other.y
}

func (l *location) ToString() string {
	return strconv.Itoa(l.x) + "|" + strconv.Itoa(l.y)
}

func (l *location) FromString(str string) *location {
	splitStr := strings.Split(str, "|")
	return &location{x: toInt(splitStr[0]), y: toInt(splitStr[1])}
}

func toInt(str string) int {
	res, _ := strconv.Atoi(str)
	return res
}

func unpack(input string) (topMap []string, trailheads []*location) {
	topMap = strings.Split(strings.Trim(input, "\n"), "\n")

	for yIndex, line := range topMap {
		for xIndex, character := range line {
			if toInt(string(character)) == 0 {
				trailheads = append(
					trailheads,
					&location{x: xIndex, y:yIndex},
				)
			}
		}
	}
	return topMap, trailheads 
}

func contains(topMap []string, trailNode *location) bool {
	if (trailNode.x < 0 || trailNode.y < 0) {
		return false
	}
	if (trailNode.x > len(topMap[0]) - 1 || trailNode.y > len(topMap) -1) {
		return false
	}
	return true
}

func locationValue(topMap []string, node *location) int {
	return toInt(string(topMap[node.y][node.x]))
}

func nextTrailNodes(topMap []string, trailNode *location) (nextNodes []*location) {
	// grab 4 location nodes
	adjacentNodes := []*location{
		&location{x: trailNode.x - 1, y: trailNode.y},
		&location{x: trailNode.x + 1, y: trailNode.y},
		&location{x: trailNode.x, y: trailNode.y - 1},
		&location{x: trailNode.x, y: trailNode.y + 1},
	}

	// for each node, check it's in the map and is slightly higher
	for _, node := range adjacentNodes {
		if !contains(topMap, node) {
			continue
		}
		trailNodeValue := locationValue(topMap, trailNode)
		adjacentNodeValue := locationValue(topMap, node)
		if adjacentNodeValue == trailNodeValue + 1 {
			nextNodes = append(nextNodes, node)
		}
	}
	return nextNodes
}

func findTrailPeaks(topMap []string, trailhead *location) (trailPeaks []*location, trailheadRating int) {
	peakMap := make(map[string]struct{})

	nextNodes := nextTrailNodes(topMap, trailhead)
	for len(nextNodes) > 0 {
		newNextNodes := []*location{}
		for _, nextNode := range nextNodes {
			nextValue := locationValue(topMap, nextNode)
			// if any nodes are 9, add to trailPeaks
			if nextValue == 9 {
				trailheadRating += 1
				nodeKey := nextNode.ToString()
				_, peakIsRecorded := peakMap[nodeKey]
				if !peakIsRecorded {
					peakMap[nodeKey] = struct{}{}
					trailPeaks = append(trailPeaks, nextNode)
				}
				continue
			}
			
			// else get nextNodes and concat
			newNextNodes = slices.Concat(
				newNextNodes,
				nextTrailNodes(topMap, nextNode),
			)
		}
		nextNodes = newNextNodes
	}
	return trailPeaks, trailheadRating
}

func main() {
	dat, _ := os.ReadFile("./assets/day10-input.txt")
	topMap, trailheads := unpack(string(dat))

	startTime := time.Now()

	trailScoreSum := 0
	trailheadRatingSum := 0
	for _, trailhead := range trailheads {
		trailPeaks, trailheadRating := findTrailPeaks(topMap, trailhead)
		trailScoreSum += len(trailPeaks)
		trailheadRatingSum += trailheadRating
	}
	fmt.Println(trailScoreSum)
	fmt.Println(trailheadRatingSum)
	fmt.Println(time.Since(startTime))
}

