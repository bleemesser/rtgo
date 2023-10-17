package structs

type Quad struct {
	tri1, tri2 *Triangle
	Mat        Material
}

func NewQuad(a, b, c, d Vec3, mat Material) *Quad {	
	tri1 := NewTriangle(b, c, a, mat)
	tri2 := NewTriangle(a, c, d, mat)
	return &Quad{tri1: tri1, tri2: tri2, Mat: mat}
}

func (q *Quad) Hit(r *Ray, tMin, tMax float64) (bool, HitRef) {
	hit1, hitRef1 := q.tri1.Hit(r, tMin, tMax)
	hit2, hitRef2 := q.tri2.Hit(r, tMin, tMax)
	hitRef2.U = 1 - hitRef2.U

	hitRef1.U = 1-hitRef1.U
	hitRef1.V = 1-hitRef1.V
	
	if hit1 && hit2 {
		if hitRef1.T < hitRef2.T {
			return true, hitRef1
		}
		return true, hitRef2
	}
	if hit1 {
		return true, hitRef1
	}
	if hit2 {
		return true, hitRef2
	}
	return false, HitRef{}
}

func (q *Quad) GetPos() Vec3 {
	return q.tri1.A
}

func (q *Quad) GetUV(p Vec3) (float64, float64) {
	return q.tri1.GetUV(p)
}

func (q *Quad) BoundingBox(t0, t1 float64) (bool, AABB) {
	// get bounding box of each triangle
	hit1, box1 := q.tri1.BoundingBox(t0, t1)
	hit2, box2 := q.tri2.BoundingBox(t0, t1)

	if hit1 && hit2 {
		return true, SurroundingBox(box1, box2)
	}
	if hit1 {
		return true, box1
	}
	if hit2 {
		return true, box2
	}
	return false, AABB{}
}