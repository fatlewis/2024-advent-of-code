package main

import (
	"fmt"
	"os"
	"strings"
	"sort"
	"maps"
	"slices"
	"container/heap"
)

// PriorityQueue of pools with max potential neighbours
type Pool struct {
	items map[string]bool
	potentials map[string]bool
}
func (p Pool) Clone() *Pool {
	return &Pool{
		items: maps.Clone(p.items),
		potentials: maps.Clone(p.potentials),
	}
}

type PoolHeap []*Pool
func (h PoolHeap) Len() int { return len(h) }
func (h PoolHeap) Less(a, b int) bool {
	return (len(h[a].items)+len(h[a].potentials)) > (len(h[b].items)+len(h[b].potentials))
}
func (h PoolHeap) Swap(a, b int) { h[a], h[b] = h[b], h[a] }

func (h *PoolHeap) Push(a any) {
	*h = append(*h, a.(*Pool))
}
func (h *PoolHeap) Pop() any {
	old := *h
	result := old[len(old)-1]
	*h = old[0:len(old)-1]
	return result
}

func unpack(input string) (map[string]map[string]bool, []string) {
	lines := strings.Split(strings.Trim(input, "\n"), "\n")
	// Map of [computerName]->map[computerName]bool
	tComputers := []string{}
	computerLinks := make(map[string]map[string]bool)
	for _, line := range lines {
		computers := strings.Split(line, "-")
		if _, exists := computerLinks[computers[0]]; !exists {
			computerLinks[computers[0]] = make(map[string]bool)
			if computers[0][0] == 't' { tComputers = append(tComputers, computers[0]) }
		}
		if _, exists := computerLinks[computers[1]]; !exists {
			computerLinks[computers[1]] = make(map[string]bool)
			if computers[1][0] == 't' { tComputers = append(tComputers, computers[1]) }
		}
		computerLinks[computers[0]][computers[1]] = true
		computerLinks[computers[1]][computers[0]] = true
	}
	return computerLinks, tComputers
}

func main() {
	dat, _ := os.ReadFile("./assets/day23-input.txt")
	computerLinks, tComputers := unpack(string(dat))

	resultMap := make(map[string]bool)
	for _, tComputer := range tComputers {
		for tLink, _ := range computerLinks[tComputer] {
			for lastLink, _ := range computerLinks[tLink] {
				if lastLink == tComputer { continue }
				// If there's a link back from the last link to the first computer
				if computerLinks[lastLink][tComputer] {
					// Add to resultMap
					result := []string{tComputer, tLink, lastLink}
					sort.Strings(result)
					resultMap[fmt.Sprintf("%v", result)] = true
				}
			}
		}
	}
	count := 0
	for range resultMap { count++ }
	fmt.Println(count)

	// Part 2 Heap approach
	poolHeap := &PoolHeap{}
	// Generate pools for each computer and add to heap
	computers := slices.Collect(maps.Keys(computerLinks))
	for _, computer := range computers {
		pool := &Pool{
			items: map[string]bool{computer:true},
			potentials: maps.Clone(computerLinks[computer]),
		}
		heap.Push(poolHeap, pool)
	}

	var maxPool *Pool
	for true {
		// Pop pool off heap
		pool := heap.Pop(poolHeap).(*Pool)
		// If no potentials, we've found the max, break.
		if len(pool.potentials) == 0 {
			maxPool = pool
			break
		}

		// Try next potential link
		nextLink := ""
		newPotentials := make(map[string]bool)
		// Filter remaining potentials based on whether it links with them
		for potential := range pool.potentials {
			if nextLink == "" {
				nextLink = potential
			} else if computerLinks[nextLink][potential] {
				newPotentials[potential] = true
			}
		}

		// Push pool with it added and every potential which links
		newPool := pool.Clone()
		newPool.items[nextLink] = true
		newPool.potentials = newPotentials

		heap.Push(poolHeap, newPool)
		// If remaining potentials are fewer than every potential
		if len(newPotentials) < len(pool.potentials)-1 {
			// Push pool without it added and with every remaining potential
			otherPool := pool.Clone()
			delete(otherPool.potentials, nextLink)
			heap.Push(poolHeap, otherPool)
		}
	}
	resultComputers := slices.Collect(maps.Keys(maxPool.items))
	sort.Strings(resultComputers)
	resultStr := strings.Join(resultComputers, ",")
	fmt.Println(resultStr)
}

