package main

import (
	"math/rand"

	primatives "github.com/markphelps/go-trace/lib"
)

func randomScene() *primatives.World {
	world := &primatives.World{}

	floor := primatives.NewSphere(0, -1000, 0, 1000,
		primatives.Lambertian{Attenuation: primatives.Color{
			R: 0.5,
			G: 0.5,
			B: 0.5}})

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
					glass := primatives.NewSphere(center.X, center.Y, center.Z, 0.2,
						primatives.Dielectric{Index: 1.5})

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
