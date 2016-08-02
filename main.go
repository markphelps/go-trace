package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	primatives "github.com/markphelps/go-trace/lib"
)

const (
	maxFov      = 120.0
	maxWidth    = 4096
	maxHeight   = 2160
	maxSamples  = 5000
	maxAperture = 0.9

	minFov      = 10.0
	minWidth    = 200
	minHeight   = 100
	minSamples  = 1
	minAperture = 0.001

	defaultFov      = 75.0
	defaultWidth    = 600
	defaultHeight   = 500
	defaultSamples  = 100
	defaultAperture = 0.01
)

type Config struct {
	nx, ny, ns    int
	ncpus         int
	aperture, fov float64
}

func main() {
	var filename string
	var x, y, z float64

	config := Config{}

	flag.Float64Var(&config.fov, "fov", defaultFov, "vertical field of view (degrees)")
	flag.IntVar(&config.nx, "width", defaultWidth, "width of image (pixels)")
	flag.IntVar(&config.ny, "height", defaultHeight, "height of image (pixels)")
	flag.IntVar(&config.ns, "samples", defaultSamples, "number of samples per pixel for AA")
	flag.Float64Var(&config.aperture, "aperture", defaultAperture, "camera aperture")
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

	config.fov = boundFloat(minFov, maxFov, config.fov)
	config.nx = boundInt(minWidth, maxWidth, config.nx)
	config.ny = boundInt(minHeight, maxHeight, config.ny)
	config.ns = boundInt(minSamples, maxSamples, config.ns)
	config.aperture = boundFloat(minAperture, maxAperture, config.aperture)
	config.ncpus = boundInt(1, runtime.NumCPU(), config.ncpus)

	lookFrom := primatives.Vector{X: x, Y: y, Z: z}
	lookAt := primatives.Vector{X: 0, Y: 0, Z: -1}

	camera := primatives.NewCamera(lookFrom,
		lookAt,
		config.fov,
		float64(config.nx)/float64(config.ny),
		config.aperture)

	start := time.Now()

	scene := randomScene()

	fmt.Printf("\nRendering %d x %d pixel scene with %d objects:\n", config.nx, config.ny, scene.Count())
	fmt.Printf("[%d cpus, %d samples/pixel, %.2f° fov, %.2f aperture]\n", config.ncpus, config.ns, config.fov, config.aperture)

	image := render(scene, camera, config)
	writePNG(filename, image)

	fmt.Printf("\n\nDone. Elapsed: %v\nOutput to: %s\n", time.Since(start), filename)
}
