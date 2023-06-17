package structs


type DiffuseLight struct {
	Albedo *Vec3
	EmissionStrength float64
}

func NewDiffuseLight(albedo Vec3, emissionStrength float64) *DiffuseLight {
	return &DiffuseLight{
		Albedo: &albedo,
		EmissionStrength: emissionStrength,
	}
}

func (d *DiffuseLight) Scatter(rIn Ray, rec HitRef) (bool, *Ray, Vec3) {
	return false, nil, Vec3{}
}

func (d *DiffuseLight) Emitted(u, v float64, p Vec3) Vec3 {
	return d.Albedo.MulScalar(d.EmissionStrength)
}