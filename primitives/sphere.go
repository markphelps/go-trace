package primitives

import (
	"math"
)

type Sphere struct {
	Center Vector
	Radius float64
}

func (s *Sphere) Hit(r *Ray, tMin float64, tMax float64) (bool, Hit) {
	oc := r.Origin.Subtract(s.Center)
	a := r.Direction.Dot(r.Direction)
	b := 2.0 * oc.Dot(r.Direction)
	c := oc.Dot(oc) - s.Radius*s.Radius
	discriminant := b*b - 4*a*c

	hit := Hit{}

	if discriminant > 0.0 {
		t := (-b - math.Sqrt(b*b-a*c)) / a
		if t < tMax && t > tMin {
			hit.T = t
			hit.P = r.Point(t)
			hit.Normal = (hit.P.Subtract(s.Center)).DivideScalar(s.Radius)
			return true, hit
		}
		t = (-b + math.Sqrt(b*b-a*c)) / a
		if t < tMax && t > tMin {
			hit.T = t
			hit.P = r.Point(t)
			hit.Normal = (hit.P.Subtract(s.Center)).DivideScalar(s.Radius)
			return true, hit
		}
	}
	return false, hit
}
