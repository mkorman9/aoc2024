package main

import (
	"bufio"
	"os"
)

type Map [][]bool

func (m Map) Set(pos Position, v bool) Map {
	m2 := make(Map, len(m))
	for i := 0; i < len(m); i++ {
		m2[i] = make([]bool, len(m[i]))
		for j := 0; j < len(m[i]); j++ {
			m2[i][j] = m[i][j]
		}
	}

	m2[pos.Row()][pos.Column()] = v

	return m2
}

type Direction int8

const (
	DirectionUp Direction = iota
	DirectionDown
	DirectionRight
	DirectionLeft
)

func (d Direction) TurnRight() Direction {
	switch d {
	case DirectionUp:
		return DirectionRight
	case DirectionRight:
		return DirectionDown
	case DirectionDown:
		return DirectionLeft
	case DirectionLeft:
		return DirectionUp
	}

	return DirectionUp
}

type Facing int8

const (
	FacingEmpty Facing = iota
	FacingWall
	FacingExit
)

type Position int32

const InvalidPosition Position = -1

func Pos(row, column int16) Position {
	return Position((int32(row) << 16) | (int32(column) & 0xffff))
}

func (p Position) Row() int16 {
	return int16((p >> 16) & 0xffff)
}

func (p Position) Column() int16 {
	return int16(p & 0xffff)
}

func (p Position) Move(direction Direction, m Map) Position {
	var p2 Position
	switch direction {
	case DirectionUp:
		p2 = Pos(p.Row()-1, p.Column())
	case DirectionDown:
		p2 = Pos(p.Row()+1, p.Column())
	case DirectionRight:
		p2 = Pos(p.Row(), p.Column()+1)
	case DirectionLeft:
		p2 = Pos(p.Row(), p.Column()-1)
	}

	if p2.Row() < 0 || p2.Row() >= int16(len(m)) || p2.Column() < 0 || p2.Column() >= int16(len(m[p2.Row()])) {
		return InvalidPosition
	}

	return p2
}

func (p Position) GetFacing(direction Direction, m Map) Facing {
	p2 := p.Move(direction, m)
	if p2 == InvalidPosition {
		return FacingExit
	}

	if m[p2.Row()][p2.Column()] {
		return FacingWall
	}

	return FacingEmpty
}

func main() {
	m, guard, err := readInput()
	if err != nil {
		return
	}

	println(simulateAndCountVisited(guard, m))
	println(simulateObstacles(guard, m))
}

func simulateAndCountVisited(guard Position, m Map) int {
	visited := make(map[Position]struct{})
	guardDirection := DirectionUp

	for running := true; running; {
		visited[guard] = struct{}{}

		facing := guard.GetFacing(guardDirection, m)
		switch facing {
		case FacingExit:
			running = false
		case FacingWall:
			guardDirection = guardDirection.TurnRight()
		case FacingEmpty:
			guard = guard.Move(guardDirection, m)
		}
	}

	return len(visited)
}

func simulateObstacles(guard Position, m Map) int {
	cycles := 0

	for i := 0; i < len(m); i++ {
		for j := 0; j < len(m[i]); j++ {
			pos := Pos(int16(i), int16(j))
			if m[i][j] {
				continue
			}

			m2 := m.Set(pos, true)
			if hasCycle(guard, m2) {
				cycles++
			}
		}
	}

	return cycles
}

func hasCycle(guard Position, m Map) bool {
	visited := make(map[Position]Position)
	guardDirection := DirectionUp

	for running := true; running; {
		facing := guard.GetFacing(guardDirection, m)
		switch facing {
		case FacingExit:
			running = false
		case FacingWall:
			guardDirection = guardDirection.TurnRight()
		case FacingEmpty:
			moved := guard.Move(guardDirection, m)
			if from, ok := visited[moved]; ok && from == guard {
				return true
			}
			visited[moved] = guard
			guard = moved
		}
	}

	return false
}

func readInput() (Map, Position, error) {
	file, err := os.Open("input.txt")
	if err != nil {
		println("Failed to open input.txt: " + err.Error())
		return nil, 0, err
	}
	defer file.Close()

	var (
		m     Map
		guard Position
	)

	scanner := bufio.NewScanner(file)
	var i int16
	for scanner.Scan() {
		line := scanner.Text()

		var row []bool
		var j int16
		for _, r := range line {
			if r == '.' {
				row = append(row, false)
			} else if r == '^' {
				guard = Pos(i, j)
				row = append(row, false)
			} else {
				row = append(row, true)
			}

			j++
		}

		m = append(m, row)
		i++
	}

	return m, guard, nil
}
