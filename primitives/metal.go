package primitives

type Metal struct {
	C Vector
}

func (m Metal) Bounce(input Ray, hit Hit) (bool, Ray) {
	direction := reflect(input.Direction, hit.Normal)
	bouncedRay := Ray{hit.Point, direction}
	bounced := direction.Dot(hit.Normal) > 0
	return bounced, bouncedRay
}

func (m Metal) Color() Vector {
	return m.C
}

func reflect(v Vector, n Vector) Vector {
	b := 2 * v.Dot(n)
	return v.Subtract(n.MultiplyScalar(b))
}
