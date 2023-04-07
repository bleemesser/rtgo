package structs

type Metal struct {
	Albedo     Vec3
	Roughness float64
}

func NewMetal(albedo Vec3, roughness float64) Metal {
	return Metal{Albedo: albedo, Roughness: roughness}
}

func reflect(i Vec3, n Vec3) Vec3 {
	return i.Sub(n.MulScalar(2 * i.Dot(n)))
}

func (m Metal) Scatter(in Ray, rec HitRef) (bool, *Ray, Vec3) {
	reflected := reflect(in.Direction, rec.Normal)
	scattered := NewRay(rec.P, reflected.Add(RandomInUnitSphere().MulScalar(m.Roughness)))
	attenuation := m.Albedo
	if scattered.Direction.Dot(rec.Normal) > 0 {
		return true, scattered, attenuation
	}
	return false, scattered, attenuation // questionable
}

func (m Metal) Emitted(u, v float64, p Vec3) Vec3 {
	return Vec3{}
}
