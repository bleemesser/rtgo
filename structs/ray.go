package structs

type Ray struct {
	Origin, Direction Vec3
}

func (r Ray) PointAt(t float64) Vec3 {
	return r.Origin.Add(r.Direction.MulScalar(t))
} // p(t) = A + t*B


