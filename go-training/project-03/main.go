package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8081", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Fatal(err)
	}
	if v, ok := r.Form["solve"]; ok {
		fmt.Fprintf(w, "%s %v\n", v[0], ok)
		var a1, b1, c1, d1, a2, b2, c2, d2, a3, b3, c3, d3 int
		var x, y, z float32
		if n, _ := fmt.Sscanf(v[0], "%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d", &a1, &b1, &c1, &d1, &a2, &b2, &c2, &d2, &a3, &b3, &c3, &d3); n == 12 {
			fmt.Fprintf(w, "system:\n%dx + %dy + %dz = %d\n%dx + %dy + %dz = %d\n%dx + %dy + %dz = %d\n", a1, b1, c1, d1, a2, b2, c2, d2, a3, b3, c3, d3)

			dx := calculateDeterminant(d1, b1, c1, d2, b2, c2, d3, b3, c3)
			dy := calculateDeterminant(a1, d1, c1, a2, d2, c2, a3, d3, c3)
			dz := calculateDeterminant(a1, b1, d1, a2, b2, d2, a3, b3, d3)
			d := calculateDeterminant(a1, b1, c1, a2, b2, c2, a3, b3, c3)

			if d != 0 {
				x = dx / d
				y = dy / d
				z = dz / d
				fmt.Fprintf(w, "solution:\nx = %.2f, y = %.2f, z = %.2f\n", x, y, z)
			} else {
				if dx == 0 || dy == 0 || dz == 0 {
					fmt.Fprintf(w, "%v\n", "solution:\nInconsistent - Indefinetly many solutions")
				} else {
					fmt.Fprintf(w, "%v\n", "solution:\nInconsistent - no solution")
				}
			}

		} else {
			fmt.Fprintf(w, "%v\n", "insufficient coefficients")
		}
	}
}

func calculateDeterminant(a1, b1, c1, a2, b2, c2, a3, b3, c3 int) float32 {
	var x, y, z float32
	x = float32((b2 * c3) - (b3 * c2))
	y = float32((a2 * c3) - (a3 * c2))
	z = float32((a2 * b3) - (a3 * b2))

	return (float32(a1) * x) - (float32(b1) * y) + (float32(c1) * z)
}
