package main

import (
	"io"
	"os"
	"regexp"
	"strconv"
)

func main() {
	input, err := readInput()
	if err != nil {
		return
	}

	println(calculateSum(input))
	println(calculateSumWithDos(input))
}

func calculateSum(input []byte) int {
	pattern := regexp.MustCompile("mul\\((\\d{1,3}),(\\d{1,3})\\)")
	matches := pattern.FindAllStringSubmatch(string(input), -1)

	sum := 0
	for _, match := range matches {
		a, err := strconv.Atoi(match[1])
		if err != nil {
			panic("Failed to convert number " + err.Error())
		}

		b, err := strconv.Atoi(match[2])
		if err != nil {
			panic("Failed to convert number " + err.Error())
		}

		sum += a * b
	}

	return sum
}

func calculateSumWithDos(input []byte) int {
	doPattern := regexp.MustCompile("do\\(\\)")
	doMatches := doPattern.FindAllStringSubmatchIndex(string(input), -1)

	dontPattern := regexp.MustCompile("don't\\(\\)")
	dontMatches := dontPattern.FindAllStringSubmatchIndex(string(input), -1)

	var enabledPartStart []int
	for _, match := range doMatches {
		enabledPartStart = append(enabledPartStart, match[1])
	}

	var disabledPartStart []int
	for _, match := range dontMatches {
		disabledPartStart = append(disabledPartStart, match[1])
	}

	var invalidRanges [][]int
	for _, disabledStartIndex := range disabledPartStart {
		found := false

		for _, enabledStartIndex := range enabledPartStart {
			if enabledStartIndex > disabledStartIndex {
				invalidRanges = append(invalidRanges, []int{disabledStartIndex, enabledStartIndex})
				found = true
				break
			}
		}

		if !found {
			invalidRanges = append(invalidRanges, []int{disabledStartIndex, len(input)})
		}
	}

	pattern := regexp.MustCompile("mul\\((\\d{1,3}),(\\d{1,3})\\)")
	matches := pattern.FindAllStringSubmatchIndex(string(input), -1)

	sum := 0
	for _, match := range matches {
		index := match[0]

		a, err := strconv.Atoi(string(input)[match[2]:match[3]])
		if err != nil {
			panic("Failed to convert number " + err.Error())
		}

		b, err := strconv.Atoi(string(input)[match[4]:match[5]])
		if err != nil {
			panic("Failed to convert number " + err.Error())
		}

		ignore := false
		for _, r := range invalidRanges {
			if index >= r[0] && index < r[1] {
				ignore = true
				break
			}
		}

		if !ignore {
			sum += a * b
		}
	}

	return sum
}

func readInput() ([]byte, error) {
	file, err := os.Open("input.txt")
	if err != nil {
		println("Failed to open input.txt: " + err.Error())
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(file)
}
