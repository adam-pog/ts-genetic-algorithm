package main

import (
    "fmt"
    "math/rand"
    "time"
    "sort"
    "math"
    "./bruteforce"
    "./drawmap"
    . "./structs"
    . "./config"
)


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
        newFitness := fittestTour.Fitness

        if newFitness != currentFitness {
            fmt.Println(fittestTour.Fitness)
            currentFitness = newFitness
        } else{
            fitnessCounter++
        }
    }

    calculateFitness(currentPop, cityMap)

    fittest := findFittest(currentPop)

    drawmap.DrawMap(fittest.Path, cityCoords)
    fmt.Println("\nPath: ", fittest.Path)
    fmt.Printf("\n\nFittest:       %f\n", fittest.Fitness)
    fmt.Printf("Best Possible: %f\n", bruteforce.FindExactSolution(cityMap))
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
                    if population[currentTour].Path[cityNum] == population[insideTour].Path[cityNum] {
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
        population[i].MutateTour()
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
    parentSlice := parent1.Path[pointA:pointB]

    count := 0
    for i := 0; count < pointA; i++{
        if !contains(parentSlice, parent2.Path[i]) {
            newTour.Path = append(newTour.Path, parent2.Path[i])
            count++
        }
    }

    for i := 0; i < len(parentSlice); i++{
        newTour.Path = append(newTour.Path, parentSlice[i])
    }

    count = NumCities
    for i := 0; count > pointB; i++{
        if !contains(parentSlice, parent2.Path[i]) && !contains(newTour.Path, parent2.Path[i]) {
            newTour.Path = append(newTour.Path, parent2.Path[i])
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
        if tours[i].Fitness < fittestTour.Fitness {
            fittestTour = tours[i]
        }
    }

    return
}

func calculateFitness(population []Tour, cityMap []map[int]float64){
    for i := 0; i < PopSize; i++ {
        tour := population[i]
        if tour.Fitness == 0.0 {
            length := len(tour.Path)

            for j := 0; j < length; j++ {
                pointA := tour.Path[j]
                pointB := tour.Path[(j+1) % length]

                tour.Fitness += cityMap[pointA][pointB]
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
                distance := math.Sqrt( math.Pow((coord1.X - coord2.X), 2) + math.Pow((coord1.Y - coord2.Y), 2) )
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
        cityCoords = append(cityCoords, Coord{float64(rand.Intn(1500)), float64(rand.Intn(1500))})
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
