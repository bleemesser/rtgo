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

func (w *World) BoundingBox(time0, time1 float64) (bool, AABB) {
	if len(w.Objects) == 0 {
		return false, AABB{}
	}
	var outbox AABB
	firstBox := true

	for _, obj := range w.Objects {
		hit, b := obj.BoundingBox(time0, time1)
		if !hit {
			return false, AABB{}
		}
		if firstBox {
			outbox = b
		} else {
			outbox = SurroundingBox(outbox, b)
		}
		firstBox = false
	}
	return true, outbox
}