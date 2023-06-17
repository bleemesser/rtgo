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
	// height = 1920
	// IMAGE OPTIONS
	aaSamples = 16000
	maxDepth  = 50
	exposure  = 1 // (samples per pixel, default 1, lower is brighter)

	height = int(width / ratio)
	// width = int(height * ratio)

	// DIVIDE IMAGE INTO PARTS FOR PARALLEL PROCESSING
	partDiv = 24 // YOUR IMAGE HEIGHT AND WIDTH MUST BE EVENLY DIVISIBLE BY THIS NUMBER
	// 720p = 20, 1080p = 24, 1440p = 40, 2160p = 12 or 24
	// 400w = 5, 800w = 10

	// FILE SETTINGS
	outputAsPng = true
	fileName    = "out" // don't add the file extension
)

var (
	// DEFINE BACKGROUND GRADIENT
	white = st.Vec3{X: 1.0, Y: 1.0, Z: 1.0}
	blue  = st.Vec3{X: 0.5, Y: 0.7, Z: 1.0}

	// cornellRed   = st.Vec3{X: 0.65, Y: 0.05, Z: 0.05}
	// cornellWhite = st.Vec3{X: 0.73, Y: 0.73, Z: 0.73}
	// cornellGreen = st.Vec3{X: 0.12, Y: 0.45, Z: 0.15}

	// SET CAMERA FOV
	cameraUp       = st.Vec3{X: 0, Y: 1, Z: 0}
	cameraLookFrom = st.Vec3{X: -5, Y: 5, Z: -40}
	cameraLookAt   = st.Vec3{X: 4, Y: 2, Z: 10}
	vFov           = 20.0
	// cameraLookFrom         = st.Vec3{X: 380, Y: 278, Z: -800}
	// cameraLookAt           = st.Vec3{X: 278, Y: 278, Z: 0}
	// vFov                   = 40.0
	focusDist         = cameraLookFrom.Sub(cameraLookAt).Length()
	aperture  float64 = 0.04

	camera = st.NewCamera(cameraLookFrom, cameraLookAt, cameraUp, vFov, ratio, focusDist, aperture)

	backgroundColor       = st.Vec3{X: 0, Y: 0, Z: 0}
	useBackgroundGradient = false

	objects = []st.Hittable{
		// floor
		// st.NewSphere(st.Vec3{X: 0, Y: -100000, Z: -20}, 100000, st.NewMetal(st.Vec3{X: 1,Y: 1,Z: 1}, 0.7)), // floor
		// st.NewSphere(st.Vec3{X: 0, Y: 1, Z: 0}, 1, st.NewMetal(st.Vec3{X: 0.9, Y: 0.9, Z: 0.9}, 0)), // centered mirror sphere
		// st.NewRectangularPlane(st.Vec3{-20,-20,-20}, st.Vec3{-20,60,-20},st.Vec3{-20,-20,20}, st.NewDiffuseLight(st.Vec3{X: 1, Y: 1, Z: 1}, 6)), // large rectangle light
		// light sphere
		st.NewSphere(st.Vec3{X: 30, Y: 60, Z: 35}, 20, st.NewDiffuseLight(st.Vec3{X: 1, Y: 1, Z: 1}, 15)),
		st.NewSphere(st.Vec3{X: 30,Y: -60,Z: 35}, 20, st.NewDiffuseLight(st.Vec3{X: 1, Y: 0, Z: 1}, 15)),
		st.NewSphere(st.Vec3{X: -30,Y: 60,Z: 35}, 20, st.NewDiffuseLight(st.Vec3{X: 0, Y: 0, Z: 1}, 15)),
		// st.NewSphere(st.Vec3{X: -20, Y: 15, Z: 35}, 10, st.NewDiffuseLight(st.Vec3{X: 1, Y: 1, Z: 1}, 3)),

		// st.NewRectangularPrism(st.NewRectangularPlane(st.Vec3{X: 0,Y: 2,Z: 0}, st.Vec3{X: 1,Y: 2,Z: 0}, st.Vec3{X: 0,Y: 2,Z: 1}, st.NewLambertian(st.Vec3{X: 0.5, Y: 0.5, Z: 0.5})), 2),
		// st.NewSphere(st.Vec3{X: 0, Y: 25, Z: 0}, 1, st.NewDiffuseLight(st.Vec3{X: 0.99, Y: 0, Z: 0}, 2)),
	}
	// mat = st.NewLambertian(st.Vec3{X: 0.7, Y: 0.7, Z: 0.7})
	// triangles = ut.LoadOBJFile("knight.obj", mat)
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
	span := 15
	for a := -span; a < span; a++ {
		for b := -span; b < span; b++ {
			chooseMat := rand.Float64()
			center := st.Vec3{X: float64(a) + 0.9*rand.Float64(), Y: 0.2, Z: float64(b) + 0.9*rand.Float64()}
			radius := 0.2
			if center.Sub(st.Vec3{X: 4, Y: radius, Z: 0}).Length() > 0.9 {
				switch {
				case chooseMat < 0.24:
					albedo := st.Vec3{X: rand.Float64() * rand.Float64(), Y: rand.Float64() * rand.Float64(), Z: rand.Float64() * rand.Float64()}
					material = st.NewDiffuseLight(albedo, randomFloat(1.8, 4))
				case chooseMat < 0.3:
					albedo := st.Vec3{X: rand.Float64() * rand.Float64(), Y: rand.Float64() * rand.Float64(), Z: rand.Float64() * rand.Float64()}
					material = st.NewLambertian(albedo)
				case chooseMat < 0.6:
					albedo := st.Vec3{X: rand.Float64() * rand.Float64(), Y: rand.Float64() * rand.Float64(), Z: rand.Float64() * rand.Float64()}
					smoothness := randomFloat(0, 1)
					material = st.NewMetal(albedo, smoothness)
				default:
					albedo := st.Vec3{X: rand.Float64() * rand.Float64(), Y: rand.Float64() * rand.Float64(), Z: rand.Float64() * rand.Float64()}
					material = st.NewTransparent(albedo, 1.5)
				}
				objects = append(objects, st.NewSphere(center, radius, material))
			}
		}
	}
	return st.World{Objects: objects}
}

func main() {
	start := time.Now()

	triangles := ut.LoadOBJFile("crystals.obj", st.NewTransparent(st.Vec3{X: 1, Y: 1, Z: 1}, 1.5))
	// set camera lookat to the average of all the triangles
	var sum st.Vec3
	for _, tri := range triangles {
		sum = sum.Add(tri.GetPos())
	}
	cameraLookAt = sum.DivScalar(float64(len(triangles)))
	camera = st.NewCamera(cameraLookFrom, cameraLookAt, cameraUp, vFov, ratio, focusDist, aperture)
	objects = append(objects, triangles...)
	world := st.World{Objects: objects}
	// world := randomScene()
	world_bvh := st.NewBVHNode(world.Objects, 0, len(world.Objects), 0, 0)
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
	if aaSamples <= 0 {
		fmt.Println("0 or fewer aaSamples specified, exiting")
		os.Exit(0)
	}
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
						col = col.Add(color(&r, &world_bvh, maxDepth))
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
