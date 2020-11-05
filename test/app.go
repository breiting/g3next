package main

import (
	"fmt"
	"time"

	"github.com/breiting/g3next/mover"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/experimental/collision"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/util"
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

	camera      *camera.Camera
	orbit       *camera.OrbitControl
	cameramover *mover.CameraPathMover
	frameRater  *util.FrameRater
	rayCaster   *collision.Raycaster
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
	a.Subscribe(window.OnCursor, a.onMouseMove)

	a.rayCaster = collision.NewRaycaster(&math32.Vector3{}, &math32.Vector3{})

	return a
}

func (a *App) onMouseMove(evname string, ev interface{}) {

	mouseEvent := ev.(*window.CursorEvent)
	width, height := a.GetFramebufferSize()

	x := 2*(mouseEvent.Xpos/float32(width)) - 1
	y := -2*(mouseEvent.Ypos/float32(height)) + 1
	a.rayCaster.SetFromCamera(a.camera, x, y)

	intersects := a.rayCaster.IntersectObjects(a.scene.Children(), true)
	if len(intersects) == 0 {
		fmt.Println("nothing found")
		return
	}
	fmt.Printf("World position: %v\n", intersects[0].Point)
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
	case window.KeySpace:
		if state == true {
			a.cameramover.ToggleEnabled()
			a.orbit.SetTarget(*a.cameramover.CurrentLocation())
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

	a.cameramover.Update(deltaTime)

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
