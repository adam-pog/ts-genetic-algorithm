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
    fmt.Println(cityMap)
    currentPop := generatePop()
    currentPop.calculateFitness(cityMap)
    fittestTour := findFittest(currentPop.tours)


    for genNum := 0; genNum < 500000; genNum++ {
        currentPop = currentPop.crossover(fittestTour)
        currentPop.mutate()
        currentPop.calculateFitness(cityMap)
        fittestTour = findFittest(currentPop.tours)

        if genNum % 50000 == 0 {
            fmt.Println(fittestTour.fitness)
        }

    }

    //fmt.Println(fitPop)
    currentPop.calculateFitness(cityMap)

    fmt.Println("Fittest: ", findFittest(currentPop.tours).fitness)
    fmt.Println("Best Possible: ", findSolution(cityMap))
}








func findSolution(cityMap []map[int]int) int {
    A := []int{}
    c := []int{}
    for i := 0; i < NumCities; i++ {
        A = append(A, i)
        c = append(c, 0)
    }

    solution := calculateTourFitness(A, cityMap)

    for i := 0; i < NumCities; {
        if c[i] < i {
            if i % 2 == 0 {
                temp := A[0]
                A[0] = A[i]
                A[i] = temp
            } else {
                temp := A[c[i]]
                A[c[i]] = A[i]
                A[i] = temp
            }
            newSol := calculateTourFitness(A, cityMap)
            if newSol < solution {
                solution = newSol
            }
            c[i]++
            i = 0
        } else {
            c[i] = 0
            i ++
        }

    }

    return solution
}


func calculateTourFitness(tour []int, cityMap []map[int]int) int{
    fitness := 0

    for j := 0; j < NumCities; j++ {
        pointA := tour[j]
        pointB := tour[(j+1) % NumCities]

        fitness += cityMap[pointA][pointB]
    }

    return fitness
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
    // for i := 0; i < NumCities; i++ {
    //     city_map = append(city_map, make(map[int]int))
    //
    //     for j := 0; j < NumCities; j++ {
    //         if j > i {
    //             city_map[i][j] = rand.Intn(1000)
    //         } else if j < i {
    //             city_map[i][j] = city_map[j][i]
    //         }
    //     }
    //
    // }

    city_map = []map[int]int{map[int]int{7:82, 9:852, 2:387, 3:477, 4:899, 5:646, 1:120, 6:883, 8:165},
     map[int]int{2:307, 3:603, 6:391, 7:31, 8:109, 0:120, 4:710, 5:99, 9:418},
     map[int]int{8:887, 9:194, 1:307, 3:165, 4:16, 6:255, 0:387, 5:720, 7:905},
     map[int]int{2:165, 6:153, 7:879, 0:477, 1:603, 4:9, 5:478, 8:821, 9:22},
     map[int]int{5:376, 6:256, 7:692, 8:532, 2:16, 3:9, 9:793, 0:899, 1:710},
     map[int]int{1:99, 2:720, 7:481, 8:921, 0:646, 3:478, 4:376, 6:366, 9:487},
     map[int]int{0:883, 3:153, 8:672, 9:905, 1:391, 2:255, 4:256, 5:366, 7:909},
     map[int]int{2:905, 3:879, 4:692, 5:481, 9:706, 0:82, 1:31, 6:909, 8:585},
     map[int]int{1:109, 3:821, 4:532, 6:672, 7:585, 9:532, 0:165, 5:921, 2:887},
     map[int]int{7:706, 8:532, 0:852, 1:418, 2:194, 4:793, 6:905, 3:22, 5:487}}

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
