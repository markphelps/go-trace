package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	primative "github.com/markphelps/go-trace/lib"
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

type config struct {
	width, height, ns int
	ncpus             int
	aperture, fov     float64
}

var (
	cfg      config
	filename string
	x, y, z  float64
)

func init() {
	flag.Float64Var(&cfg.fov, "fov", defaultFov, "vertical field of view (degrees)")
	flag.IntVar(&cfg.width, "w", defaultWidth, "width of image (pixels)")
	flag.IntVar(&cfg.height, "h", defaultHeight, "height of image (pixels)")
	flag.IntVar(&cfg.ns, "n", defaultSamples, "number of samples per pixel for AA")
	flag.Float64Var(&cfg.aperture, "a", defaultAperture, "camera aperture")
	flag.IntVar(&cfg.ncpus, "cpus", runtime.NumCPU(), "number of CPUs to use")

	flag.StringVar(&filename, "o", "out.png", "output filename")

	flag.Float64Var(&x, "x", 10, "look from X")
	flag.Float64Var(&y, "y", 4, "look from Y")
	flag.Float64Var(&z, "z", 6, "look from Z")
}

func main() {
	flag.Parse()

	if strings.ToLower(filepath.Ext(filename)) != ".png" {
		fmt.Println("Error: output must be a .png file")
		os.Exit(1)
	}

	cfg.fov = clampFloat(cfg.fov, minFov, maxFov)
	cfg.width = clampInt(cfg.width, minWidth, maxWidth)
	cfg.height = clampInt(cfg.height, minHeight, maxHeight)
	cfg.ns = clampInt(cfg.ns, minSamples, maxSamples)
	cfg.aperture = clampFloat(cfg.aperture, minAperture, maxAperture)
	cfg.ncpus = clampInt(cfg.ncpus, 1, runtime.NumCPU())

	lookFrom := primative.Vector{X: x, Y: y, Z: z}
	lookAt := primative.Vector{X: 0, Y: 0, Z: -1}

	camera := primative.NewCamera(lookFrom, lookAt, cfg.fov, float64(cfg.width)/float64(cfg.height), cfg.aperture)

	start := time.Now()

	scene := randomScene()

	fmt.Printf("\nRendering %d x %d pixel scene with %d objects:", cfg.width, cfg.height, scene.Count())
	fmt.Printf("\n[%d cpus, %d samples/pixel, %.2f° fov, %.2f aperture]\n", cfg.ncpus, cfg.ns, cfg.fov, cfg.aperture)

	ch := make(chan int, cfg.height)
	defer close(ch)

	go outputProgress(ch, cfg.height)

	image := render(scene, camera, cfg, ch)

	err := writePNG(filename, image)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Printf("\nDone. Elapsed: %v", time.Since(start))
	fmt.Printf("\nOutput to: %s\n", filename)
}

func outputProgress(ch chan int, rows int) {
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
