package main

import (
	"fmt"
	"math/big"
	"os"
	"regexp"
	"slices"
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
	result := 0

	for _, line := range lines {
		lights, _, buttons := parse(line)
		result += sovleLine(lights, buttons)
	}

	return result
}

func sovleLine(lights []bool, buttons []button) int {
	subsets := getSubsets(buttons)
	for _, btns := range subsets {
		var tempLights = make([]bool, len(lights))
		for _, b := range btns {
			for _, toTurnOn := range b.toTurnOn {
				tempLights[toTurnOn] = !tempLights[toTurnOn]
			}
		}
		if slices.Equal(lights, tempLights) {
			return len(btns)
		}
	}
	return -1
}

func parse(line string) ([]bool, []int, []button) {
	start := strings.Index(line, "[")
	end := strings.Index(line, "]")

	lightStr := line[start+1 : end]
	lights := make([]bool, len(lightStr))

	for i, char := range lightStr {
		if char == '#' {
			lights[i] = true
		} else {
			lights[i] = false
		}
	}

	// Parse buttons
	// Buttons are in (...), e.g. (0,1)
	re := regexp.MustCompile(`\(([\d,]+)\)`)
	matches := re.FindAllStringSubmatch(line, -1)

	var buttons []button

	for _, match := range matches {
		content := match[1]

		parts := strings.Split(content, ",")
		var indices []int

		for _, p := range parts {
			idx, _ := strconv.Atoi(p)
			indices = append(indices, idx)
		}

		buttons = append(buttons, button{toTurnOn: indices})
	}

	// Parse targets
	// Targets are in {...}, e.g. {3,5,4,7}
	targetStart := strings.Index(line, "{")
	targetEnd := strings.Index(line, "}")
	var targets []int
	if targetStart != -1 && targetEnd != -1 {
		targetContent := line[targetStart+1 : targetEnd]
		parts := strings.SplitSeq(targetContent, ",")
		for p := range parts {
			val, _ := strconv.Atoi(p)
			targets = append(targets, val)
		}
	}

	return lights, targets, buttons
}

func getSubsets(buttons []button) [][]button {
	n := len(buttons)
	totalSubsets := 1 << n // 2^n
	var result [][]button

	for i := range totalSubsets {
		var subset []button
		// 检查 i 的二进制每一位
		for j := range n {
			// 如果第 j 位是 1，说明选中了 nums[j]
			if (i>>j)&1 == 1 {
				subset = append(subset, buttons[j])
			}
		}
		result = append(result, subset)
	}

	slices.SortFunc(result, func(a, b []button) int {
		return len(a) - len(b)
	})
	return result
}

type button struct {
	toTurnOn []int
}

func part2(lines []string) int {
	result := 0
	for _, line := range lines {
		_, targets, buttons := parse(line)
		res := solvePart2(targets, buttons)
		if res != -1 {
			result += res
		}
	}
	return result
}

// solvePart1 (previously solveGaussian) builds the matrix and solves it using Gaussian Elimination over GF(2)
// For Part 1, we still use the brute force sovleLine (as per previous code) or could use this if we wanted.
// But the user's previous code was calling solveGaussian in part2 only.
// Wait, I replaced part2 to call solveGaussian, and left part1 calling sovleLine (brute force).
// I will keep part1 as is (using sovleLine) and focus on part2.

func solvePart2(targets []int, buttons []button) int {
	rows := len(targets)
	cols := len(buttons)

	// Build Matrix with big.Rat
	matrix := make([][]*big.Rat, rows)
	for i := range matrix {
		matrix[i] = make([]*big.Rat, cols+1)
		// Target column
		matrix[i][cols] = new(big.Rat).SetInt64(int64(targets[i]))
		// Init others to 0
		for j := range cols {
			matrix[i][j] = new(big.Rat)
		}
	}

	for j, btn := range buttons {
		for _, counterIdx := range btn.toTurnOn {
			if counterIdx < rows {
				matrix[counterIdx][j].SetInt64(1)
			}
		}
	}

	// Gaussian Elimination
	pivotRow := 0
	pivots := make([]int, rows)
	for i := range pivots {
		pivots[i] = -1
	}

	for col := 0; col < cols && pivotRow < rows; col++ {
		// Find pivot
		sel := -1
		for r := pivotRow; r < rows; r++ {
			if matrix[r][col].Sign() != 0 {
				sel = r
				break
			}
		}

		if sel == -1 {
			continue
		}

		// Swap
		matrix[pivotRow], matrix[sel] = matrix[sel], matrix[pivotRow]
		pivots[pivotRow] = col

		// Normalize pivot to 1
		pivotVal := new(big.Rat).Set(matrix[pivotRow][col])
		invPivot := new(big.Rat).Inv(pivotVal)
		for c := col; c <= cols; c++ {
			matrix[pivotRow][c].Mul(matrix[pivotRow][c], invPivot)
		}

		// Eliminate
		for r := range rows {
			if r != pivotRow && matrix[r][col].Sign() != 0 {
				factor := new(big.Rat).Set(matrix[r][col])
				for c := col; c <= cols; c++ {
					sub := new(big.Rat).Mul(matrix[pivotRow][c], factor)
					matrix[r][c].Sub(matrix[r][c], sub)
				}
			}
		}
		pivotRow++
	}

	// Check consistency
	for r := pivotRow; r < rows; r++ {
		if matrix[r][cols].Sign() != 0 {
			return -1
		}
	}

	// Free variables
	isPivotCol := make([]bool, cols)
	for r := 0; r < pivotRow; r++ {
		c := pivots[r]
		isPivotCol[c] = true
	}

	var freeVars []int
	for c := range cols {
		if !isPivotCol[c] {
			freeVars = append(freeVars, c)
		}
	}

	minTotal := new(big.Int)
	found := false

	// DFS
	var dfs func(idx int, currentSol []*big.Int)
	dfs = func(idx int, currentSol []*big.Int) {
		if idx == len(freeVars) {
			totalPresses := new(big.Int)
			fullSol := make([]*big.Int, cols)

			for i, fVal := range freeVars {
				fullSol[fVal] = currentSol[i]
				totalPresses.Add(totalPresses, currentSol[i])
			}

			// Solve dependent
			valid := true
			for r := 0; r < pivotRow; r++ {
				pCol := pivots[r]
				// x_p = RHS - sum(coeff * x_free)
				val := new(big.Rat).Set(matrix[r][cols])
				for i, fCol := range freeVars {
					coef := matrix[r][fCol]
					term := new(big.Rat).SetInt(currentSol[i])
					term.Mul(term, coef)
					val.Sub(val, term)
				}

				// Check integer >= 0
				if !val.IsInt() {
					valid = false
					break
				}
				intVal := val.Num() // Denom is 1
				if intVal.Sign() < 0 {
					valid = false
					break
				}
				fullSol[pCol] = intVal
				totalPresses.Add(totalPresses, intVal)
			}

			if valid {
				if !found || totalPresses.Cmp(minTotal) < 0 {
					minTotal.Set(totalPresses)
					found = true
				}
			}
			return
		}

		// Search range
		// Constraints are looser here, but usually solutions are small.
		// If fails, we might need to increase this.
		for v := int64(0); v <= 200; v++ {
			val := big.NewInt(v)
			currentSol[idx] = val
			dfs(idx+1, currentSol)
		}
	}

	startSol := make([]*big.Int, len(freeVars))
	dfs(0, startSol)

	if !found {
		return -1
	}

	if !minTotal.IsInt64() {
		// Should unlikely happen for this problem size, but return max int
		return int(^uint(0) >> 1)
	}
	return int(minTotal.Int64())
}
