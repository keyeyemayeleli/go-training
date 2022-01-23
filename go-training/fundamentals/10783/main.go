package main

import "fmt"

func main() {
	var a, b, c, n int
	fmt.Scan(&n)
	for j := 0; j < n; j++ {
		fmt.Scan(&a, &b)
		c = 0
		for i := a; i <= b; i++ {
			if i%2 != 0 {
				c = c + i
			}
		}
		fmt.Println(c)
	}
}
