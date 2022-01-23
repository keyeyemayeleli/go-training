package main

import (
	"fmt"
)

func main() {
	var n, c int
	fmt.Scan(&n)
	for i := 0; i < n; i++ {
		fmt.Scan(&c)
		x := make([]int, c)
		var sum int = 0
		for j := 0; j < c; j++ {
			fmt.Scan(&x[j])
			sum = sum + x[j]
		}
		var avg float32 = float32(sum) / float32(c)
		var passers int = 0
		for k := 0; k < c; k++ {
			if float32(x[k]) > avg {
				passers++
			}
		}

		var passpct = float32(passers) / float32(c) * 100
		fmt.Printf("%.3f%% \n", passpct)
		// fmt.Print(avg, passers, c, passpct)
	}
}
