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

func part1(grid []string) int {
	toRemove := 0
	rows := len(grid)
	cols := len(grid[0])

	dx := []int{-1, 0, 1, -1, 1, -1, 0, 1}
	dy := []int{-1, -1, -1, 0, 0, 1, 1, 1}

	for r := range rows {
		for c := range cols {
			cell := grid[r][c]
			if cell == '@' {
				neighborCount := 0
				for dir := range 8 {
					nr, nc := r+dy[dir], c+dx[dir]
					if nr >= 0 && nr < rows && nc >= 0 && nc < cols {
						neighbor := grid[nr][nc]
						if neighbor == cell {
							neighborCount++
						}
					}
				}
				if neighborCount < 4 {
					toRemove++
				}
			}
		}
	}
	return toRemove
}

func removeGrass(grid [][]byte) int {
	toRemove := 0

	rows := len(grid)
	cols := len(grid[0])

	dx := []int{-1, 0, 1, -1, 1, -1, 0, 1}
	dy := []int{-1, -1, -1, 0, 0, 1, 1, 1}

	for r := range rows {
		for c := range cols {
			cell := grid[r][c]
			if cell == '@' {
				neighborCount := 0
				for dir := range 8 {
					nr, nc := r+dy[dir], c+dx[dir]
					if nr >= 0 && nr < rows && nc >= 0 && nc < cols {
						neighbor := grid[nr][nc]
						if neighbor == cell {
							neighborCount++
						}
					}
				}
				if neighborCount < 4 {
					toRemove++
					grid[r][c] = 'x'
				}
			}
		}
	}
	return toRemove

}

func part2(lines []string) int {
	rows := len(lines)

	//convert string slice to byte grid
	grid := make([][]byte, rows)
	for i, row := range lines {
		grid[i] = []byte(row)
	}

	totalRemoved := 0
	for true {
		removed := removeGrass(grid)
		if removed > 0 {
			totalRemoved += removed
		} else {
			break
		}
	}
	return totalRemoved
}
