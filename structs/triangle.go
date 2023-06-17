package structs

type Triangle struct {
	A, B, C Vec3
	Mat     Material
}

func NewTriangle(a, b, c Vec3, mat Material) *Triangle {
	return &Triangle{A: a, B: b, C: c, Mat: mat}
}

func (tri *Triangle) Hit(r *Ray, tMin, tMax float64) (bool, HitRef) {	
	// Möller–Trumbore intersection algorithm
	edge1 := tri.B.Sub(tri.A)
	edge2 := tri.C.Sub(tri.A)
	h := r.Direction.Cross(edge2)
	a := edge1.Dot(h)
	if a > -0.00001 && a < 0.00001 {
		return false, HitRef{}
	}
	f := 1 / a
	s := r.Origin.Sub(tri.A)
	u := f * s.Dot(h)
	if u < 0 || u > 1 {
		return false, HitRef{}
	}
	q := s.Cross(edge1)
	v := f * r.Direction.Dot(q)
	if v < 0 || u+v > 1 {
		return false, HitRef{}
	}
	// At this stage we can compute t to find out where the intersection point is on the line.
	t := edge2.Dot(q) * f
	if t > tMin && t < tMax {
		p := r.PointAt(t)
		// Interpolate normal
		normal := tri.A.Sub(p).Cross(tri.B.Sub(p)).Normalize()
		return true, HitRef{T: t, P: p, Normal: normal, Mat: tri.Mat}
	}
	return false, HitRef{}
}

func (t *Triangle) GetPos() Vec3 {
	return t.A
}

func (t *Triangle) BoundingBox(time0, time1 float64) (bool, AABB) { // SKETCHY
	return true, AABB{
		Min: Vec3{
			X: min(min(t.A.X, t.B.X), t.C.X),
			Y: min(min(t.A.Y, t.B.Y), t.C.Y),
			Z: min(min(t.A.Z, t.B.Z), t.C.Z),
		},
		Max: Vec3{
			X: max(max(t.A.X, t.B.X), t.C.X),
			Y: max(max(t.A.Y, t.B.Y), t.C.Y),
			Z: max(max(t.A.Z, t.B.Z), t.C.Z),
		},
	}
}