package structs

import "math"

type Sphere struct {
	Center Vec3
	Radius float64
	Mat    Material
}

func NewSphere(center Vec3, radius float64, material Material) Sphere {
	return Sphere{Center: center, Radius: radius, Mat: material}
}

func (s Sphere) Hit(r *Ray, tMin float64, tMax float64) (bool, HitRef) {
	o := r.Origin.Sub(s.Center)
	a := r.Direction.Dot(r.Direction)
	halfB := o.Dot(r.Direction)
	c := o.Dot(o) - s.Radius*s.Radius
	discriminant := halfB*halfB - a*c

	rec := HitRef{Mat: s.Mat}

	if discriminant > 0 {
		t := (-halfB - math.Sqrt(discriminant)) / a
		if t < tMax && t > tMin {
			rec.T = t
			rec.P = r.PointAt(t)
			outwardNormal := (rec.P.Sub(s.Center)).DivScalar(s.Radius)
			rec.SetFaceNormal(r, outwardNormal)
			rec.Mat = s.Mat
			// rec.Normal = rec.P.Sub(s.Center).DivScalar(s.Radius)
			return true, rec
		}
		t = (-halfB + math.Sqrt(discriminant)) / a
		if t < tMax && t > tMin {
			rec.T = t
			rec.P = r.PointAt(t)
			outwardNormal := (rec.P.Sub(s.Center)).DivScalar(s.Radius)
			rec.SetFaceNormal(r, outwardNormal)
			rec.Mat = s.Mat
			return true, rec
		}

	}

	return false, rec
}

func (s Sphere) GetPos() Vec3 {
	return s.Center
}

func (s Sphere) BoundingBox(time0, time1 float64) (bool, AABB) {
	return true, AABB{
		Min: s.Center.Sub(Vec3{s.Radius, s.Radius, s.Radius}),
		Max: s.Center.Add(Vec3{s.Radius, s.Radius, s.Radius}),
	}
}