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

	// make unit vector so y is between -1.0 and 1.0
	unitDirection := r.Direction.Normalize()

	// scale t to be between 0.0 and 1.0
	t := 0.5 * (unitDirection.Y + 1.0)

	// linear blend
	// blended_value = (1 - t) * white + t * blue
	white := Vector{1.0, 1.0, 1.0}
	blue := Vector{0.5, 0.7, 1.0}

	return white.MultiplyScalar(1.0 - t).Add(blue.MultiplyScalar(t))
}

func (r Ray) HitSphere(s Sphere) bool {
	oc := r.Origin.Subtract(s.Center)
	a := r.Direction.Dot(r.Direction)
	b := 2.0 * oc.Dot(r.Direction)
	c := oc.Dot(oc) - s.Radius*s.Radius
	discriminant := b*b - 4*a*c

	return discriminant > 0
}
