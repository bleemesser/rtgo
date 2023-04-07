package structs

type Material interface {
	Scatter(in Ray, hit HitRef) (bool, *Ray, Vec3)
	Emitted(u, v float64, p Vec3) (Vec3)
}