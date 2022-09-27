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

func shuffleParents(popul [][]int, n int, m int) {
	rand.Shuffle(n, func(i1, i2 int) {
		popul[i1], popul[i2] = popul[i2], popul[i1]
	})
}

func crossOver(popul [][]int, amountPoints int, k int) {
	shuffleParents(popul, amountPoints-k, amountPoints-1)
	//скрещиваем перемешавшихся родителей между собой и порождаем потомков
	genBool1 := make([]bool, amountPoints-1) // на самом деле это хромосомный бул? состоящий из бул-генов
	genBool2 := make([]bool, amountPoints-1) // на самом деле это хромосомный бул? состоящий из бул-генов
	raNum := 0
	// заполняем булевый список "трушками"
	for j := 0; j < amountPoints-1; j++ {
		genBool1[j] = true
	}
	for j := 0; j < amountPoints-1; j++ {
		genBool2[j] = true
	}

	for i := 0; i < amountPoints-k; i = i + 2 {
		raNum = rand.Intn(amountPoints-1-1) + 1 //1 + rand() % (amountPoints - 1); // -1 (длина хромосомы)
		j1 := 0
		// первые raNum генов добавляем в потомков
		for j1 < raNum {
			popul[amountPoints-k+i][j1] = popul[i][j1]
			genBool1[popul[i][j1]-1] = false
			popul[amountPoints-k+i+1][j1] = popul[i+1][j1]
			genBool2[popul[i+1][j1]-1] = false
			j1++
		}
		// следующие [raNum; amountPoints - 1) генов добавляем в этих же потомков
		j2 := 0
		j3 := j1
		for j1 < amountPoints-1 {
			if genBool1[popul[i+1][j2]-1] {
				popul[amountPoints-k+i][j1] = popul[i+1][j2]
				j1++
			}
			j2++
		}
		// повторение всего того, но со вторым потомком
		j2 = 0
		j1 = j3
		for j1 < amountPoints-1 {
			if genBool2[popul[i][j2]-1] {
				popul[amountPoints-k+i+1][j1] = popul[i][j2]
				j1++
			}
			j2++
		}
		for j := 0; j < amountPoints-1; j++ {
			genBool1[j] = true
		}
		for j := 0; j < amountPoints-1; j++ {
			genBool2[j] = true
		}
	}
}

func toMutate(popul [][]int, amountPoints int, k int) {
	// расслабон, по сравнению со скрещиванием просто инвертируем последовательность элементов в случайном сгенеринном диапазоне
	a1, b1 := 0, 0
	for i := 0; i < 2*(amountPoints-k); i++ {
		if amountPoints-2 == 0 {
			a1 = 0
		} else {
			a1 = rand.Intn(amountPoints - 2) // rand() % (amountPoints - 2);
		}
		b1 = rand.Intn(amountPoints-1-a1) + a1 // a1 + rand() % (amountPoints - 1 - a1);
		for j := 0; j <= (b1-a1)/2; j++ {
			popul[i][a1+j], popul[i][b1-j] = popul[i][b1-j], popul[i][a1+j] // swapGens(i, a1 + j, b1 - j, popul);
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
		crossOver(popul, amountPoints, k)
		toMutate(popul, amountPoints, k)
		if theBest == dists[0] {
			summ++
		} else {
			theBest = dists[0]
			summ = 0
		}
	}
	printMtrx(graph, amountPoints)
	printMtrx(popul, 2*(amountPoints-k))
	fmt.Println(points)
	fmt.Println(dists)

}
