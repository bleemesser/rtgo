package structs

import (
	"fmt"
	"math"
	"math/rand"
	"os"
)

type Vec3 struct {
	X, Y, Z float64
}

func (v Vec3) Add(v2 Vec3) Vec3 {
	return Vec3{v.X + v2.X, v.Y + v2.Y, v.Z + v2.Z}
}

func (v Vec3) Sub(v2 Vec3) Vec3 {
	return Vec3{v.X - v2.X, v.Y - v2.Y, v.Z - v2.Z}
}

func (v Vec3) Mul(v2 Vec3) Vec3 {
	return Vec3{v.X * v2.X, v.Y * v2.Y, v.Z * v2.Z}
}

func (v Vec3) Div(v2 Vec3) Vec3 {
	return Vec3{v.X / v2.X, v.Y / v2.Y, v.Z / v2.Z}
}

func (v Vec3) MulScalar(s float64) Vec3 {
	return Vec3{v.X * s, v.Y * s, v.Z * s}
}

func (v Vec3) DivScalar(s float64) Vec3 {
	return Vec3{v.X / s, v.Y / s, v.Z / s}
}

func (v Vec3) AddScalar(s float64) Vec3 {
	return Vec3{v.X + s, v.Y + s, v.Z + s}
}

func (v Vec3) SubScalar(s float64) Vec3 {
	return Vec3{v.X - s, v.Y - s, v.Z - s}
}

func (v Vec3) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func (v Vec3) Normalize() Vec3 {
	return v.DivScalar(v.Length())
}

func (v Vec3) Dot(v2 Vec3) float64 {
	return v.X*v2.X + v.Y*v2.Y + v.Z*v2.Z
}

func RandomInUnitSphere() Vec3 {
	p := Vec3{}
	safetyIterations := 100
	min := -1.0
	max := 1.0

	for safetyIterations > 0 {
		p = Vec3{rand.Float64()*(max-min) + min, rand.Float64()*(max-min) + min, rand.Float64()*(max-min) + min}
		// fmt.Println(p.X, p.Y, p.Z)
		if p.Length() < 1.0 {
			break
		}
		safetyIterations--
	}
	if safetyIterations == 0 {
		p = Vec3{0, 0, 0}
	}
	return p
}

func clamp(x float64, min float64, max float64) float64 {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}

func WriteColor(f *os.File, col Vec3, samplesPerPixel int) {
	r := col.X
	g := col.Y
	b := col.Z

	scale := 1.0 / float64(samplesPerPixel)
	r = math.Sqrt(scale * r)
	g = math.Sqrt(scale * g)
	b = math.Sqrt(scale * b)

	ir := int(256 * clamp(r, 0.0, 0.999))
	ig := int(256 * clamp(g, 0.0, 0.999))
	ib := int(256 * clamp(b, 0.0, 0.999))

	fmt.Fprintf(f, "%d %d %d\n", ir, ig, ib)
}

func RandomUnitVector() Vec3 {
	return RandomInUnitSphere().Normalize()
}

func (v *Vec3) IsNearZero() bool {
	s := math.Pow10(-8)
	X := v.X
	Y := v.Y
	Z := v.Z
	if math.Abs(X) < s && math.Abs(Y) < s && math.Abs(Z) < s {
		return true
	}
	return false
}
