package structs

import (
	"image"
	"image/png"
	"os"
	"fmt"
)


type Texture interface {
	Value(u float64, v float64, p Vec3) Vec3
}

type SolidColor struct {
	Color Vec3
}

func NewSolidColor(c Vec3) (*SolidColor) {
	return &SolidColor{Color: c}
}

func (s *SolidColor) Value(u float64, v float64, p Vec3) Vec3 {
	return s.Color
}

type PngTexture struct {
	Image image.Image
}

func NewPngTexture(path string) (*PngTexture) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer file.Close()
	img, err := png.Decode(file)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &PngTexture{Image: img}
}

func (p *PngTexture) Value(u float64, v float64, p1 Vec3) Vec3 {
	u = Clamp(u, 0, 1)
	v = 1 - Clamp(v, 0, 1)

	i := int(u * float64(p.Image.Bounds().Max.X))
	j := int(v * float64(p.Image.Bounds().Max.Y))

	if i >= p.Image.Bounds().Max.X {
		i = p.Image.Bounds().Max.X - 1
	}
	if j >= p.Image.Bounds().Max.Y {
		j = p.Image.Bounds().Max.Y - 1
	}

	r, g, b, _ := p.Image.At(i, j).RGBA()
	return Vec3{float64(r) / 65535, float64(g) / 65535, float64(b) / 65535}
}