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

// NewRexTrackGeometry returns a new geometry information for the given REX track datablock
func NewRexTrackGeometry(track rexfile.Track) *geometry.Geometry {

	geom := new(geometry.Geometry)

	pos := math32.NewArrayF32(len(track.Points)*3, len(track.Points)*3)

	i := 0
	for _, c := range track.Points {
		pos[i] = c.Point.X()
		pos[i+1] = c.Point.Y()
		pos[i+2] = c.Point.Z()
		i += 3
	}

	geom.AddVBO(gls.NewVBO(pos).AddAttrib(gls.VertexPosition))
	return geom
}
