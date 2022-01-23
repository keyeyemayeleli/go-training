package main

import "fmt"

func main() {
	var a, b int
	_, err := fmt.Scan(&a, &b)
	for err == nil {
		println(b - a)
		_, err = fmt.Scan(&a, &b)
	}
}
