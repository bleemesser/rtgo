package utility

import (
	"github.com/g3n/engine/loader/obj"
	st "rtgo/structs"
)

// type Triangle struct {
// 	A, B, C Vec3
// 	Mat     Material
// }

// func NewTriangle(a, b, c Vec3, mat Material) *Triangle {
// 	return &Triangle{A: a, B: b, C: c, Mat: mat}
// }


// type Decoder struct {
// 	Objects   []Object             // decoded objects
// 	Matlib    string               // name of the material lib
// 	Materials map[string]*Material // maps material name to object
// 	Vertices  math32.ArrayF32      // vertices positions array
// 	Normals   math32.ArrayF32      // vertices normals
// 	Uvs       math32.ArrayF32      // vertices texture coordinates
// 	Warnings  []string             // warning messages
// 	// contains filtered or unexported fields
// }

// type Face struct {
// 	Vertices []int  // Indices to the face vertices
// 	Uvs      []int  // Indices to the face UV coordinates
// 	Normals  []int  // Indices to the face normals
// 	Material string // Material name
// 	Smooth   bool   // Smooth face
// }

// type Material struct {
// 	Name       string       // Material name
// 	Illum      int          // Illumination model
// 	Opacity    float32      // Opacity factor
// 	Refraction float32      // Refraction factor
// 	Shininess  float32      // Shininess (specular exponent)
// 	Ambient    math32.Color // Ambient color reflectivity
// 	Diffuse    math32.Color // Diffuse color reflectivity
// 	Specular   math32.Color // Specular color reflectivity
// 	Emissive   math32.Color // Emissive color
// 	MapKd      string       // Texture file linked to diffuse color
// }

// type Object struct {
// 	Name  string // Object name
// 	Faces []Face // Faces
// 	// contains filtered or unexported fields
// }


func LoadOBJFile(filename string, mat st.Material) []*st.Triangle {
	// load the obj from file, parse it into triangles

	// load the obj from file
	d, err := obj.Decode(filename, "")
	if err != nil {
		panic(err)
	}

	// parse the obj into triangles
	var triangles []*st.Triangle
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

				tri1 := st.NewTriangle(st.Vec3{float64(v1[0]), float64(v1[1]), float64(v1[2])}, st.Vec3{float64(v2[0]), float64(v2[1]), float64(v2[2])}, st.Vec3{float64(v3[0]), float64(v3[1]), float64(v3[2])}, mat)
				tri2 := st.NewTriangle(st.Vec3{float64(v1[0]), float64(v1[1]), float64(v1[2])}, st.Vec3{float64(v3[0]), float64(v3[1]), float64(v3[2])}, st.Vec3{float64(v4[0]), float64(v4[1]), float64(v4[2])}, mat)

				triangles = append(triangles, tri1, tri2)
			} else {
				// create a triangle for all other faces
				v1 := d.Vertices[face.Vertices[0]*3 : face.Vertices[0]*3+3]
				v2 := d.Vertices[face.Vertices[1]*3 : face.Vertices[1]*3+3]
				v3 := d.Vertices[face.Vertices[2]*3 : face.Vertices[2]*3+3]

				tri := st.NewTriangle(st.Vec3{float64(v1[0]), float64(v1[1]), float64(v1[2])}, st.Vec3{float64(v2[0]), float64(v2[1]), float64(v2[2])}, st.Vec3{float64(v3[0]), float64(v3[1]), float64(v3[2])}, mat)

				triangles = append(triangles, tri)
			}
		}
	}

	return triangles
}
