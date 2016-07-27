package primitives

import (
	"math"
	"math/rand"
)

var vUp = Vector{X: 0, Y: 1, Z: 0}

type Camera struct {
	lowerLeft, horizontal, vertical, origin, u, v, w Vector
	lensRadius                                       float64
}

func NewCamera(lookFrom, lookAt Vector, vFov, aspect, aperture float64) *Camera {
	c := Camera{}

	c.origin = lookFrom
	c.lensRadius = aperture / 2

	theta := vFov * math.Pi / 180
	halfHeight := math.Tan(theta / 2)
	halfWidth := aspect * halfHeight

	w := lookFrom.Subtract(lookAt).Normalize()
	u := vUp.Cross(w).Normalize()
	v := w.Cross(u)

	focusDist := lookFrom.Subtract(lookAt).Length()

	x := u.MultiplyScalar(halfWidth * focusDist)
	y := v.MultiplyScalar(halfHeight * focusDist)

	c.lowerLeft = c.origin.Subtract(x).Subtract(y).Subtract(w.MultiplyScalar(focusDist))
	c.horizontal = x.MultiplyScalar(2)
	c.vertical = y.MultiplyScalar(2)

	c.w = w
	c.u = u
	c.v = v

	return &c
}

func (c *Camera) RayAt(s, t float64) Ray {
	rd := randomInUnitDisc().MultiplyScalar(c.lensRadius)
	offset := c.u.MultiplyScalar(rd.X).Add(c.v.MultiplyScalar(rd.Y))

	horizontal := c.horizontal.MultiplyScalar(s)
	vertical := c.vertical.MultiplyScalar(t)

	origin := c.origin.Add(offset)
	direction := c.lowerLeft.Add(horizontal).Add(vertical).Subtract(c.origin).Subtract(offset)
	return Ray{origin, direction}
}

func randomInUnitDisc() Vector {
	var p Vector
	for {
		p = Vector{rand.Float64(), rand.Float64(), 0}.MultiplyScalar(2).Subtract(Vector{1, 1, 0})
		if p.Dot(p) < 1.0 {
			return p
		}
	}
}
