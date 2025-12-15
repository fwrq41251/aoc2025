package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
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

type Range struct {
	start, end int
}

func part1(lines []string) int {
	ranges := []Range{}
	freshCount := 0

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.Contains(line, "-") {
			parts := strings.Split(line, "-")
			start, _ := strconv.Atoi(parts[0])
			end, _ := strconv.Atoi(parts[1])
			ranges = append(ranges, Range{start: start, end: end})
		} else {
			num, _ := strconv.Atoi(line)
			for _, r := range ranges {
				if num >= r.start && num <= r.end {
					freshCount++
					break
				}
			}
		}
	}
	return freshCount
}

func part2(lines []string) int {
	ranges := []Range{}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "-") {
			parts := strings.Split(line, "-")
			start, _ := strconv.Atoi(parts[0])
			end, _ := strconv.Atoi(parts[1])
			ranges = append(ranges, Range{start: start, end: end})
		}
	}

	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i].start < ranges[j].start
	})

	mergedRanges := []Range{}
	mergedRanges = append(mergedRanges, ranges[0])
	for i := 1; i < len(ranges); i++ {
		current := ranges[i]

		lastIndex := len(mergedRanges) - 1
		last := mergedRanges[lastIndex]

		if current.start <= last.end {
			mergedRanges[lastIndex].end = max(last.end, current.end)
		} else {
			mergedRanges = append(mergedRanges, current)
		}
	}

	freshCount := 0
	for _, r := range mergedRanges {
		freshCount += (r.end - r.start + 1)
	}
	return freshCount
}
