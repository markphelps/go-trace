package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"

	p "github.com/markphelps/go-trace/primitives"
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
)

func check(err error, msg string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, msg, err)
		os.Exit(1)
	}
}

func color(r p.Ray, world p.Hitable, depth int) p.Vector {
	hit, record := world.Hit(r, 0.001, math.MaxFloat64)

	if hit {
		if depth < 50 {
			bounced, bouncedRay := record.Bounce(r, record)
			if bounced {
				newColor := color(bouncedRay, world, depth+1)
				return record.Material.Color().Multiply(newColor)
			}
		}
		return p.Vector{}
	}

	return gradient(r)
}

func gradient(r p.Ray) p.Vector {
	// make unit vector so y is between -1.0 and 1.0
	v := r.Direction.Normalize()

	// scale t to be between 0.0 and 1.0
	t := 0.5 * (v.Y + 1.0)

	// linear blend: blended_value = (1 - t) * white + t * blue
	return white.MultiplyScalar(1.0 - t).Add(blue.MultiplyScalar(t))
}

func createFile() *os.File {
	f, err := os.Create("out.ppm")
	check(err, "Error opening file: %v\n")

	// http://netpbm.sourceforge.net/doc/ppm.html
	_, err = fmt.Fprintf(f, "P3\n%d %d\n255\n", nx, ny)
	check(err, "Error writting to file: %v\n")
	return f
}

func writeFile(f *os.File, rgb p.Vector) {
	// get intensity of colors
	ir := int(c * math.Sqrt(rgb.X))
	ig := int(c * math.Sqrt(rgb.Y))
	ib := int(c * math.Sqrt(rgb.Z))

	_, err := fmt.Fprintf(f, "%d %d %d\n", ir, ig, ib)
	check(err, "Error writing to file: %v\n")
}

// samples rays for anti-aliasing
func sample(world *p.World, camera *p.Camera, i, j int) p.Vector {
	rgb := p.Vector{}

	for s := 0; s < ns; s++ {
		u := (float64(i) + rand.Float64()) / float64(nx)
		v := (float64(j) + rand.Float64()) / float64(ny)

		ray := camera.RayAt(u, v)
		rgb = rgb.Add(color(ray, world, 0))
	}

	// average
	return rgb.DivideScalar(float64(ns))
}

func render(world *p.World, camera *p.Camera) {
	ticker := time.NewTicker(time.Millisecond * 100)

	go func() {
		for {
			<-ticker.C
			fmt.Print(".")
		}
	}()

	f := createFile()
	defer f.Close()

	start := time.Now()

	for j := ny - 1; j >= 0; j-- {
		for i := 0; i < nx; i++ {
			rgb := sample(world, camera, i, j)
			writeFile(f, rgb)
		}
	}

	ticker.Stop()
	fmt.Printf("\nDone.\nElapsed: %v\n", time.Since(start))
}

func main() {
	camera := p.NewCamera()

	world := p.World{}

	sphere := p.Sphere{p.Vector{0, 0, -1}, 0.5, p.Lambertian{p.Vector{0.8, 0.3, 0.3}}}
	floor := p.Sphere{p.Vector{0, -100.5, -1}, 100, p.Lambertian{p.Vector{0.8, 0.8, 0.0}}}
	left := p.Sphere{p.Vector{-1, 0, -1}, 0.5, p.Metal{p.Vector{0.8, 0.8, 0.8}, 0.0}}
	right := p.Sphere{p.Vector{1, 0, -1}, 0.5, p.Metal{p.Vector{0.8, 0.6, 0.2}, 0.3}}

	world.Add(&sphere)
	world.Add(&floor)
	world.Add(&left)
	world.Add(&right)

	render(&world, &camera)
}
