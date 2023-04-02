package structs

type Material interface {
	Bounce(in Ray, hit HitRef) (bool, Ray)
	Color() Vec3
}

