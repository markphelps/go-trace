package primitives

import (
	"math"
)

type Sphere struct {
	Center Vector
	Radius float64
	Material
}

func NewSphere(x, y, z, radius float64, m Material) Sphere {
	return Sphere{Vector{x, y, z}, radius, m}
}

func (s *Sphere) Hit(r Ray, tMin float64, tMax float64) (bool, Hit) {
	oc := r.Origin.Subtract(s.Center)
	a := r.Direction.Dot(r.Direction)
	b := oc.Dot(r.Direction)
	c := oc.Dot(oc) - s.Radius*s.Radius
	discriminant := b*b - a*c

	hit := Hit{Material: s.Material}

	if discriminant > 0 {
		temp := (-b - math.Sqrt(discriminant)) / a
		if temp < tMax && temp > tMin {
			hit.T = temp
			hit.Point = r.Point(temp)
			hit.Normal = hit.Point.Subtract(s.Center).DivideScalar(s.Radius)
			return true, hit
		}
		temp = (-b + math.Sqrt(discriminant)) / a
		if temp < tMax && temp > tMin {
			hit.T = temp
			hit.Point = r.Point(temp)
			hit.Normal = hit.Point.Subtract(s.Center).DivideScalar(s.Radius)
			return true, hit
		}
	}
	return false, Hit{}
}
