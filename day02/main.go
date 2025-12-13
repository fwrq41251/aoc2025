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
	invalidSum := 0

	for _, line := range lines {
		split := strings.SplitSeq(line, ",")
		for part := range split {
			part = strings.TrimSpace(part)
			splitPart := strings.Split(part, "-")
			if len(splitPart) != 2 {
				continue
			}

			start, err1 := strconv.Atoi(strings.TrimSpace(splitPart[0]))
			end, err2 := strconv.Atoi(strings.TrimSpace(splitPart[1]))
			if err1 != nil || err2 != nil {
				continue
			}

			for i := start; i <= end; i++ {
				if isRepeatNumber(i) {
					invalidSum += i
				}
			}
		}
	}
	return invalidSum
}

func isRepeatNumber(n int) bool {
	s := strconv.Itoa(n)
	length := len(s)

	if length%2 != 0 {
		return false
	}

	mid := length / 2

	firstHalf := s[:mid]
	secondHalf := s[mid:]

	return firstHalf == secondHalf
}

func part2(lines []string) int {
	invalidSum := 0

	for _, line := range lines {
		split := strings.SplitSeq(line, ",")
		for part := range split {
			part = strings.TrimSpace(part)
			splitPart := strings.Split(part, "-")
			if len(splitPart) != 2 {
				continue
			}

			start, err1 := strconv.Atoi(strings.TrimSpace(splitPart[0]))
			end, err2 := strconv.Atoi(strings.TrimSpace(splitPart[1]))
			if err1 != nil || err2 != nil {
				continue
			}

			for i := start; i <= end; i++ {
				if isRepeatNumberV2(i) {
					invalidSum += i
				}
			}
		}
	}
	return invalidSum
}

func isRepeatNumberV2(n int) bool {
	strN := strconv.Itoa(n)
	length := len(strN)

	divisors := getDivisors(length)
	for _, d := range divisors {
		segmentLength := length / d
		segment := strN[:segmentLength]
		result := strings.Repeat(segment, d)
		if result == strN {
			return true
		}
	}
	return false
}

func getDivisors(n int) []int {
	divisors := []int{}
	for i := 2; i <= n; i++ {
		if n%i == 0 {
			divisors = append(divisors, i)
		}
	}
	return divisors
}
