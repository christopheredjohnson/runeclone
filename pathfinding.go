package main

import (
	"container/heap"
	"math"
)

type Node struct {
	Point
	G, H   float64
	F      float64
	Parent *Node
	Index  int
}

type PriorityQueue []*Node

func (pq PriorityQueue) Len() int           { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].F < pq[j].F }
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}
func (pq *PriorityQueue) Push(x any) {
	n := x.(*Node)
	n.Index = len(*pq)
	*pq = append(*pq, n)
}
func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	node := old[n-1]
	*pq = old[:n-1]
	return node
}

func heuristic(a, b Point) float64 {
	return math.Abs(float64(a.X-b.X)) + math.Abs(float64(a.Y-b.Y))
}

func neighbors(p Point) []Point {
	return []Point{
		{p.X + 1, p.Y}, {p.X - 1, p.Y},
		{p.X, p.Y + 1}, {p.X, p.Y - 1},
	}
}

func FindPath(start, goal Point, m *Map) []Point {
	open := make(PriorityQueue, 0)
	heap.Init(&open)

	startNode := &Node{Point: start, G: 0, H: heuristic(start, goal)}
	startNode.F = startNode.H
	heap.Push(&open, startNode)

	costSoFar := map[Point]float64{start: 0}
	visited := map[Point]bool{}

	for open.Len() > 0 {
		current := heap.Pop(&open).(*Node)

		if current.Point == goal {
			return reconstructPath(current)
		}

		visited[current.Point] = true

		for _, next := range neighbors(current.Point) {
			if m.GetTile(next.X, next.Y) == nil || !m.GetTile(next.X, next.Y).IsWalkable() {
				continue
			}
			if visited[next] {
				continue
			}

			newCost := costSoFar[current.Point] + 1
			if oldCost, ok := costSoFar[next]; !ok || newCost < oldCost {
				costSoFar[next] = newCost
				h := heuristic(next, goal)
				node := &Node{
					Point:  next,
					G:      newCost,
					H:      h,
					F:      newCost + h,
					Parent: current,
				}
				heap.Push(&open, node)
			}
		}
	}
	return nil
}

func reconstructPath(end *Node) []Point {
	var path []Point
	for node := end; node != nil; node = node.Parent {
		path = append([]Point{node.Point}, path...)
	}
	return path
}
