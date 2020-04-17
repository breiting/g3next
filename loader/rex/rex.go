// Copyright 2020. Bernhard Reitinger. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rex

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"

	"image"
	"image/draw"
	_ "image/jpeg" // support jpeg
	_ "image/png"  // support png

	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/texture"
	"github.com/roboticeyes/gorexfile/encoding/rexfile"
)

// Decoder is the REX file decoder
type Decoder struct {
	r io.Reader
}

// NewDecoder opens the reader and prepares everything for building the scene graph
// If file does not exist, an error will be returned
func NewDecoder(rexFile string) (*Decoder, error) {

	file, err := os.Open(rexFile)
	if err != nil {
		return nil, fmt.Errorf("Cannot open file %s", rexFile)
	}

	r := bufio.NewReader(file)

	return &Decoder{r: r}, nil
}

// NewDecoderReader creates a decoder with a reader
func NewDecoderReader(r io.Reader) *Decoder {
	return &Decoder{
		r: r,
	}
}

// vertexNormal is used to calculate smooth normals
type vertexNormal struct {
	normals []math32.Vector3
}

func addTexture(mat *material.Standard, data []byte) {

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		fmt.Println("Cannot decode image")
		return
	}

	// Converts image to RGBA format
	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		fmt.Println("unsupported stride")
		return
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	tex := texture.NewTexture2DFromRGBA(rgba)
	tex.SetWrapS(gls.REPEAT)
	tex.SetWrapT(gls.REPEAT)
	mat.AddTexture(tex)
}

// NewGroup creates and returns a group containing as children meshes.
// A group is returned even if there is only one object decoded.
func (dec *Decoder) NewGroup(name string) (*core.Node, error) {

	d := rexfile.NewDecoder(dec.r)
	_, rex, err := d.Decode()

	if err != nil && err.Error() != "unexpected EOF" {
		return core.NewNode(), fmt.Errorf("Cannot decode REX file: %v", err)
	}
	if rex == nil {
		return core.NewNode(), fmt.Errorf("Nothing to decode: %v", err)
	}

	return CreateNode(rex, name)
}

// CreateNode ...
func CreateNode(rex *rexfile.File, name string) (*core.Node, error) {

	meshGeometries := make(map[uint64]*MeshGeometry)
	pointGeometries := make(map[uint64]*PointGeometry)
	trackGeometries := make(map[uint64]*TrackGeometry)
	rexMaterials := make(map[uint64]material.IMaterial)

	group := core.NewNode()
	group.SetName(name)

	// Convert G3N geometries for all REX meshes
	for _, mesh := range rex.Meshes {
		meshGeometries[mesh.ID] = NewMeshGeometry(mesh)
	}

	// Get G3N materials for every stored geometry
	for _, mat := range rex.Materials {

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

		// Add texture if available (only diffuse textures are supported)
		if mat.KdTextureID != rexfile.NotSpecified {
			for _, v := range rex.Images {

				if v.ID == mat.KdTextureID {
					addTexture(phong, v.Data)
					break
				}
			}
		}

		rexMaterials[mat.ID] = phong
	}

	// Convert G3N geometries for all REX pointlists
	for _, pointList := range rex.PointLists {
		pointGeometries[pointList.ID] = NewPointGeometry(pointList)
	}

	// Convert G3N geometries for all REX tracks
	for _, track := range rex.Tracks {
		trackGeometries[track.ID] = NewTrackGeometry(track)
	}

	// Check if REX scenenodes are available, if yes take those, otherwise use old structure

	if len(rex.SceneNodes) > 0 {

		for _, node := range rex.SceneNodes {

			geom, ok := meshGeometries[node.GeometryID]
			if !ok {
				panic("Mesh not found")
			}

			var mat material.IMaterial
			mat, ok = rexMaterials[geom.MaterialID]
			if !ok {
				mat = material.NewStandard(&math32.Color{R: 0.5, G: 0.2, B: 0})
			}

			mesh := graphic.NewMesh(geom, mat)
			if geom.Name == "" {
				mesh.SetName(fmt.Sprintf("scenenode-%d", node.ID))
			} else {
				mesh.SetName(geom.Name)
			}
			mesh.SetPosition(node.Translation.X(), node.Translation.Y(), node.Translation.Z())
			mesh.SetScale(node.Scale.X(), node.Scale.Y(), node.Scale.Z())
			mesh.SetQuaternion(node.Rotation.X(), node.Rotation.Y(), node.Rotation.Z(), node.Rotation.W())
			group.Add(mesh)
		}
	} else {

		for id, geom := range meshGeometries {

			var mat material.IMaterial

			if geom.VertexColors {
				mat = material.NewBasic()
			} else {

				var ok bool
				mat, ok = rexMaterials[geom.MaterialID]
				if !ok {
					mat = material.NewStandard(&math32.Color{R: 0.5, G: 0.2, B: 0})
				}
			}

			mesh := graphic.NewMesh(geom, mat)
			if geom.Name == "" {
				mesh.SetName(fmt.Sprintf("mesh-%d", id))
			} else {
				mesh.SetName(geom.Name)
			}
			group.Add(mesh)
		}

	}

	// Points
	for _, geom := range pointGeometries {

		mat := material.NewBasic()
		points := graphic.NewPoints(geom, mat)
		group.Add(points)
	}
	// Tracks
	for _, geom := range trackGeometries {

		mat := material.NewStandard(&math32.Color{R: 0, G: 1, B: 0})
		lines := graphic.NewLineStrip(geom, mat)
		group.Add(lines)
	}

	return group, nil
}
