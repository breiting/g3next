package entity

import (
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
)

type TerrainMesh struct {
	*graphic.Mesh

	geometry *geometry.Geometry
}

func NewTerrainMesh(width, height float32, nrWidthSegments, nrHeightSegments int) *TerrainMesh {

	t := &TerrainMesh{}

	t.geometry = geometry.NewSegmentedPlane(width, height, nrWidthSegments, nrHeightSegments)
	t.Mesh = graphic.NewMesh(t.geometry, material.NewStandard(&math32.Color{R: 1, G: 1, B: 1}))
	return t
}

func (t *TerrainMesh) ApplyNoiseMap() {
}
