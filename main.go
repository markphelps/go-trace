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

	lowerLeft := Vector{-2.0, -1.0, -1.0}
	horizontal := Vector{4.0, 0.0, 0.0}
	vertical := Vector{0.0, 2.0, 0.0}
	origin := Vector{0.0, 0.0, 0.0}

	// writes each pixel with r/g/b values
	// from top left to bottom right
	for j := ny - 1; j >= 0; j-- {
		for i := 0; i < nx; i++ {
			u := float64(i) / float64(nx)
			v := float64(j) / float64(ny)

			hor := horizontal.MultiplyScalar(u)
			vert := vertical.MultiplyScalar(v)

			v1 := hor.Add(vert)
			v2 := lowerLeft.Add(v1)

			direction := v1.Add(v2)
			r := Ray{origin, direction}

			rgb := r.Color()

			// get intensity of colors
			ir := int(color * rgb.X)
			ig := int(color * rgb.Y)
			ib := int(color * rgb.Z)

			_, err = fmt.Fprintf(f, "%d %d %d\n", ir, ig, ib)

			check(err, "Error writing to file: %v\n")
		}
	}
}
