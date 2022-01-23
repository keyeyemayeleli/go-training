package main

import "fmt"

func main() {
	var n, a, b, c int
	fmt.Scan(&n)
	for i := 0; i < n; i++ {
		fmt.Scan(&a, &b, &c)
		if (a+b <= c) || (a+c <= b) || (b+c <= a) {
			println("Case ", i+1, ": Invalid")
		} else if (a == b) && (b == c) {
			println("Case ", i+1, ": Equilateral")
		} else if (a == b) || (b == c) {
			println("Case ", i+1, ": Isosceles")
		} else {
			println("Case ", i+1, ": Scalene")
		}
	}
}
