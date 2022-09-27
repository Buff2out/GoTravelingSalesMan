package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"strconv"
	"strings"
)

type Point struct {
	x float64
	y float64
}

func newPoint(X float64, Y float64) Point {
	return Point{
		x: X,
		y: Y,
	}
}

func fillPoints(fin []string, points []Point, amountPoints int) {
	for i := 0; i < amountPoints; i++ {
		valx, errx := strconv.ParseFloat(fin[2*i+1], 64)
		if errx != nil {
			panic(errx)
		}
		points[i].x = valx
		valy, erry := strconv.ParseFloat(fin[2*i+2], 64)
		if erry != nil {
			panic(erry)
		}
		points[i].y = valy
	}
}

func printMtrx[T any](graph [][]T, n int) {
	for i := 0; i < n; i++ {
		fmt.Println(graph[i])
	}
}

func createEmptyGraph(n int) [][]float64 {
	graph := make([][]float64, n)
	for i := 0; i < n; i++ {
		graph[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			graph[i][j] = -1
		}
	}
	for j := 0; j < n; j++ {
		graph[j][j] = 0
	}
	return graph
}

func createStartPopul(n int, m int) [][]int {
	popul := make([][]int, n)
	for i := 0; i < n; i++ {
		popul[i] = make([]int, m)
		for j := 0; j < m; j++ {
			popul[i][j] = j + 1
		}
	}
	for i := 0; i < n/2; i++ {
		rand.Shuffle(m, func(i1, j1 int) {
			popul[i][i1], popul[i][j1] = popul[i][j1], popul[i][i1]
		})
	}
	return popul
}

func fillDistsToGraph(graph [][]float64, points []Point, amountPoints int) {
	for i := 0; i < amountPoints; i++ {
		for j := i + 1; j < amountPoints; j++ {
			graph[j][i] = math.Sqrt(math.Pow((points[j].x-points[i].x), 2) + math.Pow((points[j].y-points[i].y), 2))
			graph[i][j] = graph[j][i]
		}
	}
}

func main() {
	finByte, err := ioutil.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	// else
	fin := strings.Fields(string(finByte))
	// string to int
	amountPoints, err := strconv.Atoi(fin[0])
	if err != nil {
		panic(err)
	}
	// else
	// initialising graph, popul, points, dists
	graph := createEmptyGraph(amountPoints)
	k := amountPoints % 2
	popul := createStartPopul(2*(amountPoints-k), amountPoints-1)

	points := make([]Point, amountPoints)
	fillPoints(fin, points, amountPoints)

	dists := make([]float64, 2*(amountPoints-k))
	fillDistsToGraph(graph, points, amountPoints)
	var summ int = 0
	var theBest float64 = -1
	for summ != 3 {

		if theBest == dists[0] {
			summ++
		} else {
			theBest = dists[0]
			summ = 0
		}
	}
	printMtrx(graph, amountPoints)
	printMtrx(popul, amountPoints)
	fmt.Println(points)
	fmt.Println(dists)

}
