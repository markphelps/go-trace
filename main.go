package main

import (
	"fmt"
	"os"
)

func check(e error, s string) {
	if e != nil {
		fmt.Fprintf(os.Stderr, s, e)
		os.Exit(1)
	}
}

func main() {
	// size of image x and y
	nx := 400
	ny := 300

	const color = 255.99

	f, err := os.Create("out.ppm")

	defer f.Close()

	check(err, "Error opening file: %v\n")

	// http://netpbm.sourceforge.net/doc/ppm.html
	_, err = fmt.Fprintf(f, "P3\n%d %d\n255\n", nx, ny)

	check(err, "Error writting to file: %v\n")

	// writes each pixel with r/g/b values
	// from top left to bottom right
	for j := ny - 1; j >= 0; j-- {
		for i := 0; i < nx; i++ {
			// red and green values range from
			// 0.0 to 1.0
			r := float64(i) / float64(nx)
			g := float64(j) / float64(ny)
			b := 0.2

			// get intensity of colors
			ir := int(color * r)
			ig := int(color * g)
			ib := int(color * b)

			_, err = fmt.Fprintf(f, "%d %d %d\n", ir, ig, ib)

			check(err, "Error writting to file: %v\n")
		}
	}
}
