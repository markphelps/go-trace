package primatives

import (
	"math"
)

type Sphere struct {
	Center Vector
	Radius float64
}

func (s *Sphere) Hit(r *Ray, tMin float64, tMax float64) (bool, HitRecord) {
	oc := r.Origin.Subtract(s.Center)
	a := r.Direction.Dot(r.Direction)
	b := 2.0 * oc.Dot(r.Direction)
	c := oc.Dot(oc) - s.Radius*s.Radius
	discriminant := b*b - 4*a*c

	rec := HitRecord{}

	if discriminant > 0.0 {
		t := (-b - math.Sqrt(b*b-a*c)) / a
		if t < tMax && t > tMin {
			rec.T = t
			rec.P = r.Point(t)
			rec.Normal = (rec.P.Subtract(s.Center)).DivideScalar(s.Radius)
			return true, rec
		}
		t = (-b + math.Sqrt(b*b-a*c)) / a
		if t < tMax && t > tMin {
			rec.T = t
			rec.P = r.Point(t)
			rec.Normal = (rec.P.Subtract(s.Center)).DivideScalar(s.Radius)
			return true, rec
		}
	}
	return false, rec
}
