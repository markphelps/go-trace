package main

import (
	"image"
	"math"
	"math/rand"
	"sync"
	"time"

	primatives "github.com/markphelps/go-trace/lib"
)

func color(ray primatives.Ray, hitable primatives.Hitable, rnd *rand.Rand, depth int) primatives.Color {
	hit, record := hitable.Hit(ray, 0.001, math.MaxFloat64)

	if hit {
		if depth < 50 {
			bounced, bouncedRay := record.Bounce(ray, record, rnd)
			if bounced {
				newColor := color(bouncedRay, hitable, rnd, depth+1)
				return record.Material.Color().Multiply(newColor)
			}
		}
		return primatives.Black
	}

	return primatives.Gradient(primatives.White, primatives.Blue, ray.Direction.Normalize().Y)
}

// samples rays for anti-aliasing
func sample(hitable primatives.Hitable, camera *primatives.Camera, rnd *rand.Rand, cfg config, i, j int) primatives.Color {
	rgb := primatives.Color{}

	for s := 0; s < cfg.ns; s++ {
		u := (float64(i) + rnd.Float64()) / float64(cfg.width)
		v := (float64(j) + rnd.Float64()) / float64(cfg.height)

		ray := camera.RayAt(u, v, rnd)
		rgb = rgb.Add(color(ray, hitable, rnd, 0))
	}

	// average
	return rgb.DivideScalar(float64(cfg.ns))
}

func render(hitable primatives.Hitable, camera *primatives.Camera, cfg config, ch chan int) image.Image {
	img := image.NewNRGBA(image.Rect(0, 0, cfg.width, cfg.height))

	var wg sync.WaitGroup

	for i := 0; i < cfg.ncpus; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()
			rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

			for row := i; row < cfg.height; row += cfg.ncpus {
				for col := 0; col < cfg.width; col++ {
					rgb := sample(hitable, camera, rnd, cfg, col, row)
					img.Set(col, cfg.height-row-1, rgb)
				}
				ch <- 1
			}
		}(i)
	}

	wg.Wait()
	return img
}
