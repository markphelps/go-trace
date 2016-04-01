package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"

	p "github.com/markphelps/go-trace/primatives"
)

const (
	nx = 400 // size of x
	ny = 200 // size of y
	ns = 100 // number of samples for aa
	c  = 255.99
)

var (
	white = p.Vector{1.0, 1.0, 1.0}
	blue  = p.Vector{0.5, 0.7, 1.0}

	camera = p.NewCamera()

	sphere = p.Sphere{p.Vector{0, 0, -1}, 0.5}
	floor  = p.Sphere{p.Vector{0, -100.5, -1}, 100}

	world = p.World{[]p.Hitable{&sphere, &floor}}
)

func check(e error, s string) {
	if e != nil {
		fmt.Fprintf(os.Stderr, s, e)
		os.Exit(1)
	}
}

func color(r *p.Ray, h p.Hitable) p.Vector {
	hit, record := h.Hit(r, 0.0, math.MaxFloat64)

	if hit {
		return record.Normal.AddScalar(1.0).MultiplyScalar(0.5)
	}

	// make unit vector so y is between -1.0 and 1.0
	unitDirection := r.Direction.Normalize()

	return gradient(&unitDirection)
}

func gradient(v *p.Vector) p.Vector {
	// scale t to be between 0.0 and 1.0
	t := 0.5 * (v.Y + 1.0)

	// linear blend: blended_value = (1 - t) * white + t * blue
	return white.MultiplyScalar(1.0 - t).Add(blue.MultiplyScalar(t))
}

func main() {

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
			rgb := p.Vector{}

			// sample rays for anti-aliasing
			for s := 0; s < ns; s++ {
				u := (float64(i) + rand.Float64()) / float64(nx)
				v := (float64(j) + rand.Float64()) / float64(ny)

				r := camera.RayAt(u, v)
				color := color(&r, &world)
				rgb = rgb.Add(color)
			}

			// average
			rgb = rgb.DivideScalar(float64(ns))

			// get intensity of colors
			ir := int(c * rgb.X)
			ig := int(c * rgb.Y)
			ib := int(c * rgb.Z)

			_, err = fmt.Fprintf(f, "%d %d %d\n", ir, ig, ib)
			check(err, "Error writing to file: %v\n")
		}
	}
}
