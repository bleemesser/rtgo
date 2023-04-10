package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	st "rtgo/structs"
	ut "rtgo/utility"
	"sync"
	"time"

	// progress bar
	b "github.com/schollz/progressbar/v3"
)

const (
	// SET IMAGE SIZE
	ratio = 16.0 / 9.0
	width = 1920
	// IMAGE OPTIONS
	aaSamples = 1000
	maxDepth  = 50
	exposure  = 1 // (samples per pixel, default 1, lower is brighter)

	height = int(width / ratio)

	// DIVIDE IMAGE INTO PARTS FOR PARALLEL PROCESSING
	partDiv = 24 // YOUR IMAGE HEIGHT AND WIDTH MUST BE EVENLY DIVISIBLE BY THIS NUMBER
	// 720p = 20, 1080p = 24, 1440p = 20, 2160p = 12 or 24
	// 400w = 5, 800w = 10

	// FILE SETTINGS
	outputAsPng = true
	fileName    = "out" // don't add the file extension
)

var (
	// DEFINE BACKGROUND GRADIENT
	white = st.Vec3{X: 1.0, Y: 1.0, Z: 1.0}
	blue  = st.Vec3{X: 0.5, Y: 0.7, Z: 1.0}

	cornellRed   = st.Vec3{X: 0.65, Y: 0.05, Z: 0.05}
	cornellWhite = st.Vec3{X: 0.73, Y: 0.73, Z: 0.73}
	cornellGreen = st.Vec3{X: 0.12, Y: 0.45, Z: 0.15}

	// SET CAMERA FOV
	cameraUp = st.Vec3{X: 0, Y: 1, Z: 0}
	// cameraLookFrom         = st.Vec3{X: 10, Y: 0.9, Z: 0.8}
	// cameraLookAt           = st.Vec3{X: 0, Y: 1, Z: 0}
	// vFov                   = 20.0
	cameraLookFrom         = st.Vec3{X: 380, Y: 278, Z: -800}
	cameraLookAt           = st.Vec3{X: 278, Y: 278, Z: 0}
	vFov                   = 40.0
	focusDist              = cameraLookFrom.Sub(cameraLookAt).Length()
	aperture       float64 = 0.04

	camera = st.NewCamera(cameraLookFrom, cameraLookAt, cameraUp, vFov, ratio, focusDist, aperture)

	backgroundColor       = st.Vec3{X: 0, Y: 0, Z: 0}
	useBackgroundGradient = false

	objects = []st.Hittable{
		// st.NewSphere(st.Vec3{X: 0, Y: -100000, Z: -20}, 100000, st.NewLambertian(st.Vec3{X: 0.5, Y: 0.5, Z: 0.5})), // floor
		// st.NewSphere(st.Vec3{X: 0, Y: 1, Z: 0}, 1, st.NewMetal(st.Vec3{X: 0.9, Y: 0.9, Z: 0.9}, 0)), // centered mirror sphere
		// st.NewXYRect(-20, 20, 0, 30, -20, st.NewDiffuseLight(st.Vec3{X: 1, Y: 1, Z: 1}, 6)), // large rectangle light

		// st.NewYZRect(0, 555, 0, 555, 555, st.NewLambertian(cornellGreen)), // left wall
		// st.NewYZRect(0, 555, 0, 555, 0, st.NewLambertian(cornellRed)),     // right wall
		// st.NewXZRect(0, 555, 0, 555, 0, st.NewLambertian(cornellWhite)),   // ceiling
		// st.NewXZRect(0, 555, 0, 555, 555, st.NewLambertian(cornellWhite)), // floor
		// // define the back wall as a white diffuse material
		// st.NewXYRect(0, 555, 0, 555, 555, st.NewLambertian(cornellWhite)), // back wall

		// st.NewXZRect(213, 343, 227, 332, 554, st.NewDiffuseLight(st.Vec3{X: 1, Y: 1, Z: 1}, 15)), // light

		// lambertian sphere in the box
		st.NewSphere(st.Vec3{X: 190, Y: 90, Z: 190}, 90, st.NewLambertian(st.Vec3{X: 0.73, Y: 0.73, Z: 0.73})),

		// position a light sphere up and to the left of the box, out of frame
		st.NewSphere(st.Vec3{X: 300, Y: 800, Z: -300}, 90, st.NewDiffuseLight(st.Vec3{X: 1, Y: 1, Z: 1}, 100)),

	}
	// rectangular prisms are a list of planar faces so need to be defined and added to the world separately
	rectangularPrisms = []st.RectangularPrism{
		// create two rectangular prisms inside the box
		// st.NewRectangularPrism(st.Vec3{X: 130, Y: 0, Z: 65}, st.Vec3{X: 295, Y: 165, Z: 230}, st.NewTransparent(blue, 1.5)),
		// st.NewRectangularPrism(st.Vec3{X: 265, Y: 0, Z: 295}, st.Vec3{X: 430, Y: 330, Z: 460}, st.NewLambertian(st.Vec3{X: 0.73, Y: 0.73, Z: 0.73})),

		// left wall
		st.NewRectangularPrism(st.Vec3{X: -20, Y: 0, Z: 0}, st.Vec3{X: 0, Y: 555, Z: 555}, st.NewLambertian(cornellGreen)),
		// right wall
		st.NewRectangularPrism(st.Vec3{X: 555, Y: 0, Z: 0}, st.Vec3{X: 575, Y: 555, Z: 555}, st.NewLambertian(cornellRed)),
		// ceiling
		st.NewRectangularPrism(st.Vec3{X: 0, Y: 555, Z: 0}, st.Vec3{X: 555, Y: 575, Z: 555}, st.NewLambertian(cornellWhite)),
		// floor
		st.NewRectangularPrism(st.Vec3{X: 0, Y: -20, Z: 0}, st.Vec3{X: 555, Y: 0, Z: 555}, st.NewLambertian(cornellWhite)),
		// back wall
		st.NewRectangularPrism(st.Vec3{X: 0, Y: 0, Z: 555}, st.Vec3{X: 555, Y: 555, Z: 575}, st.NewLambertian(cornellWhite)),
		// light
		st.NewRectangularPrism(st.Vec3{X: 213, Y: 554, Z: 227}, st.Vec3{X: 343, Y: 555, Z: 332}, st.NewDiffuseLight(st.Vec3{X: 1, Y: 1, Z: 1}, 15)),

	}
)

func color(r *st.Ray, h st.Hittable, depth int) st.Vec3 {
	if depth <= 0 {
		return st.Vec3{}
	}
	hit, rec := h.Hit(r, 0.001, math.MaxFloat64)
	if !hit {
		if useBackgroundGradient {
			unitDirection := r.Direction.Normalize()
			t := 0.5 * (unitDirection.Y + 1.0)
			return white.MulScalar(1.0 - t).Add(blue.MulScalar(t))
		}
		return backgroundColor
	}
	scatter, scattered, attenuation := rec.Mat.Scatter(*r, rec)
	emitted := rec.Mat.Emitted(r.U, r.V, rec.P)
	if !scatter {
		return emitted
	}
	return emitted.Add(attenuation.Mul(color(scattered, h, depth-1)))
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
				case chooseMat < 0.15:
					material = st.NewDiffuseLight(st.RandomVec(0, 1), randomFloat(1, 3))
				case chooseMat < 0.3:
					albedo := st.RandomVec(0, 1)
					material = st.NewLambertian(albedo)
				case chooseMat < 0.5:
					albedo := st.RandomVec(0, 1)
					smoothness := randomFloat(0, 1)
					material = st.NewMetal(albedo, smoothness)
				default:
					albedo := st.RandomVec(0.1, 1)
					material = st.NewTransparent(albedo, 1.5)
				}
				objects = append(objects, st.NewSphere(center, 0.2, material))
			}
		}
	}
	return st.World{Objects: objects}
}

func main() {
	start := time.Now()

	// finish creating the world by adding any rectangular prisms defined in the scene
	for _, prism := range rectangularPrisms {
		objects = st.CreateRectangularPrism(objects, prism)
	}
	world := st.World{Objects: objects}
	// world = randomScene()

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
	// add the max number of parts to active parts initially
	activeParts := make(chan bool, maxConcurrentParts)
	for i := 0; i < maxConcurrentParts; i++ {
		activeParts <- true
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
	if height*width%partArea != 0 {
		// in red, print a warning
		fmt.Println("\033[31mWarning: image size is not divisible by part size! There will be missing pixels!\033[0m")
	}
	fmt.Println("\nRendering...")

	// progress bar
	bar := b.Default(int64(height * width))

	// asynchronously, wait for there to be a value in doneCh for each part, then add a new part into active parts
	go func() {
		for range parts {
			<-doneCh
			activeParts <- true
		}
	}()

	wg := sync.WaitGroup{}
	for _, part := range parts {
		// wait for a signal to be in active parts, if there is remove it
		<-activeParts
		// add thread to waitgroup
		wg.Add(1)
		// asynchronously compute the pixels in that part
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
					// add output to image buffer
					buf[j][i] = col
					bar.Add(1)
				}
			}
			// fmt.Println("\nFinished part", part.I, part.J, part)
			// signal that the part is done, allowing another part to be dispatched
			doneCh <- true
			wg.Done()
		}(part)
	}
	// wait for all threads to finish, without this some pixels are left blank
	wg.Wait()
	// write pixels from buffer to file in correct order
	if outputAsPng {
		ut.ArrayToPng(buf, fileName+".png", exposure)
	} else {
		f, err := os.Create(fileName + ".ppm")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer f.Close()

		fmt.Fprintf(f, "P3\n%d %d\n255\n", width, height)
		for j := height - 1; j >= 0; j-- {
			for i := 0; i < width; i++ {
				col := buf[j][i]
				st.WriteColor(f, col, exposure)
			}
		}
	}
	// convert image to png instead

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
