package main

import (
	"bufio"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	reports, err := readInput()
	if err != nil {
		return
	}

	safeReports := countSafeReports(reports)
	safeReportsDampener := countSafeReportsDampener(reports)

	println(safeReports)
	println(safeReportsDampener)
}

func countSafeReports(reports [][]int) int {
	result := 0

	for _, report := range reports {
		if isSafe(report) {
			result++
		}
	}

	return result
}

func countSafeReportsDampener(reports [][]int) int {
	result := 0

	for _, report := range reports {
		if isSafe(report) {
			result++
			continue
		}

		for i := 0; i < len(report); i++ {
			report2 := skip(report, i)
			if isSafe(report2) {
				result++
				break
			}
		}
	}

	return result
}

func skip(report []int, pos int) []int {
	var result []int
	for i := 0; i < len(report); i++ {
		if i == pos {
			continue
		}

		result = append(result, report[i])
	}

	return result
}

func isSafe(report []int) bool {
	increasing := report[0] < report[1]

	for i := 1; i < len(report); i++ {
		diff := int(math.Abs(float64(report[i] - report[i-1])))
		if diff < 1 || diff > 3 {
			return false
		}
		if (increasing && report[i] < report[i-1]) || (!increasing && report[i] > report[i-1]) {
			return false
		}
	}

	return true
}

func readInput() ([][]int, error) {
	file, err := os.Open("input.txt")
	if err != nil {
		println("Failed to open input.txt: " + err.Error())
		return nil, err
	}
	defer file.Close()

	var (
		reports [][]int
	)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		splitLine := strings.Split(line, " ")

		levels := make([]int, len(splitLine))
		for i, s := range splitLine {
			n, err := strconv.Atoi(s)
			if err != nil {
				return nil, err
			}

			levels[i] = n
		}

		reports = append(reports, levels)
	}

	return reports, nil
}
