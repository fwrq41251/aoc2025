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
	sum := 0
	for _, line := range lines {
		sum += getMax(line)
	}
	return sum
}

func getMax(line string) int {
	max1 := 0
	max2 := 0
	for i := 0; i < len(line)-1; i++ {
		digit := int(line[i] - '0')
		if digit > max1 {
			max2 = 0
			max1 = digit
		} else if digit > max2 {
			max2 = digit
		}
	}

	lastDigit := int(line[len(line)-1] - '0')
	max2 = max(max2, lastDigit)
	return max1*10 + max2
}

func part2(lines []string) int {
	result := 0
	for _, line := range lines {
		maxList := getMaxList(line)
		number := getNumberByMaxList(maxList)
		result += number
	}
	return result
}

func getMaxList(line string) []int {
	maxList := []int{}
	startPos := 0
	for i := range 12 {
		tempMax := -1
		for j := startPos; j < len(line)-12+i+1; j++ {
			digit := int(line[j] - '0')
			if digit > tempMax {
				tempMax = digit
				startPos = j + 1
			}
		}
		maxList = append(maxList, tempMax)
	}
	return maxList
}

func getNumberByMaxList(maxList []int) int {
	result := 0
	multiplier := 1
	for i := len(maxList) - 1; i >= 0; i-- {
		result += maxList[i] * multiplier
		multiplier *= 10
	}
	return result
}
