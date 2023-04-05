package utility

import (
	// "image"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"

	// "fmt"
	st "rtgo/structs"
)

func convertVal(px st.Vec3, exposure float64) color.RGBA64 {
	// unscaled values between 0 and 1 inclusive
	r := px.X
	g := px.Y
	b := px.Z

	scale := 1.0 / float64(exposure)
	r = math.Sqrt(scale * r)
	g = math.Sqrt(scale * g)
	b = math.Sqrt(scale * b)

	// scaled values between 0 and 65535 inclusive
	r = 65535 * st.Clamp(r, 0, 1)
	g = 65535 * st.Clamp(g, 0, 1)
	b = 65535 * st.Clamp(b, 0, 1)

	// gamma correction
	r = 65535 * st.Clamp(math.Pow(r/65535, 1/exposure), 0, 1)
	g = 65535 * st.Clamp(math.Pow(g/65535, 1/exposure), 0, 1)
	b = 65535 * st.Clamp(math.Pow(b/65535, 1/exposure), 0, 1)

	return color.RGBA64{uint16(r), uint16(g), uint16(b), 65535}

}

func ArrayToPng(img [][]st.Vec3, filename string, exposure float64) {
	imageData := image.NewRGBA(image.Rect(0, 0, len(img[0]), len(img)))
	height := len(img)

	for y := 0; y < height; y++ {
		for x := 0; x < len(img[0]); x++ {
			imageData.Set(x, height-y-1, convertVal(img[y][x], exposure))
		}
	}

	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	png.Encode(f, imageData)
}
