package structs

type Metal struct {
	Albedo     Vec3
	Smoothness float64
}

func NewMetal(c Vec3, s float64) Metal {
	return Metal{Albedo: c, Smoothness: s}
}

func reflect(i Vec3, n Vec3) Vec3 {
	return i.Sub(n.MulScalar(2 * i.Dot(n)))
}

func (m Metal) Scatter(in Ray, rec HitRef) (bool, *Ray, Vec3) {
	reflected := reflect(in.Direction, rec.Normal)
	scattered := NewRay(rec.P, reflected)
	attenuation := m.Albedo
	if scattered.Direction.Dot(rec.Normal) > 0 {
		return true, scattered, attenuation
	}
	return false, scattered, attenuation // questionable
}
