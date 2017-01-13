package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/markphelps/go-trace/primitive"
	"github.com/markphelps/go-trace/render"
)

var Version = "0.0.0"

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
	file                         string
	x, y, z                      float64
	allowedFiletypes             = []string{".png"}
	version                      bool
)

func main() {
	BoundFloat64Var(&fov, "fov", defaultFov, minFov, maxFov, "vertical field of view (degrees)")
	BoundIntVar(&width, "w", defaultWidth, minWidth, maxWidth, "width of image (pixels)")
	BoundIntVar(&height, "h", defaultHeight, minHeight, maxHeight, "height of image (pixels)")
	BoundIntVar(&samples, "n", defaultSamples, minSamples, maxSamples, "number of samples per pixel for AA")
	BoundFloat64Var(&aperture, "a", defaultAperture, minAperture, maxAperture, "camera aperture")
	BoundIntVar(&cpus, "cpus", runtime.NumCPU(), 1, runtime.NumCPU(), "number of CPUs to use")
	FilenameVar(&file, "o", "out.png", allowedFiletypes, "output filename")

	flag.Float64Var(&x, "x", 10, "look from X")
	flag.Float64Var(&y, "y", 4, "look from Y")
	flag.Float64Var(&z, "z", 6, "look from Z")

	flag.BoolVar(&version, "version", false, "show version and exit")

	flag.Parse()

	if version {
		fmt.Printf("go-trace %s\n", Version)
		os.Exit(0)
	}

	lookFrom := primitive.Vector{X: x, Y: y, Z: z}
	lookAt := primitive.Vector{X: 0, Y: 0, Z: -1}

	camera := primitive.NewCamera(lookFrom, lookAt, fov, float64(width)/float64(height), aperture)

	start := time.Now()

	scene := render.RandomScene()

	fmt.Printf("\nRendering %d x %d pixel scene with %d objects:", width, height, scene.Count())
	fmt.Printf("\n[%d cpus, %d samples/pixel, %.2f° fov, %.2f aperture]", cpus, samples, fov, aperture)

	ch := make(chan int, height)
	defer close(ch)

	go outputProgress(ch, height)

	image := render.Do(scene, camera, cpus, samples, width, height, ch)

	if err := writePNG(file, image); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Printf("\nDone. Elapsed: %v", time.Since(start))
	fmt.Printf("\nOutput to: %s\n", file)
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

func writePNG(path string, img image.Image) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer func() {
		if cerr := file.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	err = png.Encode(file, img)
	return err
}
