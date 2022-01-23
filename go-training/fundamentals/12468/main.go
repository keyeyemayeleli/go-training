package main

import (
	"fmt"
)

func main() {
	var a, b, c, counter int
	counter = 1
	for counter == 1 {
		fmt.Scan(&a, &b)
		if (a == -1) && (b == -1) {
			counter = 0
		} else {
			if a >= b {
				c = a - b
			} else {
				c = b - a
			}
			if c > 50 {
				c = 100 - c
			}
			fmt.Println(c)
		}
	}
}
