package primitives

type Lambertian struct {
	Attenuation Color
}

func (l Lambertian) Color() Color {
	return l.Attenuation
}

func (l Lambertian) Bounce(input Ray, hit Hit) (bool, Ray) {
	direction := hit.Normal.Add(VectorInUnitSphere())
	return true, Ray{hit.Point, direction}
}
