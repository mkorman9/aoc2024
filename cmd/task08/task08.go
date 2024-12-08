package main

import (
	"bufio"
	"os"
)

type Vec2D int32

func NewVec2D(x, y int16) Vec2D {
	return Vec2D((int32(x) << 16) | (int32(y) & 0xffff))
}

func (v Vec2D) X() int16 {
	return int16((v >> 16) & 0xffff)
}

func (v Vec2D) Y() int16 {
	return int16(v & 0xffff)
}

func (v Vec2D) Reflect(pivot Vec2D) Vec2D {
	return NewVec2D(2*pivot.X()-v.X(), 2*pivot.Y()-v.Y())
}

type Frequencies map[rune][]Vec2D

func main() {
	frequencies, boundary, err := readInput()
	if err != nil {
		return
	}

	println(countAnitnodes(frequencies, boundary))
	println(countHarmonicAnitnodes(frequencies, boundary))
}

func countAnitnodes(frequencies Frequencies, boundary Vec2D) int {
	antinodes := map[Vec2D]struct{}{}

	for freq := range frequencies {
		for i := range frequencies[freq] {
			for j := range frequencies[freq] {
				if i == j {
					continue
				}

				a := frequencies[freq][i]
				b := frequencies[freq][j]
				r := a.Reflect(b)

				if r.X() >= 0 && r.X() < boundary.X() && r.Y() >= 0 && r.Y() < boundary.Y() {
					antinodes[r] = struct{}{}
				}
			}
		}
	}

	return len(antinodes)
}

func countHarmonicAnitnodes(frequencies Frequencies, boundary Vec2D) int {
	antinodes := map[Vec2D]struct{}{}

	for freq := range frequencies {
		for i := range frequencies[freq] {
			for j := range frequencies[freq] {
				if i == j {
					continue
				}

				a := frequencies[freq][i]
				b := frequencies[freq][j]
				antinodes[a] = struct{}{}
				antinodes[b] = struct{}{}

				for {
					r := a.Reflect(b)
					a = b
					b = r

					if r.X() >= 0 && r.X() < boundary.X() && r.Y() >= 0 && r.Y() < boundary.Y() {
						antinodes[r] = struct{}{}
					} else {
						break
					}
				}
			}
		}
	}

	return len(antinodes)
}

func readInput() (Frequencies, Vec2D, error) {
	file, err := os.Open("input.txt")
	if err != nil {
		println("Failed to open input.txt: " + err.Error())
		return nil, 0, err
	}
	defer file.Close()

	var (
		frequencies = make(Frequencies)
		row, column = 0, 0
	)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		column = 0
		for _, c := range line {
			if c != '.' {
				if _, ok := frequencies[c]; !ok {
					frequencies[c] = []Vec2D{}
				}

				frequencies[c] = append(frequencies[c], NewVec2D(int16(column), int16(row)))
			}

			column++
		}

		row++
	}

	return frequencies, NewVec2D(int16(column), int16(row)), nil
}
