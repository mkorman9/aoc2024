package main

import (
	"bufio"
	"os"
	"strconv"
)

type Map [][]int

type Direction int8

const (
	DirectionLeft Direction = iota
	DirectionRight
	DirectionUp
	DirectionDown
)

type Point int32

func NewPoint(x, y int16) Point {
	return Point((int32(x) << 16) | (int32(y) & 0xffff))
}

func (p Point) X() int16 {
	return int16((p >> 16) & 0xffff)
}

func (p Point) Y() int16 {
	return int16(p & 0xffff)
}

func (p Point) Move(d Direction) Point {
	switch d {
	case DirectionLeft:
		return NewPoint(p.X()-1, p.Y())
	case DirectionRight:
		return NewPoint(p.X()+1, p.Y())
	case DirectionUp:
		return NewPoint(p.X(), p.Y()-1)
	case DirectionDown:
		return NewPoint(p.X(), p.Y()+1)
	}

	return -1
}

func (p Point) CanMove(d Direction, m Map) bool {
	p2 := p.Move(d)
	if p2.Y() < 0 || p2.Y() >= int16(len(m)) || p2.X() < 0 || p2.X() >= int16(len(m[p2.Y()])) {
		return false
	}

	return m[p2.Y()][p2.X()]-m[p.Y()][p.X()] == 1
}

func main() {
	m, err := readInput()
	if err != nil {
		return
	}

	println(countTrailScores(m))
	println(countTrailRatings(m))
}

func countTrailScores(m Map) int {
	sum := 0
	trailheads := findTrailHeads(m)

	for _, trailhead := range trailheads {
		score := traverse(m, trailhead)
		sum += score
	}

	return sum
}

func countTrailRatings(m Map) int {
	sum := 0
	trailheads := findTrailHeads(m)
	trailends := findTrailEnds(m)

	for _, trailhead := range trailheads {
		for _, trailend := range trailends {
			rating := traverseTo(m, trailhead, trailend)
			sum += rating
		}
	}

	return sum
}

func traverse(m Map, start Point) int {
	score := 0
	stack := []Point{start}
	visited := make(map[Point]struct{})
	isVisited := func(p Point) bool {
		_, ok := visited[p]
		return ok
	}

	for len(stack) > 0 {
		p := stack[0]
		stack = stack[1:]

		if isVisited(p) {
			continue
		}

		if m[p.Y()][p.X()] == 9 {
			score++
		}

		if p.CanMove(DirectionLeft, m) {
			stack = append(stack, p.Move(DirectionLeft))
		}
		if p.CanMove(DirectionRight, m) {
			stack = append(stack, p.Move(DirectionRight))
		}
		if p.CanMove(DirectionUp, m) {
			stack = append(stack, p.Move(DirectionUp))
		}
		if p.CanMove(DirectionDown, m) {
			stack = append(stack, p.Move(DirectionDown))
		}

		visited[p] = struct{}{}
	}

	return score
}

func traverseTo(m Map, start, end Point) int {
	rating := 0
	stack := []Point{start}

	for len(stack) > 0 {
		p := stack[0]
		stack = stack[1:]

		if p == end {
			rating++
		}

		if p.CanMove(DirectionLeft, m) {
			stack = append(stack, p.Move(DirectionLeft))
		}
		if p.CanMove(DirectionRight, m) {
			stack = append(stack, p.Move(DirectionRight))
		}
		if p.CanMove(DirectionUp, m) {
			stack = append(stack, p.Move(DirectionUp))
		}
		if p.CanMove(DirectionDown, m) {
			stack = append(stack, p.Move(DirectionDown))
		}
	}

	return rating
}

func findTrailHeads(m Map) []Point {
	var heads []Point
	for i := range m {
		for j := range m[i] {
			if m[i][j] == 0 {
				heads = append(heads, NewPoint(int16(j), int16(i)))
			}
		}
	}
	return heads
}

func findTrailEnds(m Map) []Point {
	var ends []Point
	for i := range m {
		for j := range m[i] {
			if m[i][j] == 9 {
				ends = append(ends, NewPoint(int16(j), int16(i)))
			}
		}
	}
	return ends
}

func readInput() (Map, error) {
	file, err := os.Open("input.txt")
	if err != nil {
		println("Failed to open input.txt: " + err.Error())
		return nil, err
	}
	defer file.Close()

	var (
		m Map
	)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		var row []int
		for _, c := range line {
			d, err := strconv.Atoi(string(c))
			if err != nil {
				return nil, err
			}
			row = append(row, d)
		}

		m = append(m, row)
	}

	return m, nil
}
