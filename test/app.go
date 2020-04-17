package main

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/breiting/g3next/entity"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/util"
	"github.com/g3n/engine/util/helper"
	"github.com/g3n/engine/window"
)

const (
	gridSize  = 10
	targetFPS = 30
)

// App is the actual app logic
type App struct {
	window.IWindow
	renderer *renderer.Renderer

	root      *core.Node
	scene     *core.Node
	headLight *light.Directional

	showWireframe bool

	camera     *camera.Camera
	orbit      *camera.OrbitControl
	frameRater *util.FrameRater
}

// NewApp creates a new app
func NewApp(width, height int, title string) *App {

	a := &App{}

	// setup window
	err := window.Init(width, height, title)
	if err != nil {
		panic(err)
	}
	a.IWindow = window.Get()

	// setup renderer
	a.renderer = renderer.NewRenderer(a.Gls())
	err = a.renderer.AddDefaultShaders()
	if err != nil {
		panic(err)
	}
	glVersion := a.Gls().GetString(gls.VERSION)
	fmt.Println("OpenGL version:", glVersion)
	a.Gls().SetCheckErrors(true)
	a.frameRater = util.NewFrameRater(targetFPS)

	a.setupScene()
	a.setupLights()

	// Subscribe to events
	a.Subscribe(window.OnWindowSize, func(evname string, ev interface{}) { a.OnWindowResize() })
	a.OnWindowResize()

	a.Subscribe(window.OnKeyDown, a.onKey)
	a.Subscribe(window.OnKeyUp, a.onKey)
	// a.Subscribe(window.OnMouseDown, a.onMouseDown)
	// a.Subscribe(window.OnCursor, a.onMouseMove)
	return a
}

func (a *App) onKey(evname string, ev interface{}) {

	var state bool
	if evname == window.OnKeyDown {
		state = true
	} else {
		state = false
	}

	kev := ev.(*window.KeyEvent)
	switch kev.Key {
	case window.KeyQ:
		if state == true {
			a.Window().SetShouldClose(true)
		}
	case window.KeyF1:
		if state == true {
			a.toggleWireframe()
		}
	}
}

func (a *App) toggleWireframe() {
	a.showWireframe = !a.showWireframe
	a.toggleWireframeNode(a.scene)
}

func (a *App) toggleWireframeNode(inode core.INode) {

	node := inode.GetNode()
	for _, child := range node.Children() {
		a.toggleWireframeNode(child)
	}

	ig, ok := inode.(graphic.IGraphic)
	if !ok {
		return
	}
	gr := ig.GetGraphic()
	imat := gr.GetMaterial(0).GetMaterial()
	imat.SetWireframe(a.showWireframe)
}

// Run runs the application render loop
func (a *App) Run() {

	frameStart := time.Now()

	var frameDelta time.Duration

	for true {
		if a.Window().ShouldClose() {
			break
		}
		now := time.Now()
		frameDelta = now.Sub(frameStart)
		frameStart = now

		a.render(frameDelta)

		a.Window().SwapBuffers()
		a.Window().PollEvents()
	}

	a.Destroy()
}

func (a *App) render(deltaTime time.Duration) {

	a.frameRater.Start()

	a.Gls().Clear(gls.COLOR_BUFFER_BIT | gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT)

	camPos := a.camera.Position()
	a.headLight.SetPosition(camPos.X, camPos.Y, camPos.Z)

	err := a.renderer.Render(a.root, a.camera)
	if err != nil {
		panic(err)
	}

	a.frameRater.Wait()
}

// OnWindowResize is the default handler for resizing the window
func (a *App) OnWindowResize() {

	width, height := a.GetFramebufferSize()
	a.Gls().Viewport(0, 0, int32(width), int32(height))
	a.camera.SetAspect(float32(width) / float32(height))
}

// Window returns the original window
func (a *App) Window() *window.GlfwWindow {
	return a.IWindow.(*window.GlfwWindow)
}

func (a *App) setupLights() {
	a.headLight = light.NewDirectional(&math32.Color{R: 1, G: 1, B: 1}, 0.1)
	a.headLight.SetPosition(0, 0, 10)
	a.scene.Add(a.headLight)

	a.scene.Add(light.NewAmbient(&math32.Color{R: 1, G: 1, B: 1}, 0.5))
}

func createImage() *graphic.Mesh {
	f, err := os.Open("data/gopher.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	img := entity.NewImage(f)
	img.SetRotationX(math.Pi / 4.0)
	img.SetPositionZ(1.0)
	// img.SetRotationZ(math.Pi / 2.0)
	return img
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
