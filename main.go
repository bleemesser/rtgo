package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	st "rtgo/structs"
	"sync"
	"time"

	// progress bar
	b "github.com/schollz/progressbar/v3"
)

const (
	// SET IMAGE SIZE
	ratio = 16.0 / 9.0
	width = 3840
	// IMAGE OPTIONS
	aaSamples = 200
	maxDepth  = 40
	exposure  = 1 // (samples per pixel, default 1, lower is brighter)

	height = int(width / ratio)

	// DIVIDE IMAGE INTO PARTS FOR PARALLEL PROCESSING
	partDiv = 24 // YOUR IMAGE HEIGHT AND WIDTH MUST BE EVENLY DIVISIBLE BY THIS NUMBER
	// 1080p = 24, 1440p = 20, 2160p = 12 or 24
)

var (
	// DEFINE BACKGROUND GRADIENT
	white = st.Vec3{X: 1.0, Y: 1.0, Z: 1.0}
	blue  = st.Vec3{X: 0.5, Y: 0.7, Z: 1.0}

	// SET CAMERA FOV
	cameraUp               = st.Vec3{X: 0, Y: 1, Z: 0}
	cameraLookFrom         = st.Vec3{X: 10, Y: 0.9, Z: 0.8}
	cameraLookAt           = st.Vec3{X: 0, Y: 1, Z: 0}
	focusDist              = cameraLookFrom.Sub(cameraLookAt).Length()
	aperture       float64 = 0.04

	camera = st.NewCamera(cameraLookFrom, cameraLookAt, cameraUp, 20, ratio, focusDist, aperture)

	objects = []st.Hittable{
		st.NewSphere(st.Vec3{X: 0, Y: -100000, Z: -20}, 100000, st.NewLambertian(st.Vec3{X: 0.5, Y: 0.5, Z: 0.5})),
		st.NewSphere(st.Vec3{X: 0, Y: 1, Z: 0}, 1, st.NewMetal(st.Vec3{X: 0.9, Y: 0.9, Z: 0.9}, 0)),
	}

	world = randomScene()
)

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

func randomFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func randomScene() st.World {
	var material st.Material

	for a := -9; a < 9; a++ {
		for b := -9; b < 9; b++ {
			chooseMat := rand.Float64()
			center := st.Vec3{X: float64(a) + 0.9*rand.Float64(), Y: 0.2, Z: float64(b) + 0.9*rand.Float64()}

			if center.Sub(st.Vec3{X: 4, Y: 0.2, Z: 0}).Length() > 0.9 {
				switch {
				case chooseMat < 0.6:
					albedo := st.RandomVec(0, 1)
					material = st.NewLambertian(albedo)
				case chooseMat < 0.75:
					albedo := st.RandomVec(0, 1)
					smoothness := randomFloat(0.1, 1)
					material = st.NewMetal(albedo, smoothness)
				default:
					albedo := st.RandomVec(0, 1)
					material = st.NewTransparent(albedo, randomFloat(0.5, 3))
				}
				objects = append(objects, st.NewSphere(center, 0.2, material))
			}
		}
	}

	return st.World{Objects: objects}
}

func main() {
	start := time.Now()
	f, err := os.Create("image.ppm")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	fmt.Fprintf(f, "P3\n%d %d\n255\n", width, height)

	// 2d array of pixels as buf
	buf := make([][]st.Vec3, height)
	for i := range buf {
		buf[i] = make([]st.Vec3, width)
	}

	doneCh := make(chan bool)

	partHeight := height / partDiv
	partWidth := width / partDiv
	partArea := partHeight * partWidth
	numParts := partDiv * partDiv

	parts := make([]st.ImagePart, numParts)
	for i := 0; i < partDiv; i++ {
		for j := 0; j < partDiv; j++ {
			part := st.ImagePart{
				StartRow: j * partHeight,
				EndRow:   (j + 1) * partHeight,
				StartCol: i * partWidth,
				EndCol:   (i + 1) * partWidth,
				I:        i,
				J:        j,
			}
			parts = append(parts, part)
		}
	}

	maxConcurrentParts := 16
	activeParts := make(chan struct{}, maxConcurrentParts)
	for i := 0; i < maxConcurrentParts; i++ {
		activeParts <- struct{}{}
	}

	// print info about the image
	fmt.Println("\nImage size:", width, "x", height)
	fmt.Println("Samples per pixel:", aaSamples)
	fmt.Println("Rays: ", width*height*aaSamples)
	fmt.Println("Number of objects:", len(objects))
	fmt.Println("Exposure:", exposure)
	fmt.Println("Aperture:", aperture)
	fmt.Println("Focus distance:", focusDist)
	fmt.Println("Number of parts:", numParts)
	fmt.Println("Part size:", partWidth, "x", partHeight, "Area:", partArea)
	fmt.Println("Number of concurrent parts:", maxConcurrentParts)
	fmt.Println("\nRendering...")

	// progress bar
	bar := b.Default(int64(height * width))

	go func() {
		for range parts {
			<-doneCh
			activeParts <- struct{}{}
		}
	}()
	wg := sync.WaitGroup{}
	for _, part := range parts {
		<-activeParts
		wg.Add(1)
		go func(part st.ImagePart) {
			// fmt.Println("\nStarting part", part.I, part.J, part)
			for j := part.StartRow; j < part.EndRow; j++ {
				for i := part.StartCol; i < part.EndCol; i++ {
					col := st.Vec3{}
					for s := 0; s < aaSamples; s++ {
						u := (float64(i) + rand.Float64()) / float64(width)
						v := (float64(j) + rand.Float64()) / float64(height)
						r := camera.GetRay(u, v)
						col = col.Add(color(&r, &world, maxDepth))
					}
					col = col.DivScalar(float64(aaSamples))

					buf[j][i] = col
					bar.Add(1)
				}
			}
			// fmt.Println("\nFinished part", part.I, part.J, part)

			doneCh <- true
			wg.Done()
		}(part)
	}
	wg.Wait()
	// write pixels from buffer to file in correct order
	for j := height - 1; j >= 0; j-- {
		for i := 0; i < width; i++ {
			col := buf[j][i]
			st.WriteColor(f, col, exposure)
		}
	}

	fmt.Println("Render Time:", time.Since(start))

}

// for j := height - 1; j >= 0; j-- {
// 	for i := 0; i < width; i++ {
// 		col := st.Vec3{}

// 		for s := 0; s < aaSamples; s++ {
// 			u := (float64(i) + rand.Float64()) / float64(width)
// 			v := (float64(j) + rand.Float64()) / float64(height)
// 			r := camera.GetRay(u, v)
// 			col = col.Add(color(&r, &world, maxDepth))
// 		}

// 		col = col.DivScalar(float64(aaSamples))

// 		st.WriteColor(f, col, exposure)

// 		}
// 	}
