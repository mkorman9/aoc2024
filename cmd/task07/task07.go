package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type Calibration struct {
	Expected int
	Operands []int
}

type Operator int

const (
	OpAdd Operator = iota
	OpMul
	OpConcat
)

var operatorsCache = make(map[int][][]Operator)
var operatorsCacheWithConcat = make(map[int][][]Operator)

func main() {
	calibrations, err := readInput()
	if err != nil {
		return
	}

	println(calibrationsSum(calibrations, false))
	println(calibrationsSum(calibrations, true))
}

func calibrationsSum(calibrations []Calibration, concat bool) int {
	sum := 0

	for _, calibration := range calibrations {
		operatorCombinations := operatorsCombinations(len(calibration.Operands)-1, concat)
		for _, operators := range operatorCombinations {
			if execute(calibration.Operands, operators) == calibration.Expected {
				sum += calibration.Expected
				break
			}
		}
	}

	return sum
}

func execute(operands []int, operators []Operator) int {
	result := operands[0]

	for i := range operators {
		if operators[i] == OpAdd {
			result += operands[i+1]
		} else if operators[i] == OpMul {
			result *= operands[i+1]
		} else if operators[i] == OpConcat {
			r, err := strconv.Atoi(strconv.Itoa(result) + strconv.Itoa(operands[i+1]))
			if err != nil {
				panic(err)
			}

			result = r
		}
	}

	return result
}

func operatorsCombinations(length int, concat bool) [][]Operator {
	if length <= 0 {
		return nil
	}
	if concat {
		if c, ok := operatorsCacheWithConcat[length]; ok {
			return c
		}
	} else {
		if c, ok := operatorsCache[length]; ok {
			return c
		}
	}

	var result [][]Operator
	var combination []Operator

	var helper func(pos int)
	helper = func(pos int) {
		if pos == length {
			temp := make([]Operator, len(combination))
			copy(temp, combination)
			result = append(result, temp)
			return
		}

		combination = append(combination, OpAdd)
		helper(pos + 1)
		combination = combination[:len(combination)-1]

		combination = append(combination, OpMul)
		helper(pos + 1)
		combination = combination[:len(combination)-1]

		if concat {
			combination = append(combination, OpConcat)
			helper(pos + 1)
			combination = combination[:len(combination)-1]
		}
	}

	helper(0)

	if concat {
		operatorsCacheWithConcat[length] = result
	} else {
		operatorsCache[length] = result
	}

	return result
}

func readInput() ([]Calibration, error) {
	file, err := os.Open("input.txt")
	if err != nil {
		println("Failed to open input.txt: " + err.Error())
		return nil, err
	}
	defer file.Close()

	var (
		calibrations []Calibration
	)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ": ")

		expected, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, err
		}

		var ops []int
		for _, operand := range strings.Split(parts[1], " ") {
			op, err := strconv.Atoi(operand)
			if err != nil {
				return nil, err
			}
			ops = append(ops, op)
		}

		calibrations = append(calibrations, Calibration{expected, ops})
	}

	return calibrations, nil
}
