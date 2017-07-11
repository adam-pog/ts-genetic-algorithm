package main

import (
    "fmt"
    "math/rand"
    "time"
)

const PopSize = 5
const NumCities = 5

type Population struct {
    tours []Tour
}

type Tour struct {
    path []int
}

func main() {
    rand.Seed(time.Now().UTC().UnixNano())

    fmt.Println("Hello, World!")
    city_map := generateCityDists()
    initial_pop := generatePop()

    fmt.Println(initial_pop)

    for i := 0; i < PopSize; i++ {
        fmt.Println(city_map[i])
    }

    //calculateFitness()
}

func generatePop() (population []Tour){
    for i := 0; i < PopSize; i++ {
        population = append(population, Tour{rand.Perm(NumCities)})
    }

    return
}

func generateCityDists() (city_map []map[int]int) {
    for i := 0; i < PopSize; i++ {
        city_map = append(city_map, make(map[int]int))

        for j := 0; j < PopSize; j++ {
            if j > i {
                city_map[i][j] = rand.Intn(100)
            } else if j < i {
                city_map[i][j] = city_map[j][i]
            }
        }

    }

    return
}
// func generate_city_dists() (city_map map[int][]int) {
//     city_map = make(map[int][]int)
//
//     for i := 0; i < PopSize; i++ {
//         city_map[i] = []int{}
//
//         for j := 0; j < PopSize; j++ {
//             if j == i {
//                 city_map[i] = append(city_map[i], 0)
//             } else if j > i {
//                 city_map[i] = append(city_map[i], rand.Intn(10))
//             } else if j < i {
//                 city_map[i] = append(city_map[i], city_map[j][i])
//             }
//         }
//
//     }
//
//     return
// }
