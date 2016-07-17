package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"

	p "github.com/markphelps/go-trace/primitives"
)

const (
	c               = 255.99
	defaultNx       = 600  // size of x
	defaultNy       = 500  // size of y
	defaultNs       = 100  // number of samples for aa
	defaultAperture = 0.01 // smaller aperture means less blury
	defaultFov      = 90   // field of view
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

func createFile(nx, ny int) *os.File {
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
func sample(world *p.World, camera *p.Camera, i, j, nx, ny, ns int) p.Color {
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

func render(world *p.World, camera *p.Camera, nx, ny, ns int) {
	f := createFile(nx, ny)
	defer f.Close()

	ch := make(chan int)
	defer close(ch)

	go func() {
		for {
			if i, rendering := <-ch; rendering {
				pct := 100 * float64(i) / float64(ny)
				fmt.Printf("\r%.2f %% complete", pct)
			}
		}
	}()

	row := 1
	for j := ny - 1; j >= 0; j-- {
		ch <- row
		row++
		for i := 0; i < nx; i++ {
			rgb := sample(world, camera, i, j, nx, ny, ns)
			writeFile(f, rgb)
		}
	}
}

func randomScene() *p.World {
	world := &p.World{}

	floor := p.NewSphere(0, -1000, 0, 1000, p.Lambertian{C: p.Color{R: 0.5, G: 0.5, B: 0.5}})
	world.Add(floor)

	for a := -11; a < 11; a++ {
		for b := -11; b < 11; b++ {
			material := rand.Float64()

			center := p.Vector{
				X: float64(a) + 0.9*rand.Float64(),
				Y: 0.2,
				Z: float64(b) + 0.9*rand.Float64()}

			if center.Subtract(p.Vector{X: 4, Y: 0.2, Z: 0}).Length() > 0.9 {
				if material < 0.8 {
					lambertian := p.NewSphere(center.X, center.Y, center.Z, 0.2, p.Lambertian{C: p.Color{
						R: rand.Float64() * rand.Float64(),
						G: rand.Float64() * rand.Float64(),
						B: rand.Float64() * rand.Float64()}})

					world.Add(lambertian)
				} else if material < 0.95 {
					metal := p.NewSphere(center.X, center.Y, center.Z, 0.2,
						p.Metal{C: p.Color{
							R: 0.5 * (1.0 + rand.Float64()),
							G: 0.5 * (1.0 + rand.Float64()),
							B: 0.5 * (1.0 + rand.Float64())},
							Fuzz: 0.5 + rand.Float64()})

					world.Add(metal)
				} else {
					glass := p.NewSphere(center.X, center.Y, center.Z, 0.2, p.Dielectric{Index: 1.5})
					world.Add(glass)
				}
			}
		}
	}

	glass := p.NewSphere(0, 1, 0, 1.0, p.Dielectric{Index: 1.5})
	lambertian := p.NewSphere(-4, 1, 0, 1.0, p.Lambertian{C: p.Color{R: 0.4, G: 0.0, B: 0.1}})
	metal := p.NewSphere(4, 1, 0, 1.0, p.Metal{C: p.Color{R: 0.7, G: 0.6, B: 0.5}, Fuzz: 0.0})

	world.AddAll(glass, lambertian, metal)
	return world
}

func main() {

	fov := flag.Int("fov", defaultFov, "field of view")
	nx := flag.Int("nx", defaultNx, "width of image")
	ny := flag.Int("ny", defaultNy, "height of image")
	ns := flag.Int("ns", defaultNs, "number of samples for AA")
	aperture := flag.Float64("a", defaultAperture, "aperture")

	flag.Parse()

	lookFrom := p.Vector{X: 1, Y: 4, Z: 6}
	lookAt := p.Vector{X: 0, Y: 0, Z: -1}

	camera := p.NewCamera(lookFrom, lookAt, float64(*fov), float64(*nx)/float64(*ny), *aperture)

	start := time.Now()

	scene := randomScene()

	fmt.Printf("\nRendering %d X %d scene with %d objects...\n", *nx, *ny, scene.Size())

	render(scene, camera, *nx, *ny, *ns)

	fmt.Printf("\nDone. Elapsed: %v\n", time.Since(start))
}
