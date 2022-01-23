package main

import "fmt"

func main() {
	var v, t int
	_, err := fmt.Scan(&v, &t)
	for err == nil {
		fmt.Println(2 * v * t)
		_, err = fmt.Scan(&v, &t)
	}
}
