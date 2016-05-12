package primitives

import (
	"math"
	"math/rand"
)

type Vector struct {
	X, Y, Z float64
}

func VectorInUnitSphere() Vector {
	unitVector := Vector{1, 1, 1}
	for {
		r := Vector{rand.Float64(), rand.Float64(), rand.Float64()}
		p := r.MultiplyScalar(2.0).Subtract(unitVector)
		if p.SquaredLength() >= 1.0 {
			return p
		}
	}
}

func (v Vector) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func (v Vector) SquaredLength() float64 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

func (v Vector) Dot(o Vector) float64 {
	return v.X*o.X + v.Y*o.Y + v.Z*o.Z
}

func (v Vector) Normalize() Vector {
	l := v.Length()
	return Vector{v.X / l, v.Y / l, v.Z / l}
}

func (v Vector) Add(o Vector) Vector {
	return Vector{v.X + o.X, v.Y + o.Y, v.Z + o.Z}
}

func (v Vector) Subtract(o Vector) Vector {
	return Vector{v.X - o.X, v.Y - o.Y, v.Z - o.Z}
}

func (v Vector) AddScalar(t float64) Vector {
	return Vector{v.X + t, v.Y + t, v.Z + t}
}

func (v Vector) SubtractScalar(t float64) Vector {
	return Vector{v.X - t, v.Y - t, v.Z - t}
}

func (v Vector) MultiplyScalar(t float64) Vector {
	return Vector{v.X * t, v.Y * t, v.Z * t}
}

func (v Vector) DivideScalar(t float64) Vector {
	return Vector{v.X / t, v.Y / t, v.Z / t}
}
