package structs

type HitRef struct {
	T float64
	P, Normal Vec3
	Mat Material
	FrontFace bool
	U, V float64
}

type Hittable interface {
	Hit(r *Ray, tMin float64, tMax float64) (bool, HitRef)
	GetPos() Vec3
	BoundingBox(time0, time1 float64) (bool, AABB)
}

func (hr *HitRef) SetFaceNormal(r *Ray, outwardNormal Vec3) {
	hr.FrontFace = r.Direction.Dot(outwardNormal) < 0
	if hr.FrontFace {
		hr.Normal = outwardNormal
	} else {
		hr.Normal = outwardNormal.MulScalar(-1)
	}
}

