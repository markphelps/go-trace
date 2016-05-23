package primitives

type Material interface {
	Bounce(input Ray, hit Hit) (bool, Ray)
	Color() Vector
}
