package main

import (
	"fmt"
	"sort"
)

func main() {
	var n int
	x := make([]int, 3)
	fmt.Scan(&n)
	for i := 0; i < n; i++ {
		fmt.Scan(&x[0], &x[1], &x[2])
		sort.Ints(x)
		fmt.Printf("Case %v: %v \n", i+1, x[1])
	}
}
