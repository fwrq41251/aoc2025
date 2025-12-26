package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	content, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}

	input := string(content)
	lines := strings.Split(strings.TrimSpace(input), "\n")

	fmt.Println("Part 1:", part1(lines))
	fmt.Println("Part 2:", part2(lines))
}

func part1(lines []string) int {
	adjacencyMap := make(map[string][]string)
	for _, line := range lines {
		node, neighbors := parseLine(line)
		adjacencyMap[node] = neighbors
	}

	return getPathCount(adjacencyMap, "you", "out", make(map[string]int))
}

func parseLine(line string) (string, []string) {
	before, after, _ := strings.Cut(line, ":")
	node := before
	leftStr := strings.TrimSpace(after)
	neighbors := strings.Split(leftStr, " ")

	return node, neighbors
}

// getPathCount returns the number of distinct paths from start to end using DFS with memoization.
func getPathCount(adjacencyMap map[string][]string, start string, end string, cache map[string]int) int {
	if start == end {
		return 1
	}

	if val, found := cache[start]; found {
		return val
	}

	totalPaths := 0
	for _, neighbor := range adjacencyMap[start] {
		totalPaths += getPathCount(adjacencyMap, neighbor, end, cache)
	}

	cache[start] = totalPaths
	return totalPaths
}

func part2(lines []string) int {
	adjacencyMap := make(map[string][]string)
	for _, line := range lines {
		node, neighbors := parseLine(line)
		adjacencyMap[node] = neighbors
	}

	start := "svr"
	end := "out"
	ponit1 := "dac"
	point2 := "fft"

	dacToFftPaths := getPathCount(adjacencyMap, ponit1, point2, make(map[string]int))
	if dacToFftPaths > 0 {
		startToDacPaths := getPathCount(adjacencyMap, start, ponit1, make(map[string]int))
		fftToOutPaths := getPathCount(adjacencyMap, point2, end, make(map[string]int))
		return startToDacPaths * dacToFftPaths * fftToOutPaths
	} else {
		fftToDacPaths := getPathCount(adjacencyMap, point2, ponit1, make(map[string]int))
		if fftToDacPaths == 0 {
			return 0
		} else {
			startToFftPaths := getPathCount(adjacencyMap, start, point2, make(map[string]int))
			dacToOutPaths := getPathCount(adjacencyMap, ponit1, end, make(map[string]int))
			return startToFftPaths * fftToDacPaths * dacToOutPaths
		}
	}
}
