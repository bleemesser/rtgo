package structs

type RectangularPrism struct {
	P1, P2 Vec3
	Mat    Material
}

func CreateRectangularPrism(objects []Hittable, rect RectangularPrism) []Hittable {
	sides := make([]Hittable, 6)
	sides[0] = NewXYRect(rect.P1.X, rect.P2.X, rect.P1.Y, rect.P2.Y, rect.P1.Z, rect.Mat)
	sides[1] = NewXYRect(rect.P1.X, rect.P2.X, rect.P1.Y, rect.P2.Y, rect.P2.Z, rect.Mat)
	sides[2] = NewXZRect(rect.P1.X, rect.P2.X, rect.P1.Z, rect.P2.Z, rect.P1.Y, rect.Mat)
	sides[3] = NewXZRect(rect.P1.X, rect.P2.X, rect.P1.Z, rect.P2.Z, rect.P2.Y, rect.Mat)
	sides[4] = NewYZRect(rect.P1.Y, rect.P2.Y, rect.P1.Z, rect.P2.Z, rect.P1.X, rect.Mat)
	sides[5] = NewYZRect(rect.P1.Y, rect.P2.Y, rect.P1.Z, rect.P2.Z, rect.P2.X, rect.Mat)

	objects = append(objects, sides...)

	return objects
}

func NewRectangularPrism(p1, p2 Vec3, mat Material) RectangularPrism {
	return RectangularPrism{
		P1:  p1,
		P2:  p2,
		Mat: mat,
	}
}
