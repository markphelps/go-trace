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

func write(img image.Image) {
	file, err := os.Create(config.filename)
	check(err, "Error opening file: %v\n")
	defer file.Close()

	err = png.Encode(file, img)
	check(err, "Error writing to file: %v\n")
}

func init() {
	// seed for scene generation
	rand.Seed(time.Now().UnixNano())
}

func main() {
	x := flag.Float64("x", 10, "look from X")
	y := flag.Float64("y", 4, "look from Y")
	z := flag.Float64("z", 6, "look from Z")

	flag.Float64Var(&config.fov, "fov", 75.0, "vertical field of view (degrees)")
	flag.IntVar(&config.nx, "width", 600, "width of image")
	flag.IntVar(&config.ny, "height", 500, "height of image")
	flag.IntVar(&config.ns, "samples", 100, "number of samples for anti-aliasing")
	flag.Float64Var(&config.aperture, "aperture", 0.01, "camera aperture")
	flag.StringVar(&config.filename, "out", "out.png", "output filename")
	flag.IntVar(&config.cpus, "cpus", runtime.NumCPU(), "number of CPUs to use")

	flag.Parse()

	// upperbound to numCPU
	if config.cpus > runtime.NumCPU() {
		config.cpus = runtime.NumCPU()
	}

	// 120 degrees max fov
	config.fov = math.Min(config.fov, 120.0)

	lookFrom := primatives.Vector{X: *x, Y: *y, Z: *z}
	lookAt := primatives.Vector{X: 0, Y: 0, Z: -1}

	camera := primatives.NewCamera(lookFrom, lookAt, config.fov, float64(config.nx)/float64(config.ny), config.aperture)

	start := time.Now()

	scene := randomScene()

	fmt.Printf("\nRendering %d x %d pixel scene with %d objects...\n", config.nx, config.ny, scene.Count())
	fmt.Printf("[%d cpus, %d samples, %.2f fov]\n\n", config.cpus, config.ns, config.fov)

	image := render(scene, camera)
	write(image)

	fmt.Printf("\n\nDone. Elapsed: %v\nOutput to: %s\n", time.Since(start), config.filename)
}
