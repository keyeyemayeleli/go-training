package main

import (
	"fmt"
)

func main() {
	var a, b, c, counter int
	counter = 1
	for counter == 1 {
		fmt.Scan(&a, &b, &c)
		if (a == 0) && (b == 0) && (c == 0) {
			counter = 0
		} else if (a*a+b*b == c*c) || (c*c+b*b == a*a) || (a*a+c*c == b*b) {
			fmt.Println("right")
		} else {
			fmt.Println("wrong")
		}
	}
}
