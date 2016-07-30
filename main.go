package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"time"

	primatives "github.com/markphelps/go-trace/lib"
)

type Configuration struct {
	nx, ny, ns    int
	ncpus         int
	aperture, fov float64
}

func init() {
	// seed for scene generation
	rand.Seed(time.Now().UnixNano())
}

func main() {
	var filename string
	var x, y, z float64

	config := Configuration{}

	flag.Float64Var(&config.fov, "fov", 75.0, "vertical field of view (degrees)")
	flag.IntVar(&config.nx, "width", 600, "width of image")
	flag.IntVar(&config.ny, "height", 500, "height of image")
	flag.IntVar(&config.ns, "samples", 100, "number of samples for anti-aliasing")
	flag.Float64Var(&config.aperture, "aperture", 0.01, "camera aperture")
	flag.IntVar(&config.ncpus, "cpus", runtime.NumCPU(), "number of CPUs to use")

	flag.StringVar(&filename, "out", "out.png", "output filename")

	flag.Float64Var(&x, "x", 10, "look from X")
	flag.Float64Var(&y, "y", 4, "look from Y")
	flag.Float64Var(&z, "z", 6, "look from Z")

	flag.Parse()

	if filepath.Ext(filename) != ".png" {
		fmt.Fprintf(os.Stderr, "Error: output must be a .png file\n")
		os.Exit(1)
	}

	config.ncpus = boundInt(1, runtime.NumCPU(), config.ncpus)
	config.fov = boundFloat(10.0, 120.0, config.fov)

	lookFrom := primatives.Vector{X: x, Y: y, Z: z}
	lookAt := primatives.Vector{X: 0, Y: 0, Z: -1}

	camera := primatives.NewCamera(lookFrom,
		lookAt,
		config.fov,
		float64(config.nx)/float64(config.ny),
		config.aperture)

	start := time.Now()

	scene := randomScene()

	fmt.Printf("\nRendering %d x %d pixel scene with %d objects...\n", config.nx, config.ny, scene.Count())
	fmt.Printf("[%d cpus, %d samples, %.2f fov]\n", config.ncpus, config.ns, config.fov)

	image := render(scene, camera, config)
	writePNG(filename, image)

	fmt.Printf("\n\nDone. Elapsed: %v\nOutput to: %s\n", time.Since(start), filename)
}

func boundFloat(min, max, value float64) float64 {
	if value > max {
		value = max
	} else if value < min {
		value = min
	}
	return value
}

func boundInt(min, max, value int) int {
	if value > max {
		value = max
	} else if value < min {
		value = min
	}
	return value
}
