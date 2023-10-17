package utility

import (
	st "rtgo/structs"

	"github.com/g3n/engine/loader/obj"
)

func LoadOBJFile(filename string, mat st.Material, offset st.Vec3) []st.Hittable {
	// load the obj from file, parse it into triangles

	// load the obj from file
	d, err := obj.Decode(filename, "")
	if err != nil {
		panic(err)
	}
	// parse the obj into triangles
	var triangles []st.Hittable
	for _, obj := range d.Objects {
		for _, face := range obj.Faces {
			// check the number of vertices in the face
			numVertices := len(face.Vertices)
			if numVertices < 3 {
				continue
			}

			// split square faces into two triangles
			if numVertices == 4 {
				v1 := d.Vertices[face.Vertices[0]*3 : face.Vertices[0]*3+3]
				v2 := d.Vertices[face.Vertices[1]*3 : face.Vertices[1]*3+3]
				v3 := d.Vertices[face.Vertices[2]*3 : face.Vertices[2]*3+3]
				v4 := d.Vertices[face.Vertices[3]*3 : face.Vertices[3]*3+3]

				tri1 := st.NewTriangle(st.Vec3{X: float64(v1[0]), Y: float64(v1[1]), Z: float64(v1[2])}.Sub(offset), st.Vec3{X: float64(v2[0]), Y: float64(v2[1]), Z: float64(v2[2])}.Sub(offset), st.Vec3{X: float64(v3[0]), Y: float64(v3[1]), Z: float64(v3[2])}.Sub(offset), mat)
				// tri2 := st.NewTriangle(st.Vec3{X: float64(v1[0]), Y: float64(v1[1]), Z: float64(v1[2])}, st.Vec3{X: float64(v3[0]), Y: float64(v3[1]), Z: float64(v3[2])}, st.Vec3{X: float64(v4[0]), Y: float64(v4[1]), Z: float64(v4[2])}, mat)
				tri2 := st.NewTriangle(st.Vec3{X: float64(v1[0]), Y: float64(v1[1]), Z: float64(v1[2])}.Sub(offset), st.Vec3{X: float64(v3[0]), Y: float64(v3[1]), Z: float64(v3[2])}.Sub(offset), st.Vec3{X: float64(v4[0]), Y: float64(v4[1]), Z: float64(v4[2])}.Sub(offset), mat)

				triangles = append(triangles, tri1, tri2)
			} else {
				// create a triangle for all other faces
				v1 := d.Vertices[face.Vertices[0]*3 : face.Vertices[0]*3+3]
				v2 := d.Vertices[face.Vertices[1]*3 : face.Vertices[1]*3+3]
				v3 := d.Vertices[face.Vertices[2]*3 : face.Vertices[2]*3+3]

				// tri := st.NewTriangle(st.Vec3{X: float64(v1[0]), Y: float64(v1[1]), Z: float64(v1[2])}, st.Vec3{X: float64(v2[0]), Y: float64(v2[1]), Z: float64(v2[2])}, st.Vec3{X: float64(v3[0]), Y: float64(v3[1]), Z: float64(v3[2])}, mat)
				tri := st.NewTriangle(st.Vec3{X: float64(v1[0]), Y: float64(v1[1]), Z: float64(v1[2])}.Sub(offset), st.Vec3{X: float64(v2[0]), Y: float64(v2[1]), Z: float64(v2[2])}.Sub(offset), st.Vec3{X: float64(v3[0]), Y: float64(v3[1]), Z: float64(v3[2])}.Sub(offset), mat)

				triangles = append(triangles, tri)
			}
		}
	}

	return triangles
}
