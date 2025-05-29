package main

import (
	"fmt"
	"os"
	"strings"
	"strconv"
	"sort"
	"maps"
	"slices"
)

type Gate struct {
	operation string
	input1 string
	input2 string
	output string
}

func (g Gate) Use(wires map[string]int) int {
	if g.operation == "AND" { return wires[g.input1] & wires[g.input2] }
	if g.operation == "OR" { return wires[g.input1] | wires[g.input2] }
	// XOR
	return wires[g.input1] ^ wires[g.input2]
}
func (g Gate) Print() {
	fmt.Println(
		fmt.Sprintf(
			"%s %s %s -> %s",
			g.input1,
			g.operation,
			g.input2,
			g.output,
		),
	)
}

type GateCollection []*Gate

func (gc GateCollection) find(inputs []string, output string, operation string) []*Gate {
	result := []*Gate{}
	// Find gate in gateCollection which has inputs
	search:
	for _, gate := range gc {
		if (output != "" && gate.output != output) { continue }
		if (operation != "" && gate.operation != operation) { continue }

		for _, input := range inputs {
			if (gate.input1 != input && gate.input2 != input) { continue search }
		}

		result = append(result, gate)
	}
	return result
}

type GateSet struct {
	n string
	xor *Gate
	and *Gate
	mix *Gate
	carry *Gate
	z *Gate
}

func (gs *GateSet) Print() {
	fmt.Println("GateSet "+gs.n)
	fmt.Print("XOR: ")
	gs.xor.Print()
	fmt.Print("AND: ")
	gs.and.Print()
	fmt.Print("MIX: ")
	if (gs.mix == nil) {
		fmt.Println("nil")
	} else {
		gs.mix.Print()
	}
	fmt.Print("CARRY: ")
	gs.carry.Print()
	fmt.Print("Z: ")
	gs.z.Print()
}

func unpack(input string) (map[string]int, []string, map[string]*Gate, GateCollection) {
	// Return a map of wires with their values
	sections := strings.Split(input, "\n\n")
	inputWires := strings.Split(strings.Trim(sections[0], "\n"), "\n")
	gates := strings.Split(strings.Trim(sections[1], "\n"), "\n")

	// for each input wire, add to wires map
	wires := make(map[string]int)
	for _, inputWire := range inputWires {
		wireParts := strings.Split(inputWire, ": ")
		wireValue, _ := strconv.Atoi(wireParts[1])
		wires[wireParts[0]] = wireValue
	}
	//fmt.Println(wires)

	// for each gate, add output to wires map. create gate and map to output
	gateCollection := make(GateCollection, 0)
	gatesMap := make(map[string]*Gate)
	zWires := make([]string, 0)
	for _, gate := range gates {
		gateParts := strings.Split(gate, " ")
		outputWire := gateParts[4]
		newGate := &Gate{
			input1: gateParts[0],
			operation: gateParts[1],
			input2: gateParts[2],
			output: outputWire,
		}
		gatesMap[outputWire] = newGate
		gateCollection = append(gateCollection, newGate)
		
		// Add the output wire to the wires map
		if _, exists := wires[outputWire]; !exists { wires[outputWire] = 0 }

		// if gate output begins with z, add to zWires slice
		if outputWire[0] == 'z' {
			zWires = append(zWires, outputWire)
		}
		//fmt.Println(outputWire, newGate)
	}
	sort.Strings(zWires)
	slices.Reverse(zWires)
	return wires, zWires, gatesMap, gateCollection
}

func printRouteTo(wire string, gates map[string]*Gate) {
	toPrint := []*Gate{gates[wire]}
	for len(toPrint) > 0 {
		nextGate := toPrint[0]
		toPrint = toPrint[1:]
		if prevGate, exists := gates[nextGate.input1]; exists {
			toPrint = append(toPrint, prevGate)
		}
		if prevGate, exists := gates[nextGate.input2]; exists {
			toPrint = append(toPrint, prevGate)
		}
		nextGate.Print()
	}
}

func findGateByInput(gateType string, operation string, inputs [2][2]string, collection GateCollection) *Gate {
	gates := collection.find([]string{inputs[0][1], inputs[1][1]}, "", operation)
	if (len(gates) == 1) {
		return gates[0]
	} else if (len(gates) > 1) {
		fmt.Println("Multiple "+gateType+" gates found (returning first):")
		for _, foundGate := range gates { foundGate.Print() }
		return gates[0]
	}

	fmt.Println(gateType+" missing. Expected inputs "+inputs[0][1]+" and "+inputs[1][1]+".")
	for _, input := range inputs {
		inputGates := collection.find([]string{input[1]}, "", operation)
		if (len(inputGates) > 0) {
			fmt.Println("Found "+gateType+" with "+input[0]+" input ("+input[1]+").")
			for _, inputGate := range inputGates { inputGate.Print() }
			return inputGates[0]
		}
	}
	fmt.Println("No gate found. Returning nil.")
	return nil
}

func main() {
	dat, _ := os.ReadFile("./assets/day24-input.txt")
	wires, zWires, gates, gateCollection := unpack(string(dat))

	// while there is still a change, pulse over the wires
	isChange := true
	newWires := maps.Clone(wires)
	for isChange {
		isChange = false
		for wireName, wireValue := range newWires {
			gate, exists := gates[wireName]
			if !exists { continue }

			newValue := gate.Use(newWires)
			newWires[wireName] = newValue

			if newValue != wireValue {
				isChange = true
			}
		}
	}

	// iterate over the z wires to produce our output
	output := newWires[zWires[0]]
	for _, zWire := range zWires[1:] {
		output <<= 1
		output += newWires[zWire]
	}
	fmt.Println(output)

	// Part 2
	// for zWireN, gates for output are:
	// xN XOR yN (xorN)
	// xN-1 AND yN-1 (andN-1)
	// xorN AND andN-1 (mixN)
	// andN-1 OR mixN-1 (carryN-1)
	// xorN XOR carryN-1 (zN)
	sort.Strings(zWires)

	gateSets := make([]*GateSet, 46)
	xor0 := gateCollection.find([]string{"x00","y00"}, "", "XOR")[0]
	and0 := gateCollection.find([]string{"x00", "y00"}, "", "AND")[0]
	gateSets[0] = &GateSet{
		n: "00",
		xor: xor0,
		and: and0,
		mix: nil,
		carry: and0,
		z: xor0,
	}
	gateSets[0].Print()
	fmt.Println()

	for zIndex, zWire := range zWires {
		if (zIndex < 1 || zIndex > 44) { continue }
		fmt.Println("Performing calculations for "+zWire)

		//get the zWire number
		zWireNum := strings.TrimPrefix(zWire, "z")

		// Find xor gate
		xorGate := gateCollection.find([]string{"x"+zWireNum, "y"+zWireNum}, "", "XOR")[0]

		// Find and gate (won't exist for z45)
		andGate := gateCollection.find([]string{"x"+zWireNum, "y"+zWireNum}, "", "AND")[0]

		// Find mix gate (won't exist for z45)
		// mix := xorN AND carryN-1
		mixGate := findGateByInput(
			"mix",
			"AND",
			[2][2]string{
				[2]string{"xor", xorGate.output},
				[2]string{"carry-1", gateSets[zIndex-1].carry.output},
			},
			gateCollection,
		)

		// Find carry gate
		// carry := andN OR mixN
		carryGate := findGateByInput(
			"carry",
			"OR",
			[2][2]string{
				[2]string{"and", andGate.output},
				[2]string{"mix", mixGate.output},
			},
			gateCollection,
		)

		// Find z gate
		// z := xorN XOR carryN-1
		zGate := findGateByInput(
			"z",
			"XOR",
			[2][2]string{
				[2]string{"xor", xorGate.output},
				[2]string{"carry-1", gateSets[zIndex-1].carry.output},
			},
			gateCollection,
		)

		//create a GateSet for that number
		gateSet := GateSet{
			n: zWireNum,
			xor: xorGate,
			and: andGate,
			mix: mixGate,
			carry: carryGate,
			z: zGate,
		}
		gateSet.Print()
		gateSets[zIndex] = &gateSet
		fmt.Println()

		/**
		* Issues discovered by inspection:
		* GateSet 09:
		* XOR output is qwf, should be cnk. AND output is cnk, should be qwf.
		* GateSet 14:
		* AND output is z14, should be vhm. Z output is vhm, should be z14.
		* GateSet 27:
		* CARRY output is z27, should be mps. Z output is mps, should be z27.
		* GateSet 39:
		* MIX output is z39, should be msq. Z output is msq, should be z39.

		* Ordered list of swapped gates: cnk,mps,msq,qwf,vhm,z14,z27,z39
		**/
	}
}

