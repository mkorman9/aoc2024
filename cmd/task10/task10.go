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

type Point int64

func NewPoint(x, y int) Point {
	return Point((x << 32) | (y & 0xffffffff))
}

func (p Point) X() int {
	return int((p >> 32) & 0xffffffff)
}

func (p Point) Y() int {
	return int(p & 0xffffffff)
}

func (p Point) Move(d Direction, m Map) (Point, bool) {
	var p2 Point
	switch d {
	case DirectionLeft:
		p2 = NewPoint(p.X()-1, p.Y())
	case DirectionRight:
		p2 = NewPoint(p.X()+1, p.Y())
	case DirectionUp:
		p2 = NewPoint(p.X(), p.Y()-1)
	case DirectionDown:
		p2 = NewPoint(p.X(), p.Y()+1)
	}

	if p2.Y() < 0 || p2.Y() >= len(m) || p2.X() < 0 || p2.X() >= len(m[p2.Y()]) {
		return p2, false
	}
	if m[p2.Y()][p2.X()]-m[p.Y()][p.X()] != 1 {
		return p2, false
	}

	return p2, true
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

	for len(stack) > 0 {
		p := stack[0]
		stack = stack[1:]

		if _, ok := visited[p]; ok {
			continue
		}
		visited[p] = struct{}{}

		if m[p.Y()][p.X()] == 9 {
			score++
		}

		if p2, ok := p.Move(DirectionLeft, m); ok {
			stack = append(stack, p2)
		}
		if p2, ok := p.Move(DirectionRight, m); ok {
			stack = append(stack, p2)
		}
		if p2, ok := p.Move(DirectionUp, m); ok {
			stack = append(stack, p2)
		}
		if p2, ok := p.Move(DirectionDown, m); ok {
			stack = append(stack, p2)
		}
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

		if p2, ok := p.Move(DirectionLeft, m); ok {
			stack = append(stack, p2)
		}
		if p2, ok := p.Move(DirectionRight, m); ok {
			stack = append(stack, p2)
		}
		if p2, ok := p.Move(DirectionUp, m); ok {
			stack = append(stack, p2)
		}
		if p2, ok := p.Move(DirectionDown, m); ok {
			stack = append(stack, p2)
		}
	}

	return rating
}

func findTrailHeads(m Map) []Point {
	var heads []Point
	for i := range m {
		for j := range m[i] {
			if m[i][j] == 0 {
				heads = append(heads, NewPoint(j, i))
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
				ends = append(ends, NewPoint(j, i))
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
