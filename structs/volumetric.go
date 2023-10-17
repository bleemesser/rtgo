package structs

// import (
// 	"math"
// 	"math/rand"
// )

// type isotropic struct {
// 	albedo Texture
// }

// func (i *isotropic) Scatter(rIn *Ray, rec *HitRef) (bool, *Ray, Vec3) {
// 	scattered := Ray{
// 		Origin: rec.P,
// 		Direction: RandomInUnitSphere(),
// 		U: rIn.U,
// 		V: rIn.V,
// 	}
// 	attenuation := i.albedo.Value(rec.U, rec.V, rec.P)
// 	return true, &scattered, attenuation
// }



// type Volumetric struct {
// 	Boundary Hittable
// 	Density float64
// 	Mat Material
// }

// func (v *Volumetric) Hit(r *Ray, tMin float64, tMax float64) (bool, *HitRef) {
// 	var rec1, rec2 *HitRef
// 	if hit1, rec1 := v.Boundary.Hit(r, -math.MaxFloat64, math.MaxFloat64); !hit1 {
// 		return false, nil
// 	}

// 	if hit2, rec2 := v.Boundary.Hit(r, rec1.T+0.0001, math.MaxFloat64); !hit2 {
// 		return false, nil
// 	}

// 	if rec1.T < tMin {
// 		rec1.T = tMin
// 	}

// 	if rec2.T > tMax {
// 		rec2.T = tMax
// 	}

// 	if rec1.T >= rec2.T {
// 		return false, nil
// 	}

// 	if rec1.T < 0 {
// 		rec1.T = 0
// 	}

// 	distanceInsideBoundary := (rec2.T - rec1.T) * r.Direction.Length()
// 	hitDistance := -(1 / v.Density) * math.Log(rand.Float64())

// 	if hitDistance > distanceInsideBoundary {
// 		return false, nil
// 	}

// 	t := rec1.T + hitDistance / r.Direction.Length()
// 	p := r.PointAt(t)

// 	return true, &HitRef{
// 		T: t,
// 		P: p,
// 		Normal: Vec3{1, 0, 0},
// 		FrontFace: true,
// 		Mat: v.Mat,
// 	}
// }