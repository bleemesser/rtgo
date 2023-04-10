package structs

type XYRect struct {
	X0, X1, Y0, Y1, K float64
	Mat               Material
}

func NewXYRect(x0, x1, y0, y1, k float64, mat Material) *XYRect {
	return &XYRect{
		X0:  x0,
		X1:  x1,
		Y0:  y0,
		Y1:  y1,
		K:   k,
		Mat: mat,
	}
}

func (rect *XYRect) Hit(r *Ray, tMin float64, tMax float64) (bool, HitRef) {
	t := (rect.K - r.Origin.Z) / r.Direction.Z
	if t < tMin || t > tMax {
		return false, HitRef{}
	}
	x := r.Origin.X + t*r.Direction.X
	y := r.Origin.Y + t*r.Direction.Y

	if x < rect.X0 || x > rect.X1 || y < rect.Y0 || y > rect.Y1 {
		return false, HitRef{}
	}
	rec := HitRef{Mat: rect.Mat}
	rec.U = (x - rect.X0) / (rect.X1 - rect.X0)
	rec.V = (y - rect.Y0) / (rect.Y1 - rect.Y0)
	rec.T = t
	outwardNormal := Vec3{0, 0, 1}
	rec.SetFaceNormal(r, outwardNormal)
	rec.Mat = rect.Mat
	rec.P = r.PointAt(t)
	return true, rec
}

func (rect *XYRect) GetPos() Vec3 {
	return Vec3{}
}

type XZRect struct {
	X0, X1, Z0, Z1, K float64
	Mat               Material
}

func NewXZRect(x0, x1, z0, z1, k float64, mat Material) *XZRect {
	return &XZRect{
		X0:  x0,
		X1:  x1,
		Z0:  z0,
		Z1:  z1,
		K:   k,
		Mat: mat,
	}
}

func (rect *XZRect) Hit(r *Ray, tMin float64, tMax float64) (bool, HitRef) {
	t := (rect.K - r.Origin.Y) / r.Direction.Y
	if t < tMin || t > tMax {
		return false, HitRef{}
	}
	x := r.Origin.X + t*r.Direction.X
	z := r.Origin.Z + t*r.Direction.Z

	if x < rect.X0 || x > rect.X1 || z < rect.Z0 || z > rect.Z1 {
		return false, HitRef{}
	}
	rec := HitRef{Mat: rect.Mat}
	rec.U = (x - rect.X0) / (rect.X1 - rect.X0)
	rec.V = (z - rect.Z0) / (rect.Z1 - rect.Z0)
	rec.T = t
	outwardNormal := Vec3{0, 1, 0}
	rec.SetFaceNormal(r, outwardNormal)
	rec.Mat = rect.Mat
	rec.P = r.PointAt(t)
	return true, rec
}

func (rect *XZRect) GetPos() Vec3 {
	return Vec3{}
}

type YZRect struct {
	Y0, Y1, Z0, Z1, K float64
	Mat               Material
}

func NewYZRect(y0, y1, z0, z1, k float64, mat Material) *YZRect {
	return &YZRect{
		Y0:  y0,
		Y1:  y1,
		Z0:  z0,
		Z1:  z1,
		K:   k,
		Mat: mat,
	}
}

func (rect *YZRect) Hit(r *Ray, tMin float64, tMax float64) (bool, HitRef) {
	t := (rect.K - r.Origin.X) / r.Direction.X
	if t < tMin || t > tMax {
		return false, HitRef{}
	}
	y := r.Origin.Y + t*r.Direction.Y
	z := r.Origin.Z + t*r.Direction.Z

	if y < rect.Y0 || y > rect.Y1 || z < rect.Z0 || z > rect.Z1 {
		return false, HitRef{}
	}
	rec := HitRef{Mat: rect.Mat}
	rec.U = (y - rect.Y0) / (rect.Y1 - rect.Y0)
	rec.V = (z - rect.Z0) / (rect.Z1 - rect.Z0)
	rec.T = t
	outwardNormal := Vec3{1, 0, 0}
	rec.SetFaceNormal(r, outwardNormal)
	rec.Mat = rect.Mat
	rec.P = r.PointAt(t)
	return true, rec
}

func (rect *YZRect) GetPos() Vec3 {
	return Vec3{}
}
