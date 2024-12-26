package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"slices"
)

type Registers [3]int

func toInt(str string) int {
	res, _ := strconv.Atoi(str)
	return res
}

func unpack(input string) (Registers, []int) {
	registers := [3]int{}
	re := regexp.MustCompile(`Register .: (\d+)`)
	for registerIndex, register := range re.FindAllStringSubmatch(input, -1) {
		res, _ := strconv.Atoi(register[1])
		registers[registerIndex] = res
	}

	re = regexp.MustCompile(`Program: (.*)`)
	program := []int{}
	programInputs := strings.Split(re.FindStringSubmatch(input)[1], ",")
	for i := 0; i < len(programInputs); i += 1 {
		program = append(program, toInt(programInputs[i]))
	}
	return registers, program
}

func getComboValue(combo int, registers Registers) int {
	if (combo < 4) {
		return combo
	}
	if (combo < 7) {
		return registers[combo - 4]
	}
	// TODO combo = 7 is an error
	panic("Combo value above 6!")
}

// --- Operations ---

func pow2(exponent int, storeLocation int, registers Registers) Registers {
	pow2Result := 1
	for range exponent {
		pow2Result *= 2
	}
	newRegisters := [3]int{registers[0], registers[1], registers[2]}
	newRegisters[storeLocation] = registers[0] / pow2Result
	return newRegisters
}

func adv(operand int, registers Registers) Registers {
	return pow2(getComboValue(operand, registers), 0, registers)
}

func bdv(operand int, registers Registers) Registers {
	return pow2(getComboValue(operand, registers), 1, registers)
}

func cdv(operand int, registers Registers) Registers {
	return pow2(getComboValue(operand, registers), 2, registers)
}

func bxl(operand int, registers Registers) Registers {
	return [3]int{registers[0], registers[1] ^ operand, registers[2]}
}

func bst(operand int, registers Registers) Registers {
	return [3]int{registers[0], getComboValue(operand, registers) % 8, registers[2]}
}

func bxc(registers Registers) Registers {
	return [3]int{registers[0], registers[1] ^ registers[2], registers[2]}
}

func out(operand int, registers Registers) int {
	return getComboValue(operand, registers) % 8
}

func runProgram(program []int, registers Registers) []int {
	output := []int{}
	for i := 0; i < len(program); {
		//fmt.Printf("%d | %+v | %+v\n", i, registers, program)
		switch program[i] {
		case 0:
			// Combo
			registers = adv(program[i+1], registers)
		case 1:
			registers = bxl(program[i+1], registers)
		case 2:
			// Combo
			registers = bst(program[i+1], registers)
		case 3:
			// jnz
			if registers[0] != 0 {
				i = program[i+1]
				continue
			}
		case 4:
			registers = bxc(registers)
		case 5:
			output = append(output, out(program[i+1], registers))
			if output[len(output)-1] != program[len(output)-1] {
				break
			}
		case 6:
			// Combo
			registers = bdv(program[i+1], registers)
		case 7:
			// Combo
			registers = cdv(program[i+1], registers)
		}
		i += 2
	}
	return output
}

// Finds the initial value of Register A which produces the reverse of program
func findReverse(program []int) []int {
	prevReverses := []int{0}
	nextReverses := []int{}
	for i := 1; i <= len(program); i++ {
		programSoFar := program[len(program) - i:]

		for _, prevReverse := range prevReverses {
			for j := 0; j < 8; j++ {
				// Run program and if the output matches programSoFar, add to reverses slice
				if slices.Equal(runProgram(program, [3]int{prevReverse + j, 0, 0}), programSoFar) {
					nextReverses = append(nextReverses, prevReverse + j)
				}
			}
		}

		if i == len(program) {
			break
		}

		// For each reverse slices we want to test, bitshift it 3 times
		prevReverses = []int{}
		for _, reverse := range nextReverses {
			prevReverses = append(prevReverses, (reverse << 3))
		}
		nextReverses = []int{}
	}
	return nextReverses
}

func main() {
	dat, _ := os.ReadFile("./assets/day17-input.txt")
	registers, program := unpack(string(dat))

	fmt.Println(runProgram(program, registers))

	fmt.Println(findReverse(program))

	// Double checking the lowest result here
	fmt.Println(runProgram(program, [3]int{164542125272765, 0, 0}))
}

