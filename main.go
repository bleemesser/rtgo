package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	st "rtgo/structs"

	// progress bar
	b "github.com/schollz/progressbar/v3"
)

const (
	ratio = 16.0 / 9.0
	width = 1280

	height          = int(width / ratio)
	aaSamples       = 100
	c               = 255.99
	maxDepth        = 100
	samplesPerPixel = 1
)

var (
	white = st.Vec3{X: 1.0, Y: 1.0, Z: 1.0}
	blue  = st.Vec3{X: 0.5, Y: 0.7, Z: 1.0}

	camera  = st.NewCamera(ratio, 1, 2.5)
	objects = []st.Hittable{
		st.NewSphere(st.Vec3{X: 0, Y: -103, Z: -20}, 100, st.NewLambertian(st.Vec3{X: 0.5, Y: 0.5, Z: 0.5})),
		// st.NewSphere(st.Vec3{X: 15, Y: 2.6, Z: -20}, 3),
		st.NewSphere(st.Vec3{X: 0, Y: 0, Z: -20}, 3, st.NewMetal(st.Vec3{X: 0, Y: 0.2, Z: 0.8}, 1)),
		st.NewSphere(st.Vec3{X: 5, Y: 0, Z: -20}, 1.5, st.NewMetal(st.Vec3{X: 0.8, Y: 0.0, Z: 0.5}, 0.6)),
	}

	world = st.World{Objects: objects}

	bar = b.Default(int64(width * height))
)

// need to do material properties
// need to do lights
// make it faster

func color(r *st.Ray, h st.Hittable, depth int) st.Vec3 {
	if depth <= 0 {
		return st.Vec3{}
	}
	if hit, rec := h.Hit(r, 0.001, math.MaxFloat64); hit {
		check, scattered, attenuation := rec.Mat.Scatter(*r, rec)
		if check {
			return attenuation.Mul(color(scattered, h, depth-1))
		}
		return st.Vec3{X: 0, Y: 0, Z: 0}
	}

	unitDirection := r.Direction.Normalize()
	t := 0.5 * (unitDirection.Y + 1.0)
	return white.MulScalar(1.0 - t).Add(blue.MulScalar(t))

}

func main() {
	f, err := os.Create("image.ppm")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer f.Close()

	fmt.Fprintf(f, "P3\n%d %d\n255\n", width, height)

	for j := height - 1; j >= 0; j-- {
		for i := 0; i < width; i++ {
			col := st.Vec3{}

			for s := 0; s < aaSamples; s++ {
				u := (float64(i) + rand.Float64()) / float64(width)
				v := (float64(j) + rand.Float64()) / float64(height)
				r := camera.GetRay(u, v)
				col = col.Add(color(&r, &world, maxDepth))
			}

			col = col.DivScalar(float64(aaSamples))

			st.WriteColor(f, col, samplesPerPixel)

			bar.Add(1)

		}
	}
}
