package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/markphelps/go-trace/primitive"
	"github.com/markphelps/go-trace/render"
)

const (
	maxFov      = 120.0
	maxWidth    = 4096
	maxHeight   = 2160
	maxSamples  = 1000
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

	progressBarWidth = 80
)

var (
	aperture, fov                float64
	width, height, samples, cpus int
	filename                     string
	x, y, z                      float64
)

func main() {
	flag.Float64Var(&fov, "fov", defaultFov, "vertical field of view (degrees)")
	flag.IntVar(&width, "w", defaultWidth, "width of image (pixels)")
	flag.IntVar(&height, "h", defaultHeight, "height of image (pixels)")
	flag.IntVar(&samples, "n", defaultSamples, "number of samples per pixel for AA")
	flag.Float64Var(&aperture, "a", defaultAperture, "camera aperture")
	flag.IntVar(&cpus, "cpus", runtime.NumCPU(), "number of CPUs to use")
	flag.StringVar(&filename, "o", "out.png", "output filename")
	flag.Float64Var(&x, "x", 10, "look from X")
	flag.Float64Var(&y, "y", 4, "look from Y")
	flag.Float64Var(&z, "z", 6, "look from Z")

	flag.Parse()

	if strings.ToLower(filepath.Ext(filename)) != ".png" {
		fmt.Println("Error: output must be a .png file")
		os.Exit(1)
	}

	fov = clampFloat(fov, minFov, maxFov)
	width = clampInt(width, minWidth, maxWidth)
	height = clampInt(height, minHeight, maxHeight)
	samples = clampInt(samples, minSamples, maxSamples)
	aperture = clampFloat(aperture, minAperture, maxAperture)
	cpus = clampInt(cpus, 1, runtime.NumCPU())

	lookFrom := primitive.Vector{X: x, Y: y, Z: z}
	lookAt := primitive.Vector{X: 0, Y: 0, Z: -1}

	camera := primitive.NewCamera(lookFrom, lookAt, fov, float64(width)/float64(height), aperture)

	start := time.Now()

	scene := render.RandomScene()

	fmt.Printf("\nRendering %d x %d pixel scene with %d objects:", width, height, scene.Count())
	fmt.Printf("\n[%d cpus, %d samples/pixel, %.2f° fov, %.2f aperture]\n", cpus, samples, fov, aperture)

	ch := make(chan int, height)
	defer close(ch)

	go outputProgress(ch, height)

	image := render.Do(scene, camera, cpus, samples, width, height, ch)

	err := writePNG(filename, image)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Printf("\nDone. Elapsed: %v", time.Since(start))
	fmt.Printf("\nOutput to: %s\n", filename)
}

func outputProgress(ch <-chan int, rows int) {
	fmt.Println()
	for i := 1; i <= rows; i++ {
		<-ch
		pct := 100 * float64(i) / float64(rows)
		filled := (progressBarWidth * i) / rows
		bar := strings.Repeat("=", filled) + strings.Repeat("-", progressBarWidth-filled)
		fmt.Printf("\r[%s] %.2f%%", bar, pct)
	}
	fmt.Println()
}
