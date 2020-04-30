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

// NewRexLineSetGeometry returns a new geometry information for the given REX lineset datablock
func NewRexLineSetGeometry(points rexfile.LineSet) *geometry.Geometry {

	geom := new(geometry.Geometry)

	positions := math32.NewArrayF32(len(points.Points)*3, len(points.Points)*3)

	j := 0
	for _, c := range points.Points {
		for i := 0; i < 3; i++ {
			positions[j] = c[i]
			j++
		}
	}

	geom.AddVBO(gls.NewVBO(positions).AddAttrib(gls.VertexPosition))
	return geom
}
