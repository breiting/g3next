// Copyright 2020. Bernhard Reitinger. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rex

import (
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/math32"
	"github.com/roboticeyes/gorexfile/encoding/rexfile"
)

// MeshGeometry contains the geometry information for a REX mesh datablock
type MeshGeometry struct {
	geometry.Geometry

	Name         string
	MaterialID   uint64
	VertexColors bool // indicates if vertex colors are set, if so, we have to use basic material
}

// NewMeshGeometry returns a new geometry information for the given REX mesh datablock
func NewMeshGeometry(mesh rexfile.Mesh) *MeshGeometry {

	m := new(MeshGeometry)
	m.Geometry.Init()

	m.Name = mesh.Name
	m.MaterialID = mesh.MaterialID

	// Create buffers
	positions := math32.NewArrayF32(0, 0)
	colors := math32.NewArrayF32(0, 0)
	uvs := math32.NewArrayF32(0, 0)
	normals := math32.NewArrayF32(len(mesh.Coords)*3, len(mesh.Coords)*3)
	indices := math32.NewArrayU32(0, 0)
	tempNormals := make([]vertexNormal, len(mesh.Coords))

	for _, c := range mesh.Coords {
		vtx := math32.NewVector3(c.X(), c.Y(), c.Z())
		positions.AppendVector3(vtx)
	}

	for _, c := range mesh.Colors {
		vtx := math32.NewVector3(c.X(), c.Y(), c.Z())
		colors.AppendVector3(vtx)
	}

	for _, uv := range mesh.TexCoords {
		coords := math32.NewVector2(uv.X(), uv.Y())
		uvs.AppendVector2(coords)
	}

	for _, t := range mesh.Triangles {
		indices.Append(t.V0)
		indices.Append(t.V1)
		indices.Append(t.V2)

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
		tempNormals[t.V0].normals = append(tempNormals[t.V0].normals, n0)
		tempNormals[t.V1].normals = append(tempNormals[t.V1].normals, n1)
		tempNormals[t.V2].normals = append(tempNormals[t.V2].normals, n2)
	}

	// calculate smooth normals
	for i, n := range tempNormals {
		var sum math32.Vector3
		for _, normal := range n.normals {
			sum.Add(&normal)
		}
		sum.DivideScalar(float32(len(n.normals)))
		normals.SetVector3(i*3, sum.Normalize())
	}

	m.SetIndices(indices)
	m.AddVBO(gls.NewVBO(positions).AddAttrib(gls.VertexPosition))
	if len(colors) > 0 {
		m.AddVBO(gls.NewVBO(colors).AddAttrib(gls.VertexColor))
		m.VertexColors = true
	}
	m.AddVBO(gls.NewVBO(normals).AddAttrib(gls.VertexNormal))
	m.AddVBO(gls.NewVBO(uvs).AddAttrib(gls.VertexTexcoord))

	return m
}
