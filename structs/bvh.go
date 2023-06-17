package structs

import (
	"fmt"
	"math/rand"
	"sort"
)

type BVHNode struct {
	Left, Right Hittable
	Box         AABB
}

func (b BVHNode) Hit(r *Ray, tMin float64, tMax float64) (bool, HitRef) {
	h, _ := b.Box.Hit(r, tMin, tMax)
	if h {
		hitLeft, leftRec := b.Left.Hit(r, tMin, tMax)
		hitRight, rightRec := b.Right.Hit(r, tMin, tMax)
		if hitLeft && hitRight {
			if leftRec.T < rightRec.T {
				return true, leftRec
			} else {
				return true, rightRec
			}
		} else if hitLeft {
			return true, leftRec
		} else if hitRight {
			return true, rightRec
		}
	}
	return false, HitRef{}
}

func (b BVHNode) BoundingBox(time0, time1 float64) (bool, AABB) {
	return true, b.Box
}

func (b BVHNode) GetPos() Vec3 {
	return Vec3{}
}

func NewBVHNode(srcObjects []Hittable, start, end int, time0, time1 float64) BVHNode {
	objects := make([]Hittable, len(srcObjects))
	copy(objects, srcObjects)

	axis := rand.Intn(3)
	var comparator func(a, b Hittable) bool
	if axis == 0 {
		comparator = BoxXCompare
	} else if axis == 1 {
		comparator = BoxYCompare
	} else {
		comparator = BoxZCompare
	}

	objectSpan := end - start

	var left, right Hittable
	if objectSpan == 1 {
		left, right = objects[start], objects[start]
	} else if objectSpan == 2 {
		if comparator(objects[start], objects[start+1]) {
			left, right = objects[start], objects[start+1]
		} else {
			left, right = objects[start+1], objects[start]
		}
	} else {
		sort.Slice(objects[start:end], func(i, j int) bool {
			return comparator(objects[start+i], objects[start+j])
		})

		mid := start + objectSpan/2
		left = NewBVHNode(objects, start, mid, time0, time1)
		right = NewBVHNode(objects, mid, end, time0, time1)
	}

	bl, boxLeft := left.BoundingBox(time0, time1)
	br, boxRight := right.BoundingBox(time0, time1)
	if !bl || !br {
		fmt.Println("No bounding box in BVHNode constructor.")
	}

	return BVHNode{Left: left, Right: right, Box: SurroundingBox(boxLeft, boxRight)}
}

func BoxCompare(a, b Hittable, axis int) bool {
	ba, boxA := a.BoundingBox(0, 0)
	bb, boxB := b.BoundingBox(0, 0)
	if !ba || !bb {
		fmt.Println("No bounding box in BVHNode constructor.")
	}

	return boxA.Min.GetItemByIndex(axis) < boxB.Min.GetItemByIndex(axis)
}

func BoxXCompare(a, b Hittable) bool {
	return BoxCompare(a, b, 0)
}

func BoxYCompare(a, b Hittable) bool {
	return BoxCompare(a, b, 1)
}

func BoxZCompare(a, b Hittable) bool {
	return BoxCompare(a, b, 2)
}
