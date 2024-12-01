package main

import (
	"bufio"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	left, right, err := readInput()
	if err != nil {
		return
	}

	sort.Ints(left)
	sort.Ints(right)

	difference := calculateDifference(left, right)
	println(difference)

	similarityScore := calculateSimilarityScore(left, right)
	println(similarityScore)
}

func calculateDifference(left []int, right []int) int {
	distanceSum := 0
	for i := range len(left) {
		l := left[i]
		r := right[i]
		distance := int(math.Abs(float64(l - r)))
		distanceSum += distance
	}
	return distanceSum
}

func calculateSimilarityScore(left []int, right []int) int {
	rightOccurrences := make(map[int]int)
	for _, r := range right {
		if _, ok := rightOccurrences[r]; !ok {
			rightOccurrences[r] = 0
		}

		rightOccurrences[r]++
	}

	score := 0
	for _, l := range left {
		score += l * rightOccurrences[l]
	}
	return score
}

func readInput() ([]int, []int, error) {
	file, err := os.Open("input.txt")
	if err != nil {
		println("Failed to open input.txt: " + err.Error())
		return nil, nil, err
	}
	defer file.Close()

	var (
		leftList  []int
		rightList []int
	)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		splitLine := strings.Split(line, "   ")

		left, err := strconv.Atoi(splitLine[0])
		if err != nil {
			println("Failed to convert input to number: " + err.Error())
			return nil, nil, err
		}

		right, err := strconv.Atoi(splitLine[1])
		if err != nil {
			println("Failed to convert input to number: " + err.Error())
			return nil, nil, err
		}

		leftList = append(leftList, left)
		rightList = append(rightList, right)
	}

	return leftList, rightList, nil
}
