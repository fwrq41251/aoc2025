package main

import (
	"fmt"
	"os"
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

func part1(lines []string) int {
	currentPos := 50
	zeroCount := 0

	for _, line := range lines {
		direction := line[0]
		numberStr := line[1:]

		steps, err := strconv.Atoi(numberStr)
		if err != nil {
			fmt.Println("Error converting steps:", err)
			continue
		}

		switch direction {
		case 'R':
			currentPos += steps
		case 'L':
			currentPos -= steps
		}

		currentPos = currentPos % 100

		if currentPos < 0 {
			currentPos += 100
		}

		if currentPos == 0 {
			zeroCount++
		}
	}

	return zeroCount
}

func part2(lines []string) int {
	currentPos := 50
	zeroCount := 0

	for _, line := range lines {
		prevPos := currentPos
		direction := line[0]
		numberStr := line[1:]

		steps, err := strconv.Atoi(numberStr)
		if err != nil {
			fmt.Println("Error converting steps:", err)
			continue
		}

		switch direction {
		case 'R':
			currentPos += steps
		case 'L':
			currentPos -= steps
		}

		if currentPos == 0 {
			zeroCount++
		} else if currentPos < 0 {
			toAdd := ((-currentPos) / 100)
			if prevPos != 0 {
				toAdd++
			}
			zeroCount += toAdd
		} else if currentPos >= 100 {
			toAdd := (currentPos / 100)
			zeroCount += toAdd
		}

		currentPos = currentPos % 100

		if currentPos < 0 {
			currentPos += 100
		}
	}

	return zeroCount
}
