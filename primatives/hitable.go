package primatives

type HitRecord struct {
	T         float64
	P, Normal Vector
}

type Hitable interface {
	Hit(r *Ray, tMin float64, tMax float64) (bool, HitRecord)
}
