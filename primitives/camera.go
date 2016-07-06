package primitives

import "math"

type Camera struct {
	lowerLeft, horizontal, vertical, origin Vector
}

func NewCamera(lookFrom, lookAt, vUp Vector, vFov, aspect float64) Camera {

	theta := vFov * math.Pi / 180
	halfHeight := math.Tan(theta / 2)
	halfWidth := aspect * halfHeight

	c := Camera{}

	c.origin = lookFrom

	w := lookFrom.Subtract(lookAt).Normalize()
	u := vUp.Cross(w).Normalize()
	v := w.Cross(u)

	c.lowerLeft = c.origin.Subtract(u.MultiplyScalar(halfWidth)).Subtract(v.MultiplyScalar(halfHeight)).Subtract(w)
	c.horizontal = u.MultiplyScalar(2 * halfWidth)
	c.vertical = v.MultiplyScalar(2 * halfHeight)

	return c
}

func (c Camera) RayAt(s, t float64) Ray {
	horizontal := c.horizontal.MultiplyScalar(s)
	vertical := c.vertical.MultiplyScalar(t)

	direction := c.lowerLeft.Add(horizontal).Add(vertical).Subtract(c.origin)
	return Ray{c.origin, direction}
}
