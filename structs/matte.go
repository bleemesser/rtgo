package structs

type Lambertian struct {
	Albedo Vec3
}

func NewLambertian(c Vec3) Lambertian {
	return Lambertian{Albedo: c}
}

func (m Lambertian) Scatter(in Ray, rec HitRef) (bool, *Ray, Vec3) {
	scatterDir := rec.Normal.Add(RandomUnitVector())

	if scatterDir.IsNearZero() {
		scatterDir = rec.Normal
	}

	scattered := NewRay(rec.P, scatterDir)
	attenuation := &m.Albedo
	return true, scattered, *attenuation
}
