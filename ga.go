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

    cityCoords := generateCityCoords()
    cityMap := generateCityMap(cityCoords)
    currentPop := generatePop()

    calculateFitness(currentPop, cityMap)
    fittestTour := findFittest(currentPop)

    currentFitness := 0.0
    fmt.Println("Starting Iteration...\n")
    fitnessCounter := 0
    totalGenCount := 0
    for converged(currentPop, totalGenCount) && fitnessCounter < GenLimit{
        currentPop = crossover(currentPop, fittestTour)
        mutate(currentPop)
        calculateFitness(currentPop, cityMap)
        fittestTour = findFittest(currentPop)
        newFitness := fittestTour.Fitness


        if(totalGenCount % GraphGenCount == 0) {
          fmt.Println("Graph#", totalGenCount, fittestTour.Fitness)
          drawmap.DrawMap(fittestTour.Path, cityCoords, totalGenCount)
        }

        if newFitness != currentFitness {
            fmt.Println("Gen", totalGenCount, "-", fittestTour.Fitness)
            currentFitness = newFitness
        } else{
            fitnessCounter++
        }
        totalGenCount++
    }

    calculateFitness(currentPop, cityMap)

    fittest := findFittest(currentPop)

    drawmap.DrawMap(fittest.Path, cityCoords, totalGenCount)
    fmt.Println("\nPath: ", fittest.Path)
    fmt.Printf("\n\nFittest:       %f\n", fittest.Fitness)

    fmt.Printf("Best Possible: %f\n", bruteforce.FindExactSolution(cityMap))
}



func converged(population []Tour, count int) (halt bool){
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
      convergence := float64(geneCount[i]) / float64(PopSize)
        if (convergence) < .90 {
            halt = true
        }
    }
    if(halt == false) {
        fmt.Println("--- Population Converged ---")
    }
    return
}

func mutate(population []Tour) {
    for i := 0; i < PopSize; i++ {
        population[i].MutateTour()
    }
}

func crossover(population []Tour, fittestTour Tour)(nextPop []Tour) {
    //tournamentK select
    //nextPop.tours = append(nextPop.tours, fittestTour)
    for i := 0; i < PopSize; i++ {
        if CrossoverRate > rand.Float64() {
            parent1 := selectParent(population)
            parent2 := selectParent(population)

            nextPop = append(nextPop, crossoverTours(parent1, parent2))
        } else {
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

    return
}

func generateCityCoords()(cityCoords []Coord) {
    for i := 0; i < NumCities; i++ {
        cityCoords = append(cityCoords, Coord{float64(rand.Intn(4000)), float64(rand.Intn(4000))})
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
