package main

import (
	"fmt"
	"os"
	"strings"
	"regexp"
)

func unpack(input string) (patterns []string, designs []string) {
	lines := strings.Split(strings.Trim(input, "\n"), "\n")
	patterns = strings.Split(lines[0], ", ")
	designs = lines[2:]
	return patterns, designs
}

func variations(design string, patterns []string, cache map[string]int) int {
	if len(design) == 0 { return 1 }
	if numVariations, found := cache[design]; found {
		return numVariations
	}

	total := 0
	for _, pattern := range patterns {
		if strings.HasPrefix(design, pattern) {
			suffix := strings.TrimPrefix(design, pattern)
			suffixVariations := variations(suffix, patterns, cache)
			cache[suffix] = suffixVariations
			total += suffixVariations
		}
	}
	return total
}

func main() {
	dat, _ := os.ReadFile("./assets/day19-input.txt")
	patterns, designs := unpack(string(dat))

	// Part 1
	regexPattern := fmt.Sprintf(`^(?:(%s))*$`, strings.Join(patterns, ")|("))
	total := 0
	for _, design := range designs {
		matched, _ := regexp.Match(regexPattern, []byte(design))
		if matched { total += 1 }
	}
	fmt.Println(total)

	// Part 2
	total = 0
	cache := make(map[string]int)
	for _, design := range designs {
		total += variations(design, patterns, cache)
	}
	fmt.Println(total)
}

