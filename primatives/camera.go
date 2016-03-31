package primatives

type Camera struct {
	LowerLeft, Horizontal, Vertical, Origin Vector
}

func NewCamera() *Camera {
	c := new(Camera)
	c.LowerLeft = Vector{-2.0, -1.0, -1.0}
	c.Horizontal = Vector{4.0, 0.0, 0.0}
	c.Vertical = Vector{0.0, 2.0, 0.0}
	c.Origin = Vector{0.0, 0.0, 0.0}
	return c
}

func (c *Camera) Position(u float64, v float64) Vector {
	horizontal := c.Horizontal.MultiplyScalar(u)
	vertical := c.Vertical.MultiplyScalar(v)
	return horizontal.Add(vertical)
}

func (c *Camera) Direction(position Vector) Vector {
	return c.LowerLeft.Add(position)
}
