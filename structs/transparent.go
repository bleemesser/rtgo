package structs

import (
	"math"
	"math/rand"
)

type Transparent struct {
	Albedo Vec3
	Index  float64
	Roughness float64
}

func NewTransparent(albedo Vec3, index float64, roughness float64) Transparent {
	return Transparent{Albedo: albedo, Index: index, Roughness: roughness}
}

func (t Transparent) refract(uv Vec3, n Vec3, etaiOverEtat float64) Vec3 {
	cosTheta := uv.MulScalar(-1).Dot(n)
	rOutParallel := uv.Add(n.MulScalar(cosTheta)).MulScalar(etaiOverEtat)
	rOutPerp := n.MulScalar(-1 * math.Sqrt(1 - rOutParallel.Dot(rOutParallel)))
	return rOutParallel.Add(rOutPerp)
}

func reflectance(cos float64, refIdx float64) float64 {
	r0 := (1 - refIdx) / (1 + refIdx)
	r0 = r0 * r0
	return r0 + (1-r0)*math.Pow((1-cos), 5)
}
// TRANSPARENT SCATTER
// func (t Transparent) Scatter(in Ray, rec HitRef) (bool, *Ray, Vec3) {
// 	attenuation := t.Albedo
// 	refractionRatio := 1.0 / t.Index

// 	unitDir := in.Direction.Normalize()
// 	cos := math.Min(unitDir.MulScalar(-1).Dot(rec.Normal), 1.0)
// 	sin := math.Sqrt(1.0 - cos*cos)

// 	cannotRefract := refractionRatio*sin > 1.0
// 	var dir Vec3
// 	if cannotRefract || reflectance(cos, refractionRatio) > rand.Float64() {
// 		dir = reflect(unitDir, rec.Normal)
// 	} else {
// 		dir = t.refract(unitDir, rec.Normal, refractionRatio)
// 	}

// 	scattered := NewRay(rec.P, dir)
// 	return true, scattered, attenuation
// }

// LAMBERTIAN SCATTER
// func (m Lambertian) Scatter(in Ray, rec HitRef) (bool, *Ray, Vec3) {
// 	scatterDir := rec.Normal.Add(RandomUnitVector())

// 	if scatterDir.IsNearZero() {
// 		scatterDir = rec.Normal
// 	}

// 	scattered := NewRay(rec.P, scatterDir)
// 	attenuation := m.Albedo.Value(rec.U, rec.V, rec.P)
// 	return true, scattered, attenuation
// }

// Transparent scatter but with roughness
func (t Transparent) Scatter(in Ray, rec HitRef) (bool, *Ray, Vec3) {
	attenuation := t.Albedo
	refractionRatio := 1.0 / t.Index

	unitDir := in.Direction.Normalize()
	cos := math.Min(unitDir.MulScalar(-1).Dot(rec.Normal), 1.0)
	sin := math.Sqrt(1.0 - cos*cos)

	cannotRefract := refractionRatio*sin > 1.0
	var dir Vec3
	if cannotRefract || reflectance(cos, refractionRatio) > rand.Float64() {
		dir = reflect(unitDir, rec.Normal)
	} else {
		dir = t.refract(unitDir, rec.Normal, refractionRatio)
	}

	scattered := NewRay(rec.P, dir.Add(RandomInUnitSphere().MulScalar(t.Roughness)))
	return true, scattered, attenuation
}

func (t Transparent) Emitted(u, v float64, p Vec3) Vec3 {
	return Vec3{}
}