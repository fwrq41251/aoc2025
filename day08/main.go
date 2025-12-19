package main

import (
	"fmt"
	"os"
	"slices"
	"sort"
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
	positions := []Position{}
	for _, line := range lines {
		var pos Position
		fmt.Sscanf(line, "%d,%d,%d", &pos.x, &pos.y, &pos.z)
		positions = append(positions, pos)
	}

	connections := []Connection{}
	for i := 0; i < len(positions); i++ {
		for j := i + 1; j < len(positions); j++ {
			distSq := calculateDistance(positions[i], positions[j])
			connections = append(connections, Connection{a: i, b: j, distSq: distSq})
		}
	}

	sort.Slice(connections, func(i, j int) bool {
		return connections[i].distSq < connections[j].distSq
	})

	count := len(positions)
	uf := newUnionFind(count)
	max := 1000
	topConnections := connections[:max]
	for _, conn := range topConnections {
		uf.union(conn.a, conn.b)
	}

	sizes := []int{}
	for i := range count {
		if uf.parent[i] == i {
			sizes = append(sizes, uf.size[i])
		}
	}

	slices.SortFunc(sizes, func(a, b int) int {
		return b - a
	})

	result := 1
	for i := 0; i < 3 && i < len(sizes); i++ {
		result *= sizes[i]
	}

	return result
}

func part2(lines []string) int {
	positions := []Position{}
	for _, line := range lines {
		var pos Position
		fmt.Sscanf(line, "%d,%d,%d", &pos.x, &pos.y, &pos.z)
		positions = append(positions, pos)
	}

	connections := []Connection{}
	for i := 0; i < len(positions); i++ {
		for j := i + 1; j < len(positions); j++ {
			distSq := calculateDistance(positions[i], positions[j])
			connections = append(connections, Connection{a: i, b: j, distSq: distSq})
		}
	}

	sort.Slice(connections, func(i, j int) bool {
		return connections[i].distSq < connections[j].distSq
	})

	count := len(positions)
	uf := newUnionFind(count)

	result := 0
	i := 0

	for uf.groupCount > 1 {
		conn := connections[i]
		a := positions[conn.a]
		b := positions[conn.b]
		result = a.x * b.x
		uf.union(conn.a, conn.b)
		i++
	}

	return result
}

type Position struct {
	x int
	y int
	z int
}

type Connection struct {
	a      int
	b      int
	distSq int
}

func calculateDistance(a, b Position) int {
	return square(abs(a.x-b.x)) + square(abs(a.y-b.y)) + square(abs(a.z-b.z))
}

func square(n int) int {
	return n * n
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

type UnionFind struct {
	parent     []int
	size       []int
	groupCount int
}

func newUnionFind(n int) *UnionFind {
	parent := make([]int, n)
	size := make([]int, n)
	for i := range n {
		parent[i] = i
		size[i] = 1
	}
	return &UnionFind{parent: parent, size: size, groupCount: n}
}

func (uf *UnionFind) find(x int) int {
	if uf.parent[x] != x {
		uf.parent[x] = uf.find(uf.parent[x])
	}
	return uf.parent[x]
}

func (uf *UnionFind) union(x, y int) {
	rootX := uf.find(x)
	rootY := uf.find(y)
	if rootX != rootY {
		if uf.size[rootX] < uf.size[rootY] {
			rootX, rootY = rootY, rootX
		}
		uf.parent[rootY] = rootX
		uf.size[rootX] += uf.size[rootY]
		uf.groupCount--
	}
}
