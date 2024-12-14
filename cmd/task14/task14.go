package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
)

const (
	MapWidth  = 101
	MapHeight = 103
)

type Robot struct {
	X         int
	Y         int
	VelocityX int
	VelocityY int
}

func main() {
	robots, err := readInput()
	if err != nil {
		return
	}

	//println(simulateMovement(robots, 100))
	println(findChristmasTree(robots))
}

func simulateMovement(robots []*Robot, seconds int) int {
	robots = copyRobots(robots)

	for range seconds {
		for _, robot := range robots {
			robot.X += robot.VelocityX
			robot.Y += robot.VelocityY
			for robot.X >= MapWidth {
				robot.X -= MapWidth
			}
			for robot.X < 0 {
				robot.X += MapWidth
			}
			for robot.Y >= MapHeight {
				robot.Y -= MapHeight
			}
			for robot.Y < 0 {
				robot.Y += MapHeight
			}
		}
	}

	var (
		quadrantLeftTop     = 0
		quadrantRightTop    = 0
		quadrantLeftBottom  = 0
		quadrantRightBottom = 0
	)

	for _, robot := range robots {
		if robot.X < (MapWidth/2) && robot.Y < (MapHeight/2) {
			quadrantLeftTop++
		}
		if robot.X > (MapWidth/2) && robot.Y < (MapHeight/2) {
			quadrantRightTop++
		}
		if robot.X < (MapWidth/2) && robot.Y > (MapHeight/2) {
			quadrantLeftBottom++
		}
		if robot.X > (MapWidth/2) && robot.Y > (MapHeight/2) {
			quadrantRightBottom++
		}
	}

	return quadrantLeftTop * quadrantRightTop * quadrantLeftBottom * quadrantRightBottom
}

func findChristmasTree(robots []*Robot) int {
	for second := 1; second < 7000; second++ {
		for _, robot := range robots {
			robot.X += robot.VelocityX
			robot.Y += robot.VelocityY
			for robot.X >= MapWidth {
				robot.X -= MapWidth
			}
			for robot.X < 0 {
				robot.X += MapWidth
			}
			for robot.Y >= MapHeight {
				robot.Y -= MapHeight
			}
			for robot.Y < 0 {
				robot.Y += MapHeight
			}
		}

		var (
			avgX = 0.0
			avgY = 0.0
		)
		for _, robot := range robots {
			avgX += float64(robot.X)
			avgY += float64(robot.Y)
		}
		avgX /= float64(len(robots))
		avgY /= float64(len(robots))

		var (
			deviationX = 0.0
			deviationY = 0.0
		)
		for _, robot := range robots {
			deviationX += (float64(robot.X) - avgX) * (float64(robot.X) - avgX)
			deviationY += (float64(robot.Y) - avgY) * (float64(robot.Y) - avgY)
		}
		deviationX /= float64(len(robots))
		deviationY /= float64(len(robots))
		deviationX = math.Sqrt(deviationX)
		deviationY = math.Sqrt(deviationY)

		fmt.Printf("%d\t%f\t%f\n", second, deviationX, deviationY)
	}

	return 0
}

func copyRobots(robots []*Robot) []*Robot {
	var robots2 = make([]*Robot, len(robots))
	for i, robot := range robots {
		robots2[i] = &Robot{robot.X, robot.Y, robot.VelocityX, robot.VelocityY}
	}
	return robots2
}

func readInput() ([]*Robot, error) {
	file, err := os.Open("input.txt")
	if err != nil {
		println("Failed to open input.txt: " + err.Error())
		return nil, err
	}
	defer file.Close()

	var (
		robots []*Robot
	)

	pattern := regexp.MustCompile("p=(-?\\d+),(-?\\d+) v=(-?\\d+),(-?\\d+)")

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		matches := pattern.FindAllStringSubmatch(line, -1)

		x, err := strconv.Atoi(matches[0][1])
		if err != nil {
			return nil, err
		}
		y, err := strconv.Atoi(matches[0][2])
		if err != nil {
			return nil, err
		}
		velocityX, err := strconv.Atoi(matches[0][3])
		if err != nil {
			return nil, err
		}
		velocityY, err := strconv.Atoi(matches[0][4])
		if err != nil {
			return nil, err
		}

		robots = append(robots, &Robot{x, y, velocityX, velocityY})
	}

	return robots, nil
}
