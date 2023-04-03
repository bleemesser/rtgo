package structs

type Material interface {
	Scatter(in Ray, hit HitRef) (bool, *Ray, Vec3)
}
