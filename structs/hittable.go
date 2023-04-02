package structs

type HitRef struct {
	T float64
	P, Normal Vec3
}

type Hittable interface {
	Hit(r *Ray, tMin float64, tMax float64) (bool, HitRef)
}