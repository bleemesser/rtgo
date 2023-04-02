package structs

import "math"

type Sphere struct {
	Center Vec3
	Radius float64
}

func (s *Sphere) Hit(r *Ray, tMin float64, tMax float64) (bool, HitRef) {
	 o := r.Origin.Sub(s.Center)
	 a := r.Direction.Dot(r.Direction)
	 halfB := o.Dot(r.Direction)
	 c := o.Dot(o) - s.Radius*s.Radius
	 discriminant := halfB*halfB - a*c

	 rec := HitRef{}

	 if discriminant > 0 {
		t := (-halfB - math.Sqrt(discriminant)) / a
		if t < tMax && t > tMin {
			rec.T = t
			rec.P = r.PointAt(t)
			rec.Normal = rec.P.Sub(s.Center).DivScalar(s.Radius)
			return true, rec
		}
		t = (-halfB + math.Sqrt(discriminant)) / a
		if t < tMax && t > tMin {
			rec.T = t
			rec.P = r.PointAt(t)
			rec.Normal = rec.P.Sub(s.Center).DivScalar(s.Radius)
			return true, rec
		}

	 }
	 return false, rec
}