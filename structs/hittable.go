package structs

type HitRef struct {
	T float64
	P, Normal Vec3
	Mat Material
	FrontFace bool
}

type Hittable interface {
	Hit(r *Ray, tMin float64, tMax float64) (bool, HitRef)
}

func (hr *HitRef) SetFaceNormal(r *Ray, outwardNormal Vec3) {
	hr.FrontFace = r.Direction.Dot(outwardNormal) < 0
	if hr.FrontFace {
		hr.Normal = outwardNormal
	} else {
		hr.Normal = outwardNormal.MulScalar(-1)
	}
}

