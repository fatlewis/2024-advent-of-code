package main

import (
	"fmt"
	"os"
	"strings"
	"strconv"
	"regexp"
	"time"
)

func toInt(str string) int {
	res, _ := strconv.Atoi(str)
	return res
}

func min(a int, b int) int {
	if (a > b) {
		return b
	}
	return a
}

func max(a int, b int) int {
	if (b > a) {
		return b
	}
	return a
}

type machine struct {
	aVector [2]int
	bVector [2]int
	prizeLocation [2]int
}

type solution struct {
	aPresses int
	bPresses int
}

func (s *solution) GetPrice() int {
	return 3*s.aPresses + s.bPresses
}

func unpack(input string) []*machine {
	stringMachines := strings.Split(strings.Trim(input, "\n"), "\n\n")
	machines := []*machine{}
	for _, stringMachine := range stringMachines {
		re := regexp.MustCompile(`(?:\+|=)(\d+)`)
		values := re.FindAllStringSubmatch(stringMachine, -1)
		machines = append(machines, &machine{
			aVector: [2]int{toInt(values[0][1]), toInt(values[1][1])},
			bVector: [2]int{toInt(values[2][1]), toInt(values[3][1])},
			prizeLocation: [2]int{toInt(values[4][1]), toInt(values[5][1])},
		})
	}
	return machines
}

func findMaxQuotient(prize [2]int, button [2]int) (axis int, quotient int) {
	maxX := prize[0] / button[0]
	maxY := prize[1] / button[1]
	if (maxX < maxY) {
		return 0, maxX // min(maxX, 100)
	}
	return 1, maxY // min(maxY, 100)
}

func findMaxQuotientWithLimit(prize [2]int, button [2]int, limit int) (axis int, quotient int) {
	maxX := prize[0] / button[0]
	maxY := prize[1] / button[1]
	if (maxX < maxY) {
		return 0, min(maxX, limit) 
	}
	return 1, min(maxY, limit)
}

func findMinWinningTokensPart1(machines []*machine) int {
	minimumWinningTokens := 0
	for _, machine := range machines {
		axis, quotient := findMaxQuotientWithLimit(machine.prizeLocation, machine.aVector, 100)
		solutions := []*solution{}
		for i := 0; i <= quotient; i++ {
			remainder := machine.prizeLocation[axis] - machine.aVector[axis]*i
			// If remainder is 0, test solution with other axis using just a
			if (remainder == 0) {
				otherAxis := (axis + 1) % 2
				if (machine.prizeLocation[otherAxis] - machine.aVector[otherAxis]*i == 0) {
					solutions = append(solutions, &solution{aPresses: i, bPresses: 0})
					continue
				}
			}
			// Test if solution exists using b to make up the remainder
			if (remainder % machine.bVector[axis] == 0) {
				bQuotient := remainder / machine.bVector[axis]
				// Can't press button more than 100 times
				if (bQuotient > 100) {
					continue
				}
				// Test other axis, if true add solution
				otherAxis := (axis + 1) % 2
				otherRemainder := machine.prizeLocation[otherAxis] - machine.aVector[otherAxis]*i
				if(otherRemainder == bQuotient*machine.bVector[otherAxis]) {
					solutions = append(solutions, &solution{aPresses: i, bPresses: bQuotient})
				}
			}
		}

		// Start at 0. If there are no solutions, this ensures no credits are added to the price.
		cheapestSolution := 0
		for _, solution := range solutions {
			solutionPrice := solution.GetPrice()
			if (cheapestSolution == 0 || solutionPrice < cheapestSolution) {
				cheapestSolution = solutionPrice
			}
		}
		minimumWinningTokens += cheapestSolution
	}
	return minimumWinningTokens
}

const PRIZE_ADJUSTMENT = 10000000000000

func solveEquations(machine *machine) ([2]int, bool) {
	/*
		We need to find (x,y) from equations:
		ax + by = c
		dx + ey = f

		Where:
	*/
	a := float64(machine.aVector[0])
	b := float64(machine.bVector[0])
	c := float64(machine.prizeLocation[0])
	d := float64(machine.aVector[1])
	e := float64(machine.bVector[1])
	f := float64(machine.prizeLocation[1])
	/*
		Find y in terms of x:
		ax + by = c
		y = (c-ax)/b

		Substitute into second equation:
		dx + e(c-ax)/b = f
		bdx + ec - aex = bf
		x(bd - ae) = bf - ec
		x = (bf - ec)/(bd - ae)
	*/
	x := (b*f - e*c)/(b*d - a*e)
	// If x is not an integer, there is no solution
	if (x != float64(int(x))) {
		return [2]int{0, 0}, false
	}
	/*
		Substitute calculated x to find y:
		y = (c - ax)/b
	*/
	y := (c - a*x)/b
	// If y is not an integer, there is no solution
	if (y != float64(int(y))) {
		return [2]int{0, 0}, false
	}

	// If either x or y are negative, there is no solution
	if (int(x) < 0 || int(y) < 0) {
		return [2]int{0, 0}, false
	}
	return [2]int{int(x), int(y)}, true
}

func main() {
	dat, _ := os.ReadFile("assets/day13-input.txt")
	machines := unpack(string(dat))

	start := time.Now()
	fmt.Println(findMinWinningTokensPart1(machines))
	fmt.Println(time.Since(start))

	start = time.Now()
	minWinningTokens := 0
	for _, prizeMachine := range machines {
		adjustedMachine := &machine{
			aVector: prizeMachine.aVector,
			bVector: prizeMachine.bVector,
			prizeLocation: [2]int{
				prizeMachine.prizeLocation[0] + PRIZE_ADJUSTMENT,
				prizeMachine.prizeLocation[1] + PRIZE_ADJUSTMENT,
			},
		}
		solution, found := solveEquations(adjustedMachine)
		if (found) {
			minWinningTokens += solution[0]*3 + solution[1]
		}
	}
	fmt.Println(minWinningTokens)
	fmt.Println(time.Since(start))
}
