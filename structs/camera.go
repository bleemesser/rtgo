package structs

type Camera struct {
	LowerLeftCorner, Horizontal, Vertical, Origin Vec3
}

func NewCamera(ratio, height, focalLength float64) *Camera {
	origin := Vec3{0, 0, 0}
	horizontal := Vec3{ratio, 0, 0}
	vertical := Vec3{0, height, 0}
	lowerLeftCorner := origin.Sub(horizontal.DivScalar(2)).Sub(vertical.DivScalar(2)).Sub(Vec3{0, 0, focalLength})

	return &Camera{lowerLeftCorner, horizontal, vertical, origin}

}

func (c *Camera) GetRay(u float64, v float64) Ray {
	pos := c.Horizontal.MulScalar(u).Add(c.Vertical.MulScalar(v))
	dir := c.LowerLeftCorner.Add(pos)
	return Ray{c.Origin, dir}
}