package main

import (
	"bytes"
	"io/ioutil"
	"math"
	"net/http"
	"os"

	"github.com/breiting/g3next/entity"
	"github.com/breiting/g3next/loader/rex"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/util/helper"
	"github.com/g3n/engine/window"
)

const (
	sampleRexFileUrl = "https://github.com/roboticeyes/gorexfile/raw/master/examples/models/capsule/capsule.rex"
)

func createImage() *core.Node {
	f, err := os.Open("data/gopher.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	img := entity.NewImage(f)

	container := core.NewNode()
	container.SetPosition(0, 0, 1)
	container.SetRotationQuat(math32.NewQuaternion(0.707, 0, 0, 0.707).Inverse())
	container.Add(img)
	return container
}

func createRexFile() *core.Node {
	resp, err := http.Get(sampleRexFileUrl)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	reader := bytes.NewReader(body)

	decoder := rex.NewDecoderReader(reader)
	rexNode, err := decoder.NewGroup("rex")
	if err != nil {
		panic(err)
	}
	rexNode.SetScale(0.5, 0.5, 0.5)
	rexNode.SetRotationX(math.Pi / 2.0)
	return rexNode
}

func (a *App) setupScene() {

	a.root = core.NewNode()
	a.scene = core.NewNode()

	a.scene.SetRotationX(-math.Pi / 2.0)
	a.root.Add(helper.NewGrid(gridSize, 1, &math32.Color{R: 0.4, G: 0.4, B: 0.4}))
	a.scene.Add(helper.NewAxes(1))
	a.root.Add(a.scene)

	// begin add entities
	a.scene.Add(createImage())
	a.scene.Add(createRexFile())

	// end add entities

	a.UnsubscribeAllID(a)
	a.DisposeAllCustomCursors()
	a.SetCursor(window.ArrowCursor)

	// Set default background color
	a.Gls().ClearColor(0.6, 0.6, 0.6, 1.0)

	// Reset renderer z-sorting flag
	a.renderer.SetObjectSorting(true)

	// Camera
	w, h := a.GetSize()
	a.camera = camera.New(float32(w) / float32(h))
	a.root.Add(a.camera)
	a.orbit = camera.NewOrbitControl(a.camera)
	a.camera.SetPosition(10, 10, 10)
	a.camera.UpdateSize(5)
	a.camera.LookAt(math32.NewVector3(0, 0, 0), math32.NewVector3(0, 1, 0))
	a.camera.SetProjection(camera.Perspective)
	a.orbit.Reset()
}
