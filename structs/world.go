package structs

type World struct {
	Objects []Hittable
}

func (w *World) Hit(r *Ray, tMin float64, tMax float64) (bool, HitRef) {
	hit := false
	firstHit := tMax
	ref := HitRef{}

	for _, obj := range w.Objects {
		h, rec := obj.Hit(r, tMin, firstHit)
		if h {
			hit = true
			firstHit = rec.T
			ref = rec
		}
	}

	return hit, ref
}

func (w *World) GetPos() Vec3 {
	return Vec3{}
}