package main

import (
    "fmt"
    "math/rand"
    "time"
    "sort"
)

const PopSize = 5
const NumCities = 10

const TournamentK = 2
const CrossoverRate = 0.6
const mutationRate = 0.08

type Population struct {
    tours []Tour
}

type Tour struct {
    path []int
    fitness int
}

func main() {
    rand.Seed(time.Now().UTC().UnixNano())

    fmt.Println("Hello, World!")
    cityMap := generateCityDists()
    currentPop := generatePop()
    currentPop.calculateFitness(cityMap)
    fittestTour := findFittest(currentPop.tours)


    for genNum := 0; genNum < 500000; genNum++ {
        currentPop = currentPop.crossover(fittestTour)
        currentPop.mutate()
        currentPop.calculateFitness(cityMap)
        fittestTour = findFittest(currentPop.tours)

        if genNum % 10000 == 0 {
            fmt.Println(fittestTour.fitness)
        }

    }

    //fmt.Println(fitPop)
    currentPop.calculateFitness(cityMap)

    fmt.Println("Fittest: ", findFittest(currentPop.tours).fitness)
    fmt.Println("Best Possible: ", findSolution(cityMap))
}








func findSolution(cityMap) int {
    A := []int{}
    c := []int{}
    for i := 0; i < NumCities; i++ {
        A.path = append(A.path, i)
        c.path = append(c.path, 0)
    }

    solution := A.fitness

    for i := 0; i < NumCities; i++ {
        if c[i] < i {
            if i % 2 == 0 {
                temp := A[0]
                A[0] = A[i]
                A[i] = temp
            } else {
                temp := c[i]
                c[i] = A[i]
                A[i] = temp
            }
            fmt.Println(A)
            c[i]++
            i = 0
        } else {
            c[i] = 0
            i ++
        }

    }



}






func (population Population)mutate() {
    for i := 0; i < PopSize; i++ {
        population.tours[i].mutateTour()
    }
}


func (tour Tour) mutateTour() {
    for i := 0; i < NumCities; i++ {
        if mutationRate > rand.Float64() {
            mutateIndex := rand.Intn(NumCities)
            city1 := tour.path[i]
            tour.path[i] = tour.path[mutateIndex]
            tour.path[mutateIndex] = city1
        }
    }
}

func (population Population) crossover(fittestTour Tour)(nextPop Population) {
    //tournament select
    nextPop.tours = append(nextPop.tours, fittestTour)
    for i := 0; i < PopSize -1; i++ {
        if CrossoverRate > rand.Float64() {
            parent1 := selectParent(population)
            parent2 := selectParent(population)

            nextPop.tours = append(nextPop.tours, crossoverTours(parent1, parent2))
        } else {
            //fmt.Println("not")
            nextPop.tours = append(nextPop.tours, population.tours[i])
        }
    }

    return
}

func crossoverTours(parent1 Tour, parent2 Tour) (newTour Tour) {
    points := []int{rand.Intn(NumCities), rand.Intn(NumCities + 1)}
    sort.Ints(points)
    pointA := points[0]
    pointB := points[1]
    parentSlice := parent1.path[pointA:pointB]

    count := 0
    for i := 0; count < pointA; i++{
        if !contains(parentSlice, parent2.path[i]) {
            newTour.path = append(newTour.path, parent2.path[i])
            count++
        }
    }

    for i := 0; i < len(parentSlice); i++{
        newTour.path = append(newTour.path, parentSlice[i])
    }

    count = NumCities
    for i := 0; count > pointB; i++{
        if !contains(parentSlice, parent2.path[i]) && !contains(newTour.path, parent2.path[i]) {
            newTour.path = append(newTour.path, parent2.path[i])
            count--
        }
    }

    return
}

func selectParent(population Population) (tour Tour){
    parentTours := []Tour{}
    for i := 0; i < TournamentK; i++ {
        parentTours = append(parentTours, population.tours[rand.Intn(PopSize)])
    }

    return findFittest(parentTours)
}

func findFittest(tours []Tour)(fittestTour Tour) {
    fittestTour = tours[0]
    for i := 1; i < TournamentK; i++ {
        if tours[i].fitness < fittestTour.fitness {
            fittestTour = tours[i]
        }
    }

    return
}

func (population Population) calculateFitness(cityMap []map[int]int){
    for i := 0; i < PopSize; i++ {
        tour := population.tours[i]
        if tour.fitness == 0 {
            length := len(tour.path)

            for j := 0; j < length; j++ {
                pointA := tour.path[j]
                pointB := tour.path[(j+1) % length]

                tour.fitness += cityMap[pointA][pointB]
            }

            population.tours[i] = tour
        }
    }
    return
}

func generatePop() (population Population){
    for i := 0; i < PopSize; i++ {
        population.tours = append(population.tours, Tour{rand.Perm(NumCities), 0})
    }

    return
}

func generateCityDists() (city_map []map[int]int) {
    for i := 0; i < NumCities; i++ {
        city_map = append(city_map, make(map[int]int))

        for j := 0; j < NumCities; j++ {
            if j > i {
                city_map[i][j] = rand.Intn(1000)
            } else if j < i {
                city_map[i][j] = city_map[j][i]
            }
        }

    }

    return
}

func random(min, max int) float64 {
    return float64(rand.Intn(max - min) + min)
}

func contains(s []int, e int) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}
