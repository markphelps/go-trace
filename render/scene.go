package render

import (
	"math/rand"
	"time"

	"github.com/markphelps/go-trace/primitive"
)

func init() {
	// seed for scene generation
	rand.Seed(time.Now().UnixNano())
}

// RandomScene returns a 'random' scene
func RandomScene() *primitive.World {
	world := &primitive.World{}

	floor := primitive.NewSphere(0, -1000, 0, 1000,
		primitive.Lambertian{Attenuation: primitive.Color{
			R: 0.5,
			G: 0.5,
			B: 0.5}})

	world.Add(floor)

	for a := -11; a < 11; a++ {
		for b := -11; b < 11; b++ {
			material := rand.Float64()

			center := primitive.Vector{
				X: float64(a) + 0.9*rand.Float64(),
				Y: 0.2,
				Z: float64(b) + 0.9*rand.Float64()}

			if center.Subtract(primitive.Vector{X: 4, Y: 0.2, Z: 0}).Length() > 0.9 {
				if material < 0.8 {
					lambertian := primitive.NewSphere(center.X, center.Y, center.Z, 0.2,
						primitive.Lambertian{Attenuation: primitive.Color{
							R: rand.Float64() * rand.Float64(),
							G: rand.Float64() * rand.Float64(),
							B: rand.Float64() * rand.Float64()}})

					world.Add(lambertian)
				} else if material < 0.95 {
					metal := primitive.NewSphere(center.X, center.Y, center.Z, 0.2,
						primitive.Metal{Attenuation: primitive.Color{
							R: 0.5 * (1.0 + rand.Float64()),
							G: 0.5 * (1.0 + rand.Float64()),
							B: 0.5 * (1.0 + rand.Float64())},
							Fuzz: 0.5 + rand.Float64()})

					world.Add(metal)
				} else {
					glass := primitive.NewSphere(center.X, center.Y, center.Z, 0.2,
						primitive.Dielectric{Index: 1.5})

					world.Add(glass)
				}
			}
		}
	}

	glass := primitive.NewSphere(0, 1, 0, 1.0, primitive.Dielectric{Index: 1.5})
	lambertian := primitive.NewSphere(-4, 1, 0, 1.0, primitive.Lambertian{Attenuation: primitive.Color{R: 0.4, G: 0.0, B: 0.1}})
	metal := primitive.NewSphere(4, 1, 0, 1.0, primitive.Metal{Attenuation: primitive.Color{R: 0.7, G: 0.6, B: 0.5}, Fuzz: 0.0})

	world.AddAll(glass, lambertian, metal)
	return world
}
