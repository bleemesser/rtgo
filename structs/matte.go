package structs

type Lambertian struct {
	Albedo Texture
}

func NewLambertian(texture Texture) Lambertian {
	return Lambertian{Albedo: texture}
}

func (m Lambertian) Scatter(in Ray, rec HitRef) (bool, *Ray, Vec3) {
	scatterDir := rec.Normal.Add(RandomUnitVector())

	if scatterDir.IsNearZero() {
		scatterDir = rec.Normal
	}

	scattered := NewRay(rec.P, scatterDir)
	attenuation := m.Albedo.Value(rec.U, rec.V, rec.P)
	return true, scattered, attenuation
}

func (m Lambertian) Emitted(u, v float64, p Vec3) Vec3 {
	return Vec3{}
}