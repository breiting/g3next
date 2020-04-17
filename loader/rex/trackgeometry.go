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

// TrackGeometry contains the geometry information for a REX track datablock
type TrackGeometry struct {
	geometry.Geometry
}

// NewTrackGeometry returns a new geometry information for the given REX track datablock
func NewTrackGeometry(track rexfile.Track) *TrackGeometry {

	m := new(TrackGeometry)
	m.Geometry.Init()

	positions := math32.NewArrayF32(0, 0)

	for _, c := range track.Points {
		vtx := math32.NewVector3(c.Point.X(), c.Point.Y(), c.Point.Z())
		positions.AppendVector3(vtx)
	}

	m.AddVBO(gls.NewVBO(positions).AddAttrib(gls.VertexPosition))

	return m
}
