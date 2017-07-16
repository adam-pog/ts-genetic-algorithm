package main

import(
    "fmt"
)

const NumCities = 5
func main(){
    A := []int{}
    c := []int{}
    for i := 0; i < NumCities; i++ {
        A = append(A, i)
        c = append(c, 0)
    }
fmt.Println(A)
fmt.Println(c)
    //solution := A
    count := 1
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
            fmt.Println(A)
            count++
            c[i]++
            i = 0
        } else {
            c[i] = 0
            i ++
        }

    }
    fmt.Println(count)

}
