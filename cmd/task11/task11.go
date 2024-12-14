package main

import (
	"io"
	"os"
	"strconv"
	"strings"
)

type Stone struct {
	Value int
	Next  *Stone
	Prev  *Stone
}

func main() {
	stone, err := readInput()
	if err != nil {
		return
	}

	println(processAndCountStones(stone, 25))
	println(processAndCountStones(stone, 75))
}

func processAndCountStones(stone *Stone, blinks int) int {
	for range blinks {
		current := stone
		for current != nil {
			valueString := strconv.Itoa(current.Value)

			if current.Value == 0 {
				current.Value = 1
			} else if len(valueString)%2 == 0 {
				stoneLeftValueStr := valueString[:len(valueString)/2]
				stoneRightValueStr := valueString[len(valueString)/2:]
				stoneLeftValue, err := strconv.Atoi(stoneLeftValueStr)
				if err != nil {
					panic(err)
				}
				stoneRightValue, err := strconv.Atoi(stoneRightValueStr)
				if err != nil {
					panic(err)
				}

				current.Value = stoneLeftValue
				newStone := &Stone{Value: stoneRightValue, Next: current.Next, Prev: current}
				if current.Next != nil {
					current.Next.Prev = newStone
				}
				current.Next = newStone
				current = current.Next
			} else {
				current.Value *= 2024
			}

			current = current.Next
		}
	}

	current := stone
	count := 0
	for current != nil {
		count++
		current = current.Next
	}
	return count
}

func readInput() (*Stone, error) {
	file, err := os.Open("test.txt")
	if err != nil {
		println("Failed to open input.txt: " + err.Error())
		return nil, err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var (
		firstStone   = &Stone{}
		currentStone = firstStone
		prevStone    *Stone
	)

	parts := strings.Split(string(content), " ")
	for _, part := range parts {
		partValue, err := strconv.Atoi(part)
		if err != nil {
			return nil, err
		}

		currentStone.Value = partValue
		currentStone.Next = &Stone{}
		currentStone.Prev = prevStone
		prevStone = currentStone
		currentStone = currentStone.Next
	}

	prevStone.Next = nil

	return firstStone, nil
}
