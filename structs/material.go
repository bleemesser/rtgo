package structs

// import "math/rand"

type Material interface {
	Scatter(in Ray, hit HitRef) (bool, *Ray, Vec3)
	Emitted(u, v float64, p Vec3) (Vec3)
}

// type UniversalMaterial struct {
// 	Albedo Vec3
// 	Metalness float64
// 	Transparency float64
// 	Roughness float64
// 	Emission float64
// }

// func NewUniversalMaterial(albedo Vec3, metalness, transparency, roughness, emission float64) UniversalMaterial {
// 	return UniversalMaterial{Albedo: albedo, Metalness: metalness, Transparency: transparency, Roughness: roughness, Emission: emission}
// }

// func (m UniversalMaterial) Scatter(in Ray, hit HitRef) (bool, *Ray, Vec3) {
// 	// use the percentage of each property to determine which material to use
// 	if m.Emission