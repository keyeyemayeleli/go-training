package main

import (
	"fmt"
	"sort"
)

func main() {
	var n, c int
	fmt.Scan(&n)
	for i := 0; i < n; i++ {
		fmt.Scan(&c)
		x := make([]int, c)
		for j := 0; j < c; j++ {
			fmt.Scan(&x[j])
		}
		sort.Ints(x)
		fmt.Printf("Case %v: %v \n", i+1, x[c-1])
	}
}
