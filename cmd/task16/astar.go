package main

import (
	"container/heap"
	"math"
)

type Point struct {
	x, y int
}

type Node struct {
	point     Point
	facing    int
	cost, est int
	prev      *Node
}

type PriorityQueue []*Node

func (pq PriorityQueue) Len() int            { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool  { return pq[i].cost+pq[i].est < pq[j].cost+pq[j].est }
func (pq PriorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) { *pq = append(*pq, x.(*Node)) }
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

var directions = []Point{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}

func Heuristic(p1, p2 Point) int {
	return int(math.Abs(float64(p1.x-p2.x)) + math.Abs(float64(p1.y-p2.y)))
}

func IsValid(grid [][]rune, p Point) bool {
	return p.y >= 0 && p.y < len(grid) && p.x >= 0 && p.x < len(grid[0]) && grid[p.y][p.x] == '.'
}

func AStar(grid [][]rune, start, end Point) int {
	visited := make(map[[3]int]bool)

	pq := &PriorityQueue{}
	heap.Init(pq)
	startNode := &Node{point: start, facing: 0, cost: 0, est: Heuristic(start, end)}
	heap.Push(pq, startNode)

	for pq.Len() > 0 {
		current := heap.Pop(pq).(*Node)
		state := [3]int{current.point.x, current.point.y, current.facing}

		if visited[state] {
			continue
		}
		visited[state] = true

		if current.point == end {
			return current.cost
		}

		for i, dir := range directions {
			newPoint := Point{current.point.x + dir.x, current.point.y + dir.y}
			if !IsValid(grid, newPoint) {
				continue
			}

			turnCost := 0
			if current.facing != i {
				turnCost = 1000
			}

			neighbor := &Node{
				point:  newPoint,
				facing: i,
				cost:   current.cost + 1 + turnCost,
				est:    Heuristic(newPoint, end),
				prev:   current,
			}

			neighborState := [3]int{neighbor.point.x, neighbor.point.y, neighbor.facing}
			if !visited[neighborState] {
				heap.Push(pq, neighbor)
			}
		}
	}

	return -1
}
