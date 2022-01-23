package main

import (
	"fmt"
	"sort"
)

func main() {
	var x_arr []int
	var x int
	_, err := fmt.Scan(&x)
	for err == nil {
		x_arr = append(x_arr, x)
		sort.Ints(x_arr)
		if (len(x_arr) % 2) == 0 {
			fmt.Println(((x_arr[len(x_arr)/2]) + (x_arr[(len(x_arr)/2)-1])) / 2)
		} else {
			fmt.Println((x_arr[len(x_arr)/2]))
		}
		_, err = fmt.Scan(&x)
	}
}
