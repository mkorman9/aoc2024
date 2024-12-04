package main

import (
	"bufio"
	"os"
)

const (
	Xmas    = "XMAS"
	XmasLen = len(Xmas)
)

func main() {
	lines, err := readInput()
	if err != nil {
		return
	}

	println(countXmas(lines))
	println(countMas(lines))
}

func countXmas(lines []string) int {
	count := 0

	for i, line := range lines {
		for j := range line {
			count += countXmasInPoint(lines, i, j)
		}
	}

	return count
}

func countXmasInPoint(lines []string, y, x int) int {
	count := 0
	match := 0

	// right
	for i := x; i < x+XmasLen && i < len(lines[y]); i++ {
		if lines[y][i] != Xmas[i-x] {
			match = 0
			break
		}

		match++
	}
	if match == XmasLen {
		count += 1
	}
	match = 0

	// left
	for i := x; i > x-XmasLen && i >= 0; i-- {
		if lines[y][i] != Xmas[x-i] {
			match = 0
			break
		}

		match++
	}
	if match == XmasLen {
		count += 1
	}
	match = 0

	// down
	for i := y; i < y+XmasLen && i < len(lines); i++ {
		if lines[i][x] != Xmas[i-y] {
			match = 0
			break
		}

		match++
	}
	if match == XmasLen {
		count += 1
	}
	match = 0

	// up
	for i := y; i > y-XmasLen && i >= 0; i-- {
		if lines[i][x] != Xmas[y-i] {
			match = 0
			break
		}

		match++
	}
	if match == XmasLen {
		count += 1
	}
	match = 0

	// right-down
	i, j := x, y
	for i < x+XmasLen && i < len(lines[y]) && j < y+XmasLen && j < len(lines) {
		if lines[i][j] != Xmas[i-x] {
			match = 0
			break
		}

		match++
		i++
		j++
	}
	if match == XmasLen {
		count += 1
	}
	match = 0

	// right-up
	i, j = x, y
	for i < x+XmasLen && i < len(lines[y]) && j > y-XmasLen && j >= 0 {
		if lines[i][j] != Xmas[i-x] {
			match = 0
			break
		}

		match++
		i++
		j--
	}
	if match == XmasLen {
		count += 1
	}
	match = 0

	// left-down
	i, j = x, y
	for i > x-XmasLen && i >= 0 && j < y+XmasLen && j < len(lines) {
		if lines[i][j] != Xmas[x-i] {
			match = 0
			break
		}

		match++
		i--
		j++
	}
	if match == XmasLen {
		count += 1
	}
	match = 0

	// left-up
	i, j = x, y
	for i > x-XmasLen && i >= 0 && j > y-XmasLen && j >= 0 {
		if lines[i][j] != Xmas[x-i] {
			match = 0
			break
		}

		match++
		i--
		j--
	}
	if match == XmasLen {
		count += 1
	}
	match = 0

	return count
}

func countMas(lines []string) int {
	count := 0

	for i, line := range lines {
		for j := range line {
			count += countMasInPoint(lines, i, j)
		}
	}

	return count
}

func countMasInPoint(lines []string, y, x int) int {
	if lines[y][x] != 'A' || x-1 < 0 || x+1 >= len(lines[y]) || y-1 < 0 || y+1 >= len(lines) {
		return 0
	}

	if lines[y-1][x-1] == 'M' && lines[y+1][x+1] == 'S' && lines[y-1][x+1] == 'M' && lines[y+1][x-1] == 'S' {
		return 1
	}
	if lines[y-1][x-1] == 'S' && lines[y+1][x+1] == 'M' && lines[y-1][x+1] == 'M' && lines[y+1][x-1] == 'S' {
		return 1
	}
	if lines[y-1][x-1] == 'M' && lines[y+1][x+1] == 'S' && lines[y-1][x+1] == 'S' && lines[y+1][x-1] == 'M' {
		return 1
	}
	if lines[y-1][x-1] == 'S' && lines[y+1][x+1] == 'M' && lines[y-1][x+1] == 'S' && lines[y+1][x-1] == 'M' {
		return 1
	}

	return 0
}

func readInput() ([]string, error) {
	file, err := os.Open("input.txt")
	if err != nil {
		println("Failed to open input.txt: " + err.Error())
		return nil, err
	}
	defer file.Close()

	var (
		lines []string
	)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	return lines, nil
}
