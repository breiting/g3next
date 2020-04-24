// Copyright 2020. Bernhard Reitinger. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rex

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/breiting/g3next/entity"
	"github.com/breiting/g3next/geom"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
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

	return CreateRexNode(rex, name)
}

func CreateRexNode(rex *rexfile.File, name string) (*core.Node, error) {

	var meshes []*entity.RexMesh

	group := core.NewNode()
	group.SetName(name)

	for _, mesh := range rex.Meshes {
		meshes = append(meshes, entity.NewRexMesh(mesh))
	}

	for _, mat := range rex.Materials {
		for _, v := range meshes {
			v.OfferMaterial(mat)
		}
	}

	for _, img := range rex.Images {
		for _, v := range meshes {
			v.OfferTexture(img)
		}
	}

	for _, m := range meshes {
		group.Add(m)
	}

	for _, pointList := range rex.PointLists {
		mat := material.NewBasic()
		points := graphic.NewPoints(geom.NewRexPointGeometry(pointList), mat)
		group.Add(points)
	}

	for _, track := range rex.Tracks {
		mat := material.NewStandard(&math32.Color{R: 0, G: 1, B: 0})
		lines := graphic.NewLineStrip(geom.NewRexTrackGeometry(track), mat)
		group.Add(lines)
	}

	// TODO currently unsupported
	// if len(rex.SceneNodes) > 0 {
	// 	for _, node := range rex.SceneNodes {
	//
	// 		geom, ok := meshGeometries[node.GeometryID]
	// 		if !ok {
	// 			panic("Mesh not found")
	// 		}
	//
	// 		var mat material.IMaterial
	// 		mat, ok = rexMaterials[rex.Meshes[node.GeometryID].MaterialID]
	// 		if !ok {
	// 			mat = material.NewStandard(&math32.Color{R: 0.5, G: 0.2, B: 0})
	// 		}
	//
	// 		mesh := graphic.NewMesh(geom, mat)
	// 		if rex.Meshes[node.GeometryID].Name == "" {
	// 			mesh.SetName(fmt.Sprintf("scenenode-%d", node.ID))
	// 		} else {
	// 			mesh.SetName(rex.Meshes[node.GeometryID].Name)
	// 		}
	// 		mesh.SetPosition(node.Translation.X(), node.Translation.Y(), node.Translation.Z())
	// 		mesh.SetScale(node.Scale.X(), node.Scale.Y(), node.Scale.Z())
	// 		mesh.SetQuaternion(node.Rotation.X(), node.Rotation.Y(), node.Rotation.Z(), node.Rotation.W())
	// 		group.Add(mesh)
	// 	}
	// }

	return group, nil
}
