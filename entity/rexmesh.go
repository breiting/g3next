package entity

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg" // support jpeg
	_ "image/png"  // support png
	"math/rand"

	"github.com/breiting/g3next/geom"
	"github.com/breiting/g3next/mat"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/texture"
	"github.com/roboticeyes/gorexfile/encoding/rexfile"
)

var (
	defaultMeshMaterial = &math32.Color{R: 1, G: 0, B: 1}
)

type RexMesh struct {
	*graphic.Mesh

	data     rexfile.Mesh
	material rexfile.Material
	phong    *material.Standard
}

func NewRexMesh(data rexfile.Mesh) *RexMesh {

	mesh := &RexMesh{
		data: data,
	}
	geom := geom.NewRexMeshGeometry(data)

	// select material based on vertex coloring
	if geom.VBO(gls.VertexColor) != nil {
		mesh.Mesh = graphic.NewMesh(geom, material.NewBasic())
	} else {
		mesh.Mesh = graphic.NewMesh(geom, material.NewStandard(
			&math32.Color{
				R: rand.Float32(),
				G: rand.Float32(),
				B: rand.Float32(),
			},
		))
	}
	if data.Name != "" {
		mesh.Mesh.SetName(data.Name)
	} else {
		mesh.Mesh.SetName(fmt.Sprintf("rexmesh-%d", data.ID))
	}
	return mesh
}

// OfferMaterial is used by the caller to propagate a material. If the
// materialID is matching, the material is replaced
func (m *RexMesh) OfferMaterial(material rexfile.Material) {

	if material.ID != m.data.MaterialID {
		return
	}

	// Ignore StandardMaterial if VertexColor is available
	if m.Graphic.GetGeometry().VBO(gls.VertexColor) != nil {
		fmt.Printf("Ignoring Material since vertex colors are available (mesh: %d)\n", m.data.ID)
		return
		// } else {
		// 	fmt.Printf("Set material %d for mesh %d\n", material.ID, m.data.ID)
	}

	m.phong = mat.NewRexStandardMaterial(material)
	m.Graphic.ClearMaterials()
	m.Graphic.AddMaterial(m, m.phong, 0, 0)
	m.material = material
}

// OfferTexture is used by the caller to propagate a texture for the material. If the
// ID is matching, the texture is applied
func (m *RexMesh) OfferTexture(img rexfile.Image) {

	if img.ID != m.material.KdTextureID {
		return
	}

	if m.phong == nil {
		return
	}

	decodedImg, _, err := image.Decode(bytes.NewReader(img.Data))
	if err != nil {
		fmt.Println("Cannot decode texture", img.ID)
		return
	}

	// Converts image to RGBA format
	rgba := image.NewRGBA(decodedImg.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		fmt.Println("Unsupported stride information for texture", img.ID)
		return
	}
	draw.Draw(rgba, rgba.Bounds(), decodedImg, image.Point{0, 0}, draw.Src)

	tex := texture.NewTexture2DFromRGBA(rgba)
	tex.SetWrapS(gls.REPEAT)
	tex.SetWrapT(gls.REPEAT)
	m.phong.AddTexture(tex)
}
