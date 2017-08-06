package bruteforce

func FindExactSolution(cityMap []map[int]float64, numCities int) float64 {
    A := []int{}
    c := []int{}
    for i := 0; i < numCities; i++ {
        A = append(A, i)
        c = append(c, 0)
    }

    solution := calculateTourFitness(A, cityMap, numCities)
    count := 0
    for i := 0; i < numCities; {
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
            count++

            // if(count % 1000000 == 0) {
            //     fmt.Println("\nIteration:", count)
            // }

            newSol := calculateTourFitness(A, cityMap, numCities)
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


func calculateTourFitness(tour []int, cityMap []map[int]float64, numCities int) float64{
    fitness := 0.0

    for j := 0; j < numCities; j++ {
        pointA := tour[j]
        pointB := tour[(j+1) % numCities]

        fitness += cityMap[pointA][pointB]
    }

    return fitness
}
