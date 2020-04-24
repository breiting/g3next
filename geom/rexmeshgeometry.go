// Copyright 2020. Bernhard Reitinger. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package geom

import (
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/math32"
	"github.com/roboticeyes/gorexfile/encoding/rexfile"
)

// NewRexMeshGeometry returns a new geometry information for the given REX mesh datablock
func NewRexMeshGeometry(mesh rexfile.Mesh) *geometry.Geometry {

	geom := new(geometry.Geometry)

	positions := math32.NewArrayF32(len(mesh.Coords)*3, len(mesh.Coords)*3)
	colors := math32.NewArrayF32(len(mesh.Colors)*3, len(mesh.Colors)*3)
	uvs := math32.NewArrayF32(len(mesh.TexCoords)*2, len(mesh.TexCoords)*2)
	normals := math32.NewArrayF32(len(mesh.Coords)*3, len(mesh.Coords)*3)
	indices := math32.NewArrayU32(len(mesh.Triangles)*3, len(mesh.Triangles)*3)
	tempNormals := make([][]math32.Vector3, len(mesh.Coords))

	j := 0
	for _, c := range mesh.Coords {
		for i := 0; i < 3; i++ {
			positions[j] = c[i]
			j++
		}
	}

	j = 0
	for _, c := range mesh.Colors {
		for i := 0; i < 3; i++ {
			colors[j] = c[i]
			j++
		}
	}

	j = 0
	for _, uv := range mesh.TexCoords {
		for i := 0; i < 2; i++ {
			uvs[j] = uv[i]
			j++
		}
	}

	j = 0
	for _, t := range mesh.Triangles {
		indices[j] = t.V0
		indices[j+1] = t.V1
		indices[j+2] = t.V2
		j += 3

		// calculate normals per face
		var v0, v1, v2 math32.Vector3
		var n0, n1, n2 math32.Vector3
		var sub1, sub2 math32.Vector3
		positions.GetVector3(int(t.V0*3), &v0)
		positions.GetVector3(int(t.V1*3), &v1)
		positions.GetVector3(int(t.V2*3), &v2)
		n0.CrossVectors(sub1.SubVectors(&v1, &v0), sub2.SubVectors(&v2, &v0)).Normalize()
		n1.CrossVectors(sub1.SubVectors(&v2, &v1), sub2.SubVectors(&v0, &v1)).Normalize()
		n2.CrossVectors(sub1.SubVectors(&v0, &v2), sub2.SubVectors(&v1, &v2)).Normalize()
		tempNormals[t.V0] = append(tempNormals[t.V0], n0)
		tempNormals[t.V1] = append(tempNormals[t.V1], n1)
		tempNormals[t.V2] = append(tempNormals[t.V2], n2)
	}

	// calculate smooth normals
	for i, n := range tempNormals {
		var sum math32.Vector3
		for _, normal := range n {
			sum.Add(&normal)
		}
		sum.DivideScalar(float32(len(n)))
		normals.SetVector3(i*3, sum.Normalize())
	}

	geom.SetIndices(indices)
	geom.AddVBO(gls.NewVBO(positions).AddAttrib(gls.VertexPosition))
	if len(colors) > 0 {
		geom.AddVBO(gls.NewVBO(colors).AddAttrib(gls.VertexColor))
	}
	geom.AddVBO(gls.NewVBO(normals).AddAttrib(gls.VertexNormal))
	geom.AddVBO(gls.NewVBO(uvs).AddAttrib(gls.VertexTexcoord))

	return geom
}
