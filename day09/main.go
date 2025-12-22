package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

type point struct {
	x, y int
}

func main() {
	content, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}

	input := string(content)
	lines := strings.Split(strings.TrimSpace(input), "\n")

	points := []point{}
	for _, line := range lines {
		var x, y int
		fmt.Sscanf(line, "%d,%d", &x, &y)
		points = append(points, point{x, y})
	}

	fmt.Println("Part 1:", part1(points))
	fmt.Println("Part 2:", part2(points))
}

func part1(points []point) int {
	maxArea := 0
	for i := range points {
		for j := i + 1; j < len(points); j++ {
			area := getArea(points[i], points[j])
			if area > maxArea {
				maxArea = area
			}
		}
	}
	return maxArea
}

func part2(points []point) int {
	// Coordinate Compression
	// We need to represent intervals.
	// Each point (x, y) effectively occupies [x, x+1) x [y, y+1)
	// Segments connect these.

	// Collect unique coordinates needed for the grid boundaries.
	// For every coordinate x, we might have boundaries at x and x+1.
	// We also need padding to ensure we can flood fill from outside.

	uniqueX := make(map[int]bool)
	uniqueY := make(map[int]bool)

	for _, p := range points {
		uniqueX[p.x] = true
		uniqueX[p.x+1] = true
		uniqueY[p.y] = true
		uniqueY[p.y+1] = true
	}

	// Add bounds for padding
	minX, maxX := points[0].x, points[0].x
	minY, maxY := points[0].y, points[0].y
	for _, p := range points {
		if p.x < minX {
			minX = p.x
		}
		if p.x > maxX {
			maxX = p.x
		}
		if p.y < minY {
			minY = p.y
		}
		if p.y > maxY {
			maxY = p.y
		}
	}
	uniqueX[minX-1] = true
	uniqueX[maxX+2] = true
	uniqueY[minY-1] = true
	uniqueY[maxY+2] = true

	sortedX := make([]int, 0, len(uniqueX))
	for x := range uniqueX {
		sortedX = append(sortedX, x)
	}
	sort.Ints(sortedX)

	sortedY := make([]int, 0, len(uniqueY))
	for y := range uniqueY {
		sortedY = append(sortedY, y)
	}
	sort.Ints(sortedY)

	xMap := make(map[int]int)
	for i, x := range sortedX {
		xMap[x] = i
	}
	yMap := make(map[int]int)
	for i, y := range sortedY {
		yMap[y] = i
	}

	W := len(sortedX) - 1
	H := len(sortedY) - 1

	// grid[y][x] represents the rectangular area:
	// X range: [sortedX[x], sortedX[x+1])
	// Y range: [sortedY[y], sortedY[y+1])
	// 0: Empty/Outside, 1: Boundary/Inside
	grid := make([][]int, H)
	for i := range grid {
		grid[i] = make([]int, W)
	}

	// Draw boundaries
	count := len(points)
	for i := range count {
		p1 := points[i]
		p2 := points[(i+1)%count]

		// Determine the continuous range covered by this segment
		// Segment from tile p1 to p2 (inclusive).
		// X range: [min(p1.x, p2.x), max(p1.x, p2.x) + 1)
		// Y range: [min(p1.y, p2.y), max(p1.y, p2.y) + 1)

		xStart, xEnd := p1.x, p2.x
		if xStart > xEnd {
			xStart, xEnd = xEnd, xStart
		}
		xEnd++ // Exclusive upper bound

		yStart, yEnd := p1.y, p2.y
		if yStart > yEnd {
			yStart, yEnd = yEnd, yStart
		}
		yEnd++ // Exclusive upper bound

		// Map to grid indices
		ixStart := xMap[xStart]
		ixEnd := xMap[xEnd]
		iyStart := yMap[yStart]
		iyEnd := yMap[yEnd]

		for y := iyStart; y < iyEnd; y++ {
			for x := ixStart; x < ixEnd; x++ {
				grid[y][x] = 1 // Mark as boundary
			}
		}
	}

	// Flood fill from outside (0, 0)
	// Since we added padding min-1, the cell (0,0) is guaranteed to be outside.
	floodFill(grid, 0, 0, 2) // Mark outside as 2

	// Build Weighted Prefix Sum
	// ps[y][x] stores the sum of AREAS of filled tiles in region (0,0) to (x,y)
	ps := make([][]int, H+1)
	for i := range ps {
		ps[i] = make([]int, W+1)
	}

	for y := range H {
		for x := range W {
			val := 0
			if grid[y][x] != 2 { // It is Boundary (1) or Inside (0 which wasn't reached by flood fill)
				width := sortedX[x+1] - sortedX[x]
				height := sortedY[y+1] - sortedY[y]
				val = width * height
			}

			ps[y+1][x+1] = val + ps[y][x+1] + ps[y+1][x] - ps[y][x]
		}
	}

	// Check all pairs
	maxArea := 0
	for i := range count {
		for j := i + 1; j < count; j++ {
			p1 := points[i]
			p2 := points[j]

			geoArea := getArea(p1, p2)
			if geoArea <= maxArea {
				continue
			}

			xStart, xEnd := p1.x, p2.x
			if xStart > xEnd {
				xStart, xEnd = xEnd, xStart
			}
			xEnd++ // Exclusive

			yStart, yEnd := p1.y, p2.y
			if yStart > yEnd {
				yStart, yEnd = yEnd, yStart
			}
			yEnd++ // Exclusive

			// Query PS
			// Map geometric coords to grid indices
			ixStart := xMap[xStart]
			ixEnd := xMap[xEnd]
			iyStart := yMap[yStart]
			iyEnd := yMap[yEnd]

			filledArea := querySum(ps, ixStart, iyStart, ixEnd, iyEnd)

			if filledArea == geoArea {
				maxArea = geoArea
			}
		}
	}

	return maxArea
}

func floodFill(grid [][]int, x, y, val int) {
	H := len(grid)
	W := len(grid[0])

	q := []point{{x, y}}
	grid[y][x] = val

	for len(q) > 0 {
		curr := q[len(q)-1]
		q = q[:len(q)-1]

		dirs := []point{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
		for _, d := range dirs {
			nx, ny := curr.x+d.x, curr.y+d.y
			if nx >= 0 && ny >= 0 && nx < W && ny < H {
				if grid[ny][nx] == 0 {
					grid[ny][nx] = val
					q = append(q, point{nx, ny})
				}
			}
		}
	}
}

func querySum(ps [][]int, x1, y1, x2, y2 int) int {
	// Query rectangular region [x1, x2) x [y1, y2) in grid indices
	// ps is 1-based size (W+1, H+1)
	return ps[y2][x2] - ps[y1][x2] - ps[y2][x1] + ps[y1][x1]
}

func getArea(p1, p2 point) int {
	return (abs(p1.x-p2.x) + 1) * (abs(p1.y-p2.y) + 1)
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
