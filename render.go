package main

import (
	"fmt"
	"image"
	"math"
	"math/rand"
	"strings"
	"time"

	primatives "github.com/markphelps/go-trace/lib"
)

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
func sample(world *primatives.World, camera *primatives.Camera, rnd *rand.Rand, config Config, i, j int) primatives.Color {
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

func render(world *primatives.World, camera *primatives.Camera, config Config) image.Image {
	img := image.NewNRGBA(image.Rect(0, 0, config.nx, config.ny))

	ch := make(chan int, config.ny)
	defer close(ch)

	for i := 0; i < config.ncpus; i++ {
		go func(i int) {
			rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
			for row := i; row < config.ny; row += config.ncpus {
				for col := 0; col < config.nx; col++ {
					rgb := sample(world, camera, rnd, config, col, row)
					img.Set(col, config.ny-row-1, rgb)
				}
				ch <- 1
			}
		}(i)
	}

	fmt.Println()
	for i := 1; i <= config.ny; i++ {
		<-ch
		pct := 100 * float64(i) / float64(config.ny)
		filled := (80 * i) / config.ny
		bar := strings.Repeat("=", filled) + strings.Repeat("-", 80-filled)
		fmt.Printf("\r[%s] %.2f%%", bar, pct)
	}
	return img
}
