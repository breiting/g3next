package mat

import (
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/roboticeyes/gorexfile/encoding/rexfile"
)

// NewRexStandardMaterial generates a new material
func NewRexStandardMaterial(mat rexfile.Material) *material.Standard {

	phong := material.NewStandard(&math32.Color{
		R: mat.KdRgb.X(),
		G: mat.KdRgb.Y(),
		B: mat.KdRgb.Z()})

	phong.SetSpecularColor(&math32.Color{
		R: mat.KsRgb.X(),
		G: mat.KsRgb.Y(),
		B: mat.KsRgb.Z()})

	if mat.Ns != 0 {
		phong.SetShininess(mat.Ns)
	}

	phong.SetSide(material.SideFront)
	phong.SetOpacity(mat.Alpha)

	return phong
}
