package main

import (
	"bufio"
	"os"
)

func main() {
	maze, err := readInput()
	if err != nil {
		return
	}

	start, end := Point{}, Point{}
	for y := 0; y < len(maze); y++ {
		for x := 0; x < len(maze[y]); x++ {
			if maze[y][x] == 'S' {
				start = Point{x, y}
				maze[y][x] = '.'
			} else if maze[y][x] == 'E' {
				end = Point{x, y}
				maze[y][x] = '.'
			}
		}
	}

	println(AStar(maze, start, end))
}

func readInput() ([][]rune, error) {
	file, err := os.Open("input.txt")
	if err != nil {
		println("Failed to open input.txt: " + err.Error())
		return nil, err
	}
	defer file.Close()

	var (
		maze [][]rune
	)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		maze = append(maze, []rune(line))
	}

	return maze, nil
}
