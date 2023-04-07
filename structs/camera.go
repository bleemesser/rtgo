package structs

import (
	// "fmt"
	"math"
)


type Camera struct {
	Origin          Vec3
	LowerLeftCorner Vec3
	Horizontal      Vec3
	Vertical        Vec3
	LensRadius      float64
	Aperture        float64
	FocusDist       float64
}

func NewCamera(lookFrom, lookAt, vUp Vec3, vfov, aspectRatio, focusDist, aperture float64) *Camera {
	// vfov = vertical fov, degrees
	theta := vfov * math.Pi / 180
	halfHeight := math.Tan(theta / 2)
	height := halfHeight * 2
	width := aspectRatio * height

	w := lookFrom.Sub(lookAt).Normalize()
	u := vUp.Cross(w).Normalize()
	v := w.Cross(u)

	origin := lookFrom
	horizontal := u.MulScalar(width).MulScalar(focusDist)
	vertical := v.MulScalar(height).MulScalar(focusDist)
	lowerLeftCorner := origin.Sub(horizontal.DivScalar(2)).Sub(vertical.DivScalar(2)).Sub(w.MulScalar(focusDist))
	radius := aperture / 2

	return &Camera{
		Origin:          origin,
		LowerLeftCorner: lowerLeftCorner,
		Horizontal:      horizontal,
		Vertical:        vertical,
		LensRadius: radius,
		Aperture: aperture,
		FocusDist: focusDist,
	}
}

func (c *Camera) GetRay(u float64, v float64) Ray {
	// return Ray{c.Origin, c.LowerLeftCorner.Add(c.Horizontal.MulScalar(u).Add(c.Vertical.MulScalar(v).Sub(c.Origin)))}
	rand := RandomInUnitDisk().MulScalar(c.LensRadius)
	offset := rand.X * u + v * rand.Y

	return Ray{
		Origin: c.Origin.AddScalar(offset),
		Direction: c.LowerLeftCorner.Add(c.Horizontal.MulScalar(u)).Add(c.Vertical.MulScalar(v)).Sub(c.Origin).SubScalar(offset),
		U: u,
		V: v,
	}
}
