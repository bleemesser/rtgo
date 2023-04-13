package structs

// import "fmt"
// import "sort"
// import "math"

// type RectangularPrism struct {
// 	// Vertices
// 	A, B, C, D, E, F, G, H Vec3
// 	// Material
// 	Mat Material
// }
// func OrderPointsClockwise(points []Vec3) []Vec3 {
// 	// Compute the centroid of the points
// 	var centroid Vec3
// 	for _, p := range points {
// 		centroid = centroid.Add(p)
// 	}
// 	centroid = centroid.MulScalar(1.0 / float64(len(points)))

// 	// Compute the normal vector of the plane defined by the points
// 	v1 := points[0].Sub(points[1])
// 	v2 := points[1].Sub(points[2])
// 	normal := v1.Cross(v2).Normalize()
// 	if normal.Dot(centroid.Sub(points[0])) < 0 {
// 		normal = normal.MulScalar(-1)
// 	}

// 	// Choose a reference direction perpendicular to the normal vector
// 	refDir := normal.Cross(points[1].Sub(points[0]))

// 	// Compute the angle between each point and the reference direction
// 	angles := make([]float64, len(points))
// 	for i, p := range points {
// 		v := p.Sub(centroid)
// 		angle := math.Atan2(normal.Dot(refDir.Cross(v)), refDir.Dot(v))
// 		if angle < 0 {
// 			angle += 2 * math.Pi
// 		}
// 		angles[i] = angle
// 	}

// 	// Sort the points by their angles
// 	sortedPoints := make([]Vec3, len(points))
// 	copy(sortedPoints, points)
// 	sort.Slice(sortedPoints, func(i, j int) bool {
// 		return angles[i] < angles[j]
// 	})

// 	return sortedPoints
// }

// // define a rectangular prism using the minimum number of vertices needed to be able to have it at any orientation
// func NewRectangularPrism(r1 *RectangularPlane, x float64) *RectangularPrism {
// 	// r1 is the bottom plane, x is the height of the prism.
// 	// the prism is defined by the bottom plane and the top plane
// 	// which is the bottom plane translated by x in the direction of the normal of the bottom plane

// 	// sort the vertices of the bottom plane in clockwise order
// 	// so that the normal of every plane is pointing outwards

// 	points := []Vec3{r1.A, r1.B, r1.C, r1.D}
// 	points = OrderPointsClockwise(points)
// 	r1.A = points[0]
// 	r1.B = points[1]
// 	r1.C = points[2]
// 	r1.D = points[3]


// 	// get the normal of the bottom plane
// 	normal := r1.Normal()
// 	fmt.Println(r1.A, r1.B, r1.C, r1.D)
// 	fmt.Println(x)
// 	fmt.Println(normal)
// 	// translate the bottom plane by x in the direction of the normal
	
// 	r2 := RectangularPlane{}

// 	points2 := []Vec3{r1.A.Add(normal.MulScalar(x)), r1.B.Add(normal.MulScalar(x)), r1.C.Add(normal.MulScalar(x)), r1.D.Add(normal.MulScalar(x))}
// 	points2 = OrderPointsClockwise(points2)
// 	r2.A = points2[0]
// 	r2.B = points2[1]
// 	r2.C = points2[2]
// 	r2.D = points2[3]

// 	rect := RectangularPrism{A: r1.A, B: r1.B, C: r1.C, D: r1.D, E: r2.A, F: r2.B, G: r2.C, H: r2.D, Mat: r1.Mat}
// 	fmt.Println(rect.A, rect.B, rect.C, rect.D, rect.E, rect.F, rect.G, rect.H)
// 	return &rect
// }

// func (rect *RectangularPrism) Hit(r *Ray, tMin, tMax float64) (bool, HitRef) {
// 	// split the prism into 6 rectangular planes
// 	// and check if the ray hits any of them
// 	// if it does, return the closest hit

// 	// temporarily set the material of each rectangular plane to something different and identifiable
// 	planes := []*RectangularPlane{
// 		MakeRectangularPlane(rect.A, rect.B, rect.C, rect.D, rect.Mat),
// 		MakeRectangularPlane(rect.A, rect.B, rect.F, rect.E, rect.Mat),
// 		// MakeRectangularPlane(rect.B, rect.C, rect.G, rect.F, rect.Mat),
// 		// MakeRectangularPlane(rect.C, rect.D, rect.H, rect.G, rect.Mat),
// 		// MakeRectangularPlane(rect.D, rect.A, rect.E, rect.H, rect.Mat),
// 		// MakeRectangularPlane(rect.E, rect.F, rect.G, rect.H, rect.Mat),

// 	}

// 	var hit bool
// 	var hitRef HitRef

// 	for _, plane := range planes {

// 		hit, hitRef = plane.Hit(r, tMin, tMax)
// 		if hit {
// 			return true, hitRef
// 		}
// 	}

// 	return false, HitRef{}
// }

// func (r *RectangularPrism) GetPos() Vec3 {
// 	return r.A
// }
