package primitives

import (
	"math"
	"math/rand"
)

type Dielectric struct {
	Index float64
}

func (d Dielectric) Color() Color {
	return Color{1.0, 1.0, 1.0}
}

func (d Dielectric) Bounce(input Ray, hit Hit) (bool, Ray) {
	var outwardNormal Vector
	var niOverNt, cosine float64

	if input.Direction.Dot(hit.Normal) > 0 {
		outwardNormal = hit.Normal.MultiplyScalar(-1)
		niOverNt = d.Index

		a := input.Direction.Dot(hit.Normal) * d.Index
		b := input.Direction.Length()

		cosine = a / b
	} else {
		outwardNormal = hit.Normal
		niOverNt = 1.0 / d.Index

		a := input.Direction.Dot(hit.Normal) * d.Index
		b := input.Direction.Length()

		cosine = -a / b
	}

	var success bool
	var refracted Vector
	var reflectProbability float64

	if success, refracted = input.Direction.Refract(outwardNormal, niOverNt); success {
		reflectProbability = d.schlick(cosine)
	} else {
		reflectProbability = 1.0
	}

	if rand.Float64() < reflectProbability {
		reflected := input.Direction.Reflect(hit.Normal)
		return true, Ray{hit.Point, reflected}
	}

	return true, Ray{hit.Point, refracted}
}

func (d Dielectric) schlick(cosine float64) float64 {
	r0 := (1 - d.Index) / (1 + d.Index)
	r0 = r0 * r0
	return r0 + (1-r0)*math.Pow((1-cosine), 5)
}
