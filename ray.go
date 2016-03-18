package main

type Ray struct {
	Origin, Direction Vector
}

func (r *Ray) Point(t float64) Vector {
	b := r.Direction.MultiplyScalar(t)
	a := r.Origin
	return a.Add(b)
}

func (r *Ray) Color() Vector {
	unitDirection := r.Direction.Normalize()
	t := 0.5 * (unitDirection.Y + 1.0)

	a := Vector{1.0, 1.0, 1.0}
	a = a.MultiplyScalar(1.0 - t)

	b := Vector{0.5, 0.7, 1.0}
	b = b.MultiplyScalar(t)

	return a.Add(b)
}
