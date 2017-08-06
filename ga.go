package main

import (
    "fmt"
    "math/rand"
    "time"
    "sort"
    "math"
    "./bruteforce"
)

const PopSize = 200
const NumCities = 10

const TournamentK = 20
const GenLimit = 500

const CrossoverRate = 0.6
const mutationRate = 0.03


type Tour struct {
    path []int
    fitness float64
}

type Coord struct {
    x, y float64
}

func main() {
    rand.Seed(time.Now().UTC().UnixNano())

    fmt.Println("Generating City Coords...")
    cityCoords := generateCityCoords()
    fmt.Println("Generating City Map...")
    cityMap := generateCityMap(cityCoords)


    fmt.Println("Initializing Population...")
    currentPop := generatePop()
    calculateFitness(currentPop, cityMap)
    fittestTour := findFittest(currentPop)

    currentFitness := 0.0
    fmt.Println("Starting iteration...\n")
    fitnessCounter := 0
    for converged(currentPop) && fitnessCounter < GenLimit{
        currentPop = crossover(currentPop, fittestTour)
        mutate(currentPop)
        calculateFitness(currentPop, cityMap)
        fittestTour = findFittest(currentPop)
        newFitness := fittestTour.fitness

        if newFitness != currentFitness {
            fmt.Println(fittestTour.fitness)
            currentFitness = newFitness
        } else{
            fitnessCounter++
        }
    }

    calculateFitness(currentPop, cityMap)

    fittest := findFittest(currentPop)
    fmt.Println("\nPath: ", fittest.path)
    fmt.Printf("\n\nFittest:       %f\n", fittest.fitness)
    fmt.Printf("Best Possible: %f\n", bruteforce.FindExactSolution(cityMap, NumCities))
}



func converged(population []Tour) (halt bool){
    geneCount := []int{}

    for i := 0; i < NumCities; i++ {
        geneCount = append(geneCount, 0)
    }

    for cityNum := 0; cityNum < NumCities; cityNum++ {
        for currentTour := 0; currentTour < PopSize; currentTour++ {
            sameCity := 1
            for insideTour := 0; insideTour < PopSize; insideTour++ {
                if insideTour != currentTour {
                    if population[currentTour].path[cityNum] == population[insideTour].path[cityNum] {
                        sameCity += 1
                    }
                }
            }

            if geneCount[cityNum] < sameCity {
                geneCount[cityNum] = sameCity
            }
        }
    }

    for i := 0; i < NumCities; i++ {
        if (float64(geneCount[i]) / float64(PopSize)) < .90 {
            halt = true
        }
    }
    if(halt == false) {
        fmt.Println("Converged!!!!")
    }
    return
}

func mutate(population []Tour) {
    for i := 0; i < PopSize; i++ {
        population[i].mutateTour()
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

func crossover(population []Tour, fittestTour Tour)(nextPop []Tour) {
    //tournament select
    //nextPop.tours = append(nextPop.tours, fittestTour)
    for i := 0; i < PopSize; i++ {
        if CrossoverRate > rand.Float64() {
            parent1 := selectParent(population)
            parent2 := selectParent(population)

            nextPop = append(nextPop, crossoverTours(parent1, parent2))
        } else {
            //fmt.Println("not")
            nextPop = append(nextPop, population[i])
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

func selectParent(population []Tour) (tour Tour){
    parentTours := []Tour{}
    for i := 0; i < TournamentK; i++ {
        parentTours = append(parentTours, population[rand.Intn(PopSize)])
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

func calculateFitness(population []Tour, cityMap []map[int]float64){
    for i := 0; i < PopSize; i++ {
        tour := population[i]
        if tour.fitness == 0.0 {
            length := len(tour.path)

            for j := 0; j < length; j++ {
                pointA := tour.path[j]
                pointB := tour.path[(j+1) % length]

                tour.fitness += cityMap[pointA][pointB]
            }

            population[i] = tour
        }
    }
    return
}

func generatePop() (population []Tour){
    for i := 0; i < PopSize; i++ {
        population = append(population, Tour{rand.Perm(NumCities), 0})
    }

    return
}

// func generateCityMap() (city_map []map[int]int) {
//     // for i := 0; i < NumCities; i++ {
//     //     city_map = append(city_map, make(map[int]int))
//     //
//     //     for j := 0; j < NumCities; j++ {
//     //         if j > i {
//     //             city_map[i][j] = rand.Intn(1000)
//     //         } else if j < i {
//     //             city_map[i][j] = city_map[j][i]
//     //         }
//     //     }
//     //
//     // }
//
//     city_map = []map[int]int{map[int]int{7:82, 9:852, 2:387, 3:477, 4:899, 5:646, 1:120, 6:883, 8:165},
//      map[int]int{2:307, 3:603, 6:391, 7:31, 8:109, 0:120, 4:710, 5:99, 9:418},
//      map[int]int{8:887, 9:194, 1:307, 3:165, 4:16, 6:255, 0:387, 5:720, 7:905},
//      map[int]int{2:165, 6:153, 7:879, 0:477, 1:603, 4:9, 5:478, 8:821, 9:22},
//      map[int]int{5:376, 6:256, 7:692, 8:532, 2:16, 3:9, 9:793, 0:899, 1:710},
//      map[int]int{1:99, 2:720, 7:481, 8:921, 0:646, 3:478, 4:376, 6:366, 9:487},
//      map[int]int{0:883, 3:153, 8:672, 9:905, 1:391, 2:255, 4:256, 5:366, 7:909},
//      map[int]int{2:905, 3:879, 4:692, 5:481, 9:706, 0:82, 1:31, 6:909, 8:585},
//      map[int]int{1:109, 3:821, 4:532, 6:672, 7:585, 9:532, 0:165, 5:921, 2:887},
//      map[int]int{7:706, 8:532, 0:852, 1:418, 2:194, 4:793, 6:905, 3:22, 5:487}}
//
//     return
// }

func generateCityMap(cityCoords []Coord) (city_map []map[int]float64) {
    for i := 0; i < NumCities; i++ {
        city_map = append(city_map, make(map[int]float64))

        for j := 0; j < NumCities; j++ {
            if j > i {
                coord1 := cityCoords[i]
                coord2 := cityCoords[j]
                distance := math.Sqrt( math.Pow((coord1.x - coord2.x), 2) + math.Pow((coord1.y - coord2.y), 2) )
                city_map[i][j] = distance
            } else if j < i {
                city_map[i][j] = city_map[j][i]
            }
        }

    }

    // city_map = []map[int]int{map[int]int{7:82, 9:852, 2:387, 3:477, 4:899, 5:646, 1:120, 6:883, 8:165},
    //  map[int]int{2:307, 3:603, 6:391, 7:31, 8:109, 0:120, 4:710, 5:99, 9:418},
    //  map[int]int{8:887, 9:194, 1:307, 3:165, 4:16, 6:255, 0:387, 5:720, 7:905},
    //  map[int]int{2:165, 6:153, 7:879, 0:477, 1:603, 4:9, 5:478, 8:821, 9:22},
    //  map[int]int{5:376, 6:256, 7:692, 8:532, 2:16, 3:9, 9:793, 0:899, 1:710},
    //  map[int]int{1:99, 2:720, 7:481, 8:921, 0:646, 3:478, 4:376, 6:366, 9:487},
    //  map[int]int{0:883, 3:153, 8:672, 9:905, 1:391, 2:255, 4:256, 5:366, 7:909},
    //  map[int]int{2:905, 3:879, 4:692, 5:481, 9:706, 0:82, 1:31, 6:909, 8:585},
    //  map[int]int{1:109, 3:821, 4:532, 6:672, 7:585, 9:532, 0:165, 5:921, 2:887},
    //  map[int]int{7:706, 8:532, 0:852, 1:418, 2:194, 4:793, 6:905, 3:22, 5:487}}

    return
}

func generateCityCoords()(cityCoords []Coord) {
    for i := 0; i < NumCities; i++ {
        cityCoords = append(cityCoords, Coord{float64(rand.Intn(1000)), float64(rand.Intn(1000))})
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
