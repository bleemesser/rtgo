package structs

// type RectangularPlane struct {
// 	// Vertices
// 	A, B, C, D Vec3
// 	// Material
// 	Mat Material
// }

// func NewRectangularPlane(a, b, c Vec3, mat Material) *RectangularPlane {
// 	d := a.Add(c.Sub(b))
// 	points := []Vec3{a, b, c, d}
// 	points = OrderPointsClockwise(points)
// 	a = points[0]
// 	b = points[1]
// 	c = points[2]
// 	d = points[3]

// 	return &RectangularPlane{A: a, B: b, C: c, D: d, Mat: mat}
// }

// func MakeRectangularPlane(a,b,c,d Vec3, mat Material) *RectangularPlane {
// 	points := []Vec3{a, b, c, d}
// 	points = OrderPointsClockwise(points)
// 	a = points[0]
// 	b = points[1]
// 	c = points[2]
// 	d = points[3]
// 	return &RectangularPlane{A: a, B: b, C: c, D: d, Mat: mat}
// }

// func (rect *RectangularPlane) Hit(r *Ray, tMin, tMax float64) (bool, HitRef) {
// 	tri1 := NewTriangle(rect.A, rect.B, rect.C, rect.Mat)
// 	tri2 := NewTriangle(rect.A, rect.C, rect.D, rect.Mat)

// 	hit1, hitRef1 := tri1.Hit(r, tMin, tMax)
// 	hit2, hitRef2 := tri2.Hit(r, tMin, tMax)

// 	if hit1 && hit2 {
// 		if hitRef1.T < hitRef2.T {
// 			return true, hitRef1
// 		} else {
// 			return true, hitRef2
// 		}
// 	}
// 	if hit1 {
// 		return true, hitRef1
// 	}
// 	if hit2 {
// 		return true, hitRef2
// 	}
// 	return false, HitRef{}


// }

// func (r *RectangularPlane) GetPos() Vec3 {
// 	return r.A
// }

// func (plane *RectangularPlane) Normal() Vec3 {
//     return plane.B.Sub(plane.A).Cross(plane.C.Sub(plane.A)).Normalize()
// }