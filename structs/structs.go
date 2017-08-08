package structs

import (
    "math/rand"
    . "../config"
)


type Tour struct {
    Path []int
    Fitness float64
}

type Coord struct {
    X, Y float64
}

func (tour Tour) MutateTour() {
    for i := 0; i < NumCities; i++ {
        if MutationRate > rand.Float64() {
            mutateIndex := rand.Intn(NumCities)
            city1 := tour.Path[i]
            tour.Path[i] = tour.Path[mutateIndex]
            tour.Path[mutateIndex] = city1
        }
    }
}
