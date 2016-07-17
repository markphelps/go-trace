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
	c = 255.99
)

type config struct {
	nx, ny, ns    int
	aperture, fov float64
	filename      string
}

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

func createFile(cfg config) *os.File {
	f, err := os.Create(cfg.filename)
	check(err, "Error opening file: %v\n")

	// http://netpbm.sourceforge.net/doc/ppm.html
	_, err = fmt.Fprintf(f, "P3\n%d %d\n255\n", cfg.nx, cfg.ny)
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
func sample(world *p.World, camera *p.Camera, cfg config, i, j int) p.Color {
	rgb := p.Color{}

	for s := 0; s < cfg.ns; s++ {
		u := (float64(i) + rand.Float64()) / float64(cfg.nx)
		v := (float64(j) + rand.Float64()) / float64(cfg.ny)

		ray := camera.RayAt(u, v)
		rgb = rgb.Add(color(ray, world, 0))
	}

	// average
	return rgb.DivideScalar(float64(cfg.ns))
}

func render(world *p.World, camera *p.Camera, cfg config) {
	f := createFile(cfg)
	defer f.Close()

	ch := make(chan int)
	defer close(ch)

	go func() {
		for {
			if i, rendering := <-ch; rendering {
				pct := 100 * float64(i) / float64(cfg.ny)
				fmt.Printf("\r%.2f %% complete", pct)
			}
		}
	}()

	row := 1
	for j := cfg.ny - 1; j >= 0; j-- {
		ch <- row
		row++
		for i := 0; i < cfg.nx; i++ {
			rgb := sample(world, camera, cfg, i, j)
			writeFile(f, rgb)
		}
	}
}

func scene() *p.World {
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

	cfg := config{}

	flag.Float64Var(&cfg.fov, "fov", 90.0, "field of view")
	flag.IntVar(&cfg.nx, "nx", 600, "width of image")
	flag.IntVar(&cfg.ny, "ny", 500, "height of image")
	flag.IntVar(&cfg.ns, "ns", 100, "number of samples for AA")
	flag.Float64Var(&cfg.aperture, "a", 0.01, "aperture")
	flag.StringVar(&cfg.filename, "o", "out.ppm", "output filename")

	flag.Parse()

	lookFrom := p.Vector{X: 1, Y: 4, Z: 6}
	lookAt := p.Vector{X: 0, Y: 0, Z: -1}

	camera := p.NewCamera(lookFrom, lookAt, cfg.fov, float64(cfg.nx)/float64(cfg.ny), cfg.aperture)

	start := time.Now()

	scene := scene()

	fmt.Printf("\nRendering %d x %d pixel scene with %d objects...\n", cfg.nx, cfg.ny, scene.Count())

	render(scene, camera, cfg)

	fmt.Printf("\nDone. Elapsed: %v\n", time.Since(start))
}
