package main

type Ray struct {
	Origin, Direction Vector
}

func (r Ray) Point(t float64) Vector {
	b := r.Direction.MultiplyScalar(t)
	a := r.Origin
	return a.Add(b)
}

func (r Ray) Color() Vector {
	sphere := Sphere{Center: Vector{0, 0, -1}, Radius: 0.5}

	if r.HitSphere(sphere) {
		return Vector{1.0, 0.0, 0.0} // red
	}

	unitDirection := r.Direction.Normalize()
	t := 0.5 * (unitDirection.Y + 1.0)

	white := Vector{1.0, 1.0, 1.0}.MultiplyScalar(1.0 - t)
	blue := Vector{0.5, 0.7, 1.0}.MultiplyScalar(t)

	return white.Add(blue)
}

func (r Ray) HitSphere(s Sphere) bool {
	oc := r.Origin.Subtract(s.Center)
	a := r.Direction.Dot(r.Direction)
	b := 2.0 * oc.Dot(r.Direction)
	c := oc.Dot(oc) - s.Radius*s.Radius
	discriminant := b*b - 4*a*c

	return discriminant > 0
}
