package structs

type AABB struct {
	Min, Max Vec3
}

func min(a, b float64) float64 {
	if a < b {
		return a
	} else {
		return b
	}
}
func max(a, b float64) float64 {
	if a > b {
		return a
	} else {
		return b
	}
}

func (a AABB) Hit(r *Ray, tMin float64, tMax float64) (bool, HitRef) {
	for i := 0; i < 3; i++ {
		d := 1.0 / r.Direction.GetItemByIndex(i)
		t0 := (a.Min.GetItemByIndex(i) - r.Origin.GetItemByIndex(i)) * d
		t1 := (a.Max.GetItemByIndex(i) - r.Origin.GetItemByIndex(i)) * d
		if d < 0.0 {
			t0, t1 = t1, t0
		}
		if t0 > tMin {
			tMin = t0
		}
		if t1 < tMax {
			tMax = t1
		}
		if tMax <= tMin {
			return false, HitRef{}
		}
	}
	return true, HitRef{}
}

func SurroundingBox(box0, box1 AABB) AABB {
	small := Vec3{
		min(box0.Min.X, box1.Min.X),
		min(box0.Min.Y, box1.Min.Y),
		min(box0.Min.Z, box1.Min.Z),
	}
	big := Vec3{
		max(box0.Max.X, box1.Max.X),
		max(box0.Max.Y, box1.Max.Y),
		max(box0.Max.Z, box1.Max.Z),
	}
	return AABB{small, big}
}
