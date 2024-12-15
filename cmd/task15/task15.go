package main

import (
	"io"
	"os"
	"strings"
)

type Map [][]int32

func (m Map) Copy() Map {
	var m2 = make(Map, len(m))
	for i := range m {
		m2[i] = make([]int32, len(m[i]))
		for j := range m[i] {
			m2[i][j] = m[i][j]
		}
	}
	return m2
}

func (m Map) Set(m2 Map) {
	for i := range m {
		for j := range m[i] {
			m[i][j] = m2[i][j]
		}
	}
}

type Moves string

type Position int32

func Pos(x, y int) Position {
	return Position((int32(x) << 16) | (int32(y) & 0xffff))
}

func (p Position) X() int {
	return int((p >> 16) & 0xffff)
}

func (p Position) Y() int {
	return int(p & 0xffff)
}

func (p Position) Move(direction int32) Position {
	switch direction {
	case '<':
		return Pos(p.X()-1, p.Y())
	case '>':
		return Pos(p.X()+1, p.Y())
	case '^':
		return Pos(p.X(), p.Y()-1)
	case 'v':
		return Pos(p.X(), p.Y()+1)
	}

	return p
}

func main() {
	m, moves, err := readInput()
	if err != nil {
		return
	}

	scaledUpMap := scaleUpMap(m)

	println(executeMoves(m, moves))
	println(executeMovesScaledUp(scaledUpMap, moves))
}

func executeMoves(m Map, moves Moves) int {
	robot := findRobot(m)
	m[robot.Y()][robot.X()] = '.'

	for _, move := range moves {
		target := robot.Move(move)

		// walk into wall
		if m[target.Y()][target.X()] == '#' {
			continue
		}

		// push boxes
		if m[target.Y()][target.X()] == 'O' {
			t := target
			for m[t.Y()][t.X()] == 'O' {
				t = t.Move(move)
			}

			if m[t.Y()][t.X()] == '.' {
				// we have a space to push
				m[target.Y()][target.X()] = '.'
				m[t.Y()][t.X()] = 'O'
			} else if m[t.Y()][t.X()] == '#' {
				// all boxes are against the wall
				continue
			}
		}

		robot = target
	}

	// calculate coordinates sum
	sum := 0
	for i := range m {
		for j := range m[i] {
			if m[i][j] == 'O' {
				sum += i*100 + j
			}
		}
	}

	return sum
}

func executeMovesScaledUp(m Map, moves Moves) int {
	robot := findRobot(m)
	m[robot.Y()][robot.X()] = '.'

	for _, move := range moves {
		isHorizontal := true
		if move == '^' || move == 'v' {
			isHorizontal = false
		}

		target := robot.Move(move)

		// walk into wall
		if m[target.Y()][target.X()] == '#' {
			continue
		}

		// push boxes
		if m[target.Y()][target.X()] == '[' || m[target.Y()][target.X()] == ']' {
			if isHorizontal {
				t := target
				for m[t.Y()][t.X()] == '[' || m[t.Y()][t.X()] == ']' {
					t = t.Move(move)
				}

				if m[t.Y()][t.X()] == '.' {
					// we have a space to push
					m[target.Y()][target.X()] = '.'
					t2 := target
					next := ']'
					for t2 != t {
						t2 = t2.Move(move)
						if m[t2.Y()][t2.X()] == '[' {
							m[t2.Y()][t2.X()] = ']'
							next = '['
						} else if m[t2.Y()][t2.X()] == ']' {
							m[t2.Y()][t2.X()] = '['
							next = ']'
						} else if m[t2.Y()][t2.X()] == '.' {
							m[t2.Y()][t2.X()] = next
						}
					}

				} else if m[t.Y()][t.X()] == '#' {
					// all boxes are against the wall
					continue
				}
			} else {
				if !moveVertical(m, target, move) {
					continue
				}
			}
		}

		robot = target
	}

	// calculate coordinates sum
	sum := 0
	for i := range m {
		for j := range m[i] {
			if m[i][j] == '[' {
				sum += i*100 + j
			}
		}
	}

	return sum
}

func moveVertical(m Map, position Position, direction int32) bool {
	var (
		left  Position
		right Position
	)
	if m[position.Y()][position.X()] == '[' {
		left = position
		right = position.Move('>')
	} else if m[position.Y()][position.X()] == ']' {
		left = position.Move('<')
		right = position
	} else {
		return false
	}

	leftMoved := left.Move(direction)
	rightMoved := right.Move(direction)
	if m[leftMoved.Y()][leftMoved.X()] == '.' && m[rightMoved.Y()][rightMoved.X()] == '.' {
		m[leftMoved.Y()][leftMoved.X()] = '['
		m[rightMoved.Y()][rightMoved.X()] = ']'
		m[left.Y()][left.X()] = '.'
		m[right.Y()][right.X()] = '.'
		return true
	} else if m[leftMoved.Y()][leftMoved.X()] == '#' || m[rightMoved.Y()][rightMoved.X()] == '#' {
		return false
	} else if m[leftMoved.Y()][leftMoved.X()] == '[' && m[rightMoved.Y()][rightMoved.X()] == ']' {
		if !moveVertical(m, leftMoved, direction) {
			return false
		}
		m[leftMoved.Y()][leftMoved.X()] = '['
		m[rightMoved.Y()][rightMoved.X()] = ']'
		m[left.Y()][left.X()] = '.'
		m[right.Y()][right.X()] = '.'
		return true
	}

	if m[leftMoved.Y()][leftMoved.X()] == ']' && m[rightMoved.Y()][rightMoved.X()] == '.' {
		if !moveVertical(m, leftMoved, direction) {
			return false
		}
		m[leftMoved.Y()][leftMoved.X()] = '['
		m[rightMoved.Y()][rightMoved.X()] = ']'
		m[left.Y()][left.X()] = '.'
		m[right.Y()][right.X()] = '.'
		return true
	} else if m[leftMoved.Y()][leftMoved.X()] == '.' && m[rightMoved.Y()][rightMoved.X()] == '[' {
		if !moveVertical(m, rightMoved, direction) {
			return false
		}
		m[leftMoved.Y()][leftMoved.X()] = '['
		m[rightMoved.Y()][rightMoved.X()] = ']'
		m[left.Y()][left.X()] = '.'
		m[right.Y()][right.X()] = '.'
		return true
	} else if m[leftMoved.Y()][leftMoved.X()] == ']' && m[rightMoved.Y()][rightMoved.X()] == '[' {
		m2 := m.Copy()
		if !moveVertical(m2, leftMoved, direction) {
			return false
		}
		if !moveVertical(m2, rightMoved, direction) {
			return false
		}
		m2[leftMoved.Y()][leftMoved.X()] = '['
		m2[rightMoved.Y()][rightMoved.X()] = ']'
		m2[left.Y()][left.X()] = '.'
		m2[right.Y()][right.X()] = '.'
		m.Set(m2)
		return true
	}

	return false
}

func findRobot(m Map) Position {
	for i, row := range m {
		for j, cell := range row {
			if cell == '@' {
				return Pos(j, i)
			}
		}
	}

	return 0
}

func scaleUpMap(m Map) Map {
	var m2 Map
	for _, row := range m {
		var row2 []int32

		for _, cell := range row {
			if cell == '.' {
				row2 = append(row2, '.')
				row2 = append(row2, '.')
			} else if cell == '#' {
				row2 = append(row2, '#')
				row2 = append(row2, '#')
			} else if cell == 'O' {
				row2 = append(row2, '[')
				row2 = append(row2, ']')
			} else if cell == '@' {
				row2 = append(row2, '@')
				row2 = append(row2, '.')
			}
		}

		m2 = append(m2, row2)
	}

	return m2
}

func readInput() (Map, Moves, error) {
	file, err := os.Open("input.txt")
	if err != nil {
		println("Failed to open input.txt: " + err.Error())
		return nil, "", err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, "", err
	}

	var (
		m Map
	)

	parts := strings.Split(string(content), "\n\n")
	for _, line := range strings.Split(parts[0], "\n") {
		var row []int32
		for _, cell := range line {
			row = append(row, cell)
		}

		m = append(m, row)
	}

	return m, Moves(parts[1]), nil
}
