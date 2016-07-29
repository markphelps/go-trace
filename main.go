package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"runtime"
	"time"

	"image"
	"image/png"

	primatives "github.com/markphelps/go-trace/lib"
)

var config struct {
	nx, ny, ns    int
	cpus          int
	aperture, fov float64
	filename      string
}

func check(err error, msg string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, msg, err)
		os.Exit(1)
	}
}

func color(r primatives.Ray, world primatives.Hitable, depth int) primatives.Color {
	hit, record := world.Hit(r, 0.001, math.MaxFloat64)

	if hit {
		if depth < 50 {
			bounced, bouncedRay := record.Bounce(r, record)
			if bounced {
				newColor := color(bouncedRay, world, depth+1)
				return record.Material.Color().Multiply(newColor)
			}
		}
		return primatives.Black
	}

	return primatives.Gradient(primatives.White, primatives.Blue, r.Direction.Normalize().Y)
}

// samples rays for anti-aliasing
func sample(world *primatives.World, camera *primatives.Camera, rnd *rand.Rand, i, j int) primatives.Color {
	rgb := primatives.Color{}

	for s := 0; s < config.ns; s++ {
		u := (float64(i) + rnd.Float64()) / float64(config.nx)
		v := (float64(j) + rnd.Float64()) / float64(config.ny)

		ray := camera.RayAt(u, v, rnd)
		rgb = rgb.Add(color(ray, world, 0))
	}

	// average
	return rgb.DivideScalar(float64(config.ns))
}

func render(world *primatives.World, camera *primatives.Camera) image.Image {
	img := image.NewNRGBA(image.Rect(0, 0, config.nx, config.ny))

	ch := make(chan int, config.ny)
	defer close(ch)

	for i := 0; i < config.cpus; i++ {
		go func(i int) {
			rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
			for row := i; row < config.ny; row += config.cpus {
				for col := 0; col < config.nx; col++ {
					rgb := sample(world, camera, rnd, col, row)
					img.Set(col, config.ny-row-1, rgb.Sqrt())
				}
				ch <- 1
			}
		}(i)
	}

	fmt.Print("0.00% complete")
	for i := 0; i < config.ny; i++ {
		<-ch
		pct := 100 * float64(i+1) / float64(config.ny)
		fmt.Printf("\r%.2f %% complete", pct)
	}
	return img
}

func write(img image.Image) {
	file, err := os.Create(config.filename)
	check(err, "Error opening file: %v\n")
	defer file.Close()

	err = png.Encode(file, img)
	check(err, "Error writing to file: %v\n")
}

func scene() *primatives.World {
	world := &primatives.World{}

	floor := primatives.NewSphere(0, -1000, 0, 1000, primatives.Lambertian{Attenuation: primatives.Color{R: 0.5, G: 0.5, B: 0.5}})
	world.Add(floor)

	for a := -11; a < 11; a++ {
		for b := -11; b < 11; b++ {
			material := rand.Float64()

			center := primatives.Vector{
				X: float64(a) + 0.9*rand.Float64(),
				Y: 0.2,
				Z: float64(b) + 0.9*rand.Float64()}

			if center.Subtract(primatives.Vector{X: 4, Y: 0.2, Z: 0}).Length() > 0.9 {
				if material < 0.8 {
					lambertian := primatives.NewSphere(center.X, center.Y, center.Z, 0.2,
						primatives.Lambertian{Attenuation: primatives.Color{
							R: rand.Float64() * rand.Float64(),
							G: rand.Float64() * rand.Float64(),
							B: rand.Float64() * rand.Float64()}})

					world.Add(lambertian)
				} else if material < 0.95 {
					metal := primatives.NewSphere(center.X, center.Y, center.Z, 0.2,
						primatives.Metal{Attenuation: primatives.Color{
							R: 0.5 * (1.0 + rand.Float64()),
							G: 0.5 * (1.0 + rand.Float64()),
							B: 0.5 * (1.0 + rand.Float64())},
							Fuzz: 0.5 + rand.Float64()})

					world.Add(metal)
				} else {
					glass := primatives.NewSphere(center.X, center.Y, center.Z, 0.2, primatives.Dielectric{Index: 1.5})
					world.Add(glass)
				}
			}
		}
	}

	glass := primatives.NewSphere(0, 1, 0, 1.0, primatives.Dielectric{Index: 1.5})
	lambertian := primatives.NewSphere(-4, 1, 0, 1.0, primatives.Lambertian{Attenuation: primatives.Color{R: 0.4, G: 0.0, B: 0.1}})
	metal := primatives.NewSphere(4, 1, 0, 1.0, primatives.Metal{Attenuation: primatives.Color{R: 0.7, G: 0.6, B: 0.5}, Fuzz: 0.0})

	world.AddAll(glass, lambertian, metal)
	return world
}

func init() {
	// seed for scene generation
	rand.Seed(time.Now().UnixNano())
}

func main() {
	lookFrom := primatives.Vector{}

	flag.Float64Var(&lookFrom.X, "x", 10, "look from X")
	flag.Float64Var(&lookFrom.Y, "y", 4, "look from Y")
	flag.Float64Var(&lookFrom.Z, "z", 6, "look from Z")

	flag.Float64Var(&config.fov, "fov", 75.0, "vertical field of view (degrees)")
	flag.IntVar(&config.nx, "width", 600, "width of image")
	flag.IntVar(&config.ny, "height", 500, "height of image")
	flag.IntVar(&config.ns, "samples", 100, "number of samples for anti-aliasing")
	flag.Float64Var(&config.aperture, "aperture", 0.01, "camera aperture")
	flag.StringVar(&config.filename, "out", "out.png", "output filename")
	flag.IntVar(&config.cpus, "cpus", runtime.NumCPU(), "number of CPUs to use")

	flag.Parse()

	if config.cpus > runtime.NumCPU() {
		config.cpus = runtime.NumCPU()
	}

	lookAt := primatives.Vector{X: 0, Y: 0, Z: -1}

	camera := primatives.NewCamera(lookFrom, lookAt, config.fov, float64(config.nx)/float64(config.ny), config.aperture)

	start := time.Now()

	scene := scene()

	fmt.Printf("\nRendering %d x %d pixel scene with %d objects [%d cpus, %d samples]...\n", config.nx, config.ny, scene.Count(), config.cpus, config.ns)

	image := render(scene, camera)
	write(image)

	fmt.Printf("\n\nDone. Elapsed: %v\nOutput to: %s\n", time.Since(start), config.filename)
}
