package primatives

type World struct {
	Elements []Hitable
}

func (w *World) Hit(r *Ray, tMin float64, tMax float64) (bool, Hit) {
	hitAnything := false
	closest := tMax
	record := Hit{}

	for _, element := range w.Elements {
		hit, tempRecord := element.Hit(r, tMin, closest)

		if hit {
			hitAnything = true
			closest = tempRecord.T
			record = tempRecord
		}
	}
	return hitAnything, record
}
