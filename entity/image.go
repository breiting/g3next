// Copyright 2020. Bernhard Reitinger. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package entity

import (
	"bytes"
	"image"
	"image/draw"
	_ "image/jpeg" // support jpeg
	_ "image/png"  // support png
	"io"

	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/texture"
	"github.com/roboticeyes/gorexfile/encoding/rexfile"
)

const (
	imageWidthInMeters = 0.4
)

// NewImage returns a new mesh for the given REX image datablock
func NewImageFromRex(block rexfile.Image) *graphic.Mesh {
	return NewImage(bytes.NewReader(block.Data))
}

// NewImage returns a new mesh containing the specified image which can be read by the io.Reader
// interface (currently supported are PNG and JPG)
func NewImage(r io.Reader) *graphic.Mesh {

	img, _, err := image.Decode(r)
	if err != nil {
		panic(err)
	}

	// Converts image to RGBA format
	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		panic("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	tex := texture.NewTexture2DFromRGBA(rgba)
	xs := float32(img.Bounds().Size().X)
	ys := float32(img.Bounds().Size().Y)
	geom := geometry.NewPlane(imageWidthInMeters, imageWidthInMeters*(ys/xs))
	mat := material.NewStandard(&math32.Color{R: 1, G: 1, B: 1})
	mat.SetSide(material.SideDouble)
	mat.AddTexture(tex)
	mesh := graphic.NewMesh(geom, mat)
	return mesh
}
