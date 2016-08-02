package primitives

import (
	"math"
	"math/rand"
)

type Material interface {
	Bounce(input Ray, hit Hit, rnd *rand.Rand) (bool, Ray)
	Color() Color
}

type Lambertian struct {
	Attenuation Color
}

func (l Lambertian) Color() Color {
	return l.Attenuation
}

func (l Lambertian) Bounce(input Ray, hit Hit, rnd *rand.Rand) (bool, Ray) {
	direction := hit.Normal.Add(VectorInUnitSphere(rnd))
	return true, Ray{hit.Point, direction}
}

type Metal struct {
	Attenuation Color
	Fuzz        float64
}

func (m Metal) Color() Color {
	return m.Attenuation
}

func (m Metal) Bounce(input Ray, hit Hit, rnd *rand.Rand) (bool, Ray) {
	direction := input.Direction.Reflect(hit.Normal)
	bouncedRay := Ray{hit.Point, direction.Add(VectorInUnitSphere(rnd).MultiplyScalar(m.Fuzz))}
	bounced := direction.Dot(hit.Normal) > 0
	return bounced, bouncedRay
}

type Dielectric struct {
	Index float64
}

func (d Dielectric) Color() Color {
	return Color{1.0, 1.0, 1.0}
}

func (d Dielectric) Bounce(input Ray, hit Hit, rnd *rand.Rand) (bool, Ray) {
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

	if rnd.Float64() < reflectProbability {
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
