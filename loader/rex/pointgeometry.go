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

// PointGeometry contains the geometry information for a REX point datablock
type PointGeometry struct {
	geometry.Geometry
}

// NewPointGeometry returns a new geometry information for the given REX point datablock
func NewPointGeometry(points rexfile.PointList) *PointGeometry {

	m := new(PointGeometry)
	m.Geometry.Init()

	positions := math32.NewArrayF32(0, 0)
	colors := math32.NewArrayF32(0, 0)

	for _, c := range points.Points {
		vtx := math32.NewVector3(c.X(), c.Y(), c.Z())
		positions.AppendVector3(vtx)
	}

	for _, c := range points.Colors {
		vtx := math32.NewVector3(c.X(), c.Y(), c.Z())
		colors.AppendVector3(vtx)
	}

	m.AddVBO(gls.NewVBO(positions).AddAttrib(gls.VertexPosition))
	if len(colors) > 0 {
		m.AddVBO(gls.NewVBO(colors).AddAttrib(gls.VertexColor))
	}

	return m
}
