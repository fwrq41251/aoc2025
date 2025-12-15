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
	length := len(lines)
	splitLines := [][]string{}
	for _, line := range lines {
		parts := strings.Fields(line)
		splitLines = append(splitLines, parts)
	}

	result := 0
	cols := len(splitLines[0])
	for i := range cols {
		operator := splitLines[length-1][i]
		acc, _ := strconv.Atoi(splitLines[0][i])

		for j := 1; j < length-1; j++ {
			value, _ := strconv.Atoi(splitLines[j][i])
			if operator == "+" {
				acc += value
			}
			if operator == "*" {
				acc *= value
			}
		}
		result += acc
	}

	return result
}

func part2(grid []string) int {
	cols := 0
	for _, line := range grid {
		if len(line) > cols {
			cols = len(line)
		}
	}
	//add emoty column at the end
	cols += 1
	for i := range grid {
		if len(grid[i]) < cols {
			grid[i] = grid[i] + strings.Repeat(" ", cols-len(grid[i]))
		}
	}

	result := 0
	numbers := []int{}
	start := 0
	end := 0

	for col := 0; col < cols; col++ {
		emptyCol := true
		var numBuilder strings.Builder

		for row := 0; row < len(grid)-1; row++ {
			char := grid[row][col]
			if char != ' ' {
				emptyCol = false
				numBuilder.WriteByte(char)
			}
		}

		if emptyCol {
			end = col
			op := getOperator(grid[len(grid)-1], start, end)

			var acc int
			if op == "*" {
				acc = 1
			} else {
				acc = 0
			}

			for _, n := range numbers {
				switch op {
				case "+":
					acc += n
				case "*":
					acc *= n
				}
			}

			result += acc
			numbers = []int{}
			start = col + 1
		} else {
			numStr := numBuilder.String()
			if numStr != "" {
				num, _ := strconv.Atoi(numStr)
				numbers = append(numbers, num)
			}
		}
	}
	return result
}

func getOperator(line string, start int, end int) string {
	for i := start; i < end; i++ {
		if line[i] == '+' || line[i] == '*' {
			return string(line[i])
		}
	}
	return ""
}
