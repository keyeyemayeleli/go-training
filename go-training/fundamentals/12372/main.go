package main

import "fmt"

func main() {
	var n, l, w, h int
	fmt.Scan(&n)
	for i := 0; i < n; i++ {
		fmt.Scan(&l, &w, &h)
		if (h <= 20) && (w <= 20) && (l <= 20) {
			fmt.Printf("Case %v: good\n", i+1)
		} else {
			fmt.Printf("Case %v: bad\n", i+1)
		}
	}
}
