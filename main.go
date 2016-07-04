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

func check(err error, msg string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, msg, err)
		os.Exit(1)
	}
}

func color(r p.Ray, world p.Hitable, depth int) p.Color {
	hit, record := world.Hit(r, 0.001, math.MaxFloat64)

	if hit {
		if depth < 50 {
			bounced, bouncedRay := record.Bounce(r, record)
			if bounced {
				newColor := color(bouncedRay, world, depth+1)
				return record.Material.Color().Multiply(newColor)
			}
		}
		return p.Black
	}

	return p.Gradient(p.White, p.Blue, r.Direction.Normalize().Y)
}

func createFile() *os.File {
	f, err := os.Create("out.ppm")
	check(err, "Error opening file: %v\n")

	// http://netpbm.sourceforge.net/doc/ppm.html
	_, err = fmt.Fprintf(f, "P3\n%d %d\n255\n", nx, ny)
	check(err, "Error writting to file: %v\n")
	return f
}

func writeFile(f *os.File, rgb p.Color) {
	// get intensity of colors with gamma-2 correction
	ir := int(c * math.Sqrt(rgb.R))
	ig := int(c * math.Sqrt(rgb.G))
	ib := int(c * math.Sqrt(rgb.B))

	_, err := fmt.Fprintf(f, "%d %d %d\n", ir, ig, ib)
	check(err, "Error writing to file: %v\n")
}

// samples rays for anti-aliasing
func sample(world *p.World, camera *p.Camera, i, j int) p.Color {
	rgb := p.Color{}

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
	camera := p.NewCamera(90, float64(nx)/float64(ny))

	world := p.World{}

	radius := math.Cos(math.Pi / 4)

	blue := p.NewSphere(-radius, 0, -1, radius, p.Lambertian{p.Color{0, 0, 1}})
	red := p.NewSphere(radius, 0, -1, radius, p.Lambertian{p.Color{1, 0, 0}})

	world.AddAll(&blue, &red)

	render(&world, &camera)
}
