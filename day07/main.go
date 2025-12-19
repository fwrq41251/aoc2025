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
	splitCounter := 0
	previousLine := make(map[int]string)

	for i := range lines {
		line := lines[i]
		if i == 0 {
			for j := 0; j < len(line); j++ {
				previousLine[j] = string(line[j])
			}
			continue
		}

		for j := 0; j < len(line); j++ {
			char := string(line[j])
			prevousChar := previousLine[j]
			if char == "^" {
				if prevousChar == "|" {
					previousLine[j-1] = "|"
					previousLine[j+1] = "|"
					splitCounter++
				}
				previousLine[j] = "^"
			} else if char == "." && (prevousChar == "S" || prevousChar == "|") {
				previousLine[j] = "|"
			} else {
				previousLine[j] = char
			}
		}
	}
	return splitCounter
}

func part2(lines []string) int {
	if len(lines) == 0 {
		return 0
	}

	width := len(lines[0])
	currentCounts := make([]int, width)

	totalTimelines := 0
	for j := range width {
		if lines[0][j] == 'S' {
			currentCounts[j] = 1
			totalTimelines = 1
			break
		}
	}

	for i := 1; i < len(lines); i++ {
		line := lines[i]
		nextCounts := make([]int, width)

		for col, count := range currentCounts {
			if count == 0 {
				continue
			}

			if col >= len(line) {
				continue
			}

			char := line[col]

			if char == '^' {
				totalTimelines += count

				if col-1 >= 0 {
					nextCounts[col-1] += count
				}
				if col+1 < width {
					nextCounts[col+1] += count
				}
			} else {
				nextCounts[col] += count
			}
		}

		currentCounts = nextCounts
	}

	return totalTimelines
}
