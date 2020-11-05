package mover

import (
	"time"

	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/math32"
)

type CameraPathMover struct {
	*core.Node

	cam     *camera.Camera // Controlled camera
	up      math32.Vector3 // The orbit axis (Y+)
	enabled bool

	location     *math32.Vector3
	path         []*math32.Vector3
	velocity     *math32.Vector3
	acceleration *math32.Vector3
	pathId       int

	start time.Time

	maxSpeed float32
	maxForce float32
	Alive    bool
	started  bool
}

func NewCameraPathMover(path []*math32.Vector3, cam *camera.Camera) *CameraPathMover {

	if len(path) < 1 {
		panic("camera path mover must not be zero length")
	}
	cpc := &CameraPathMover{
		Node:         core.NewNode(),
		cam:          cam,
		path:         path,
		location:     math32.NewVector3(cam.Position().X, cam.Position().Y, cam.Position().Z),
		velocity:     math32.NewVec3(),
		acceleration: math32.NewVec3(),
		maxSpeed:     10.0,
		maxForce:     1.0, // for direction changes
		up:           *math32.NewVector3(0, 1, 0),
		enabled:      false,
		pathId:       0,
	}

	return cpc
}

func (c *CameraPathMover) CurrentLocation() *math32.Vector3 {
	return c.location
}

func (c *CameraPathMover) applyForce(force *math32.Vector3) {

	c.acceleration.Add(force)
}

func (c *CameraPathMover) Start() {
	c.location = math32.NewVector3(c.cam.Position().X, c.cam.Position().Y, c.cam.Position().Z)
	c.enabled = true
}

func (c *CameraPathMover) Stop() {
	c.enabled = false
}

func (c *CameraPathMover) ToggleEnabled() {
	if c.enabled {
		c.Stop()
	} else {
		c.Start()
	}
}

func (c *CameraPathMover) seek(deltaTime time.Duration) {

	desired := c.path[c.pathId].Clone().Sub(c.location)
	mag := desired.Length()
	desired.Normalize()

	maxUnitPerFrame := c.maxSpeed * float32(deltaTime.Seconds())
	if mag < maxUnitPerFrame {
		c.pathId++
		if c.pathId > len(c.path)-1 {
			c.enabled = false
			return
		}
	}
	// if mag < 10 {
	// 	desired.MultiplyScalar(0.2)
	// } else {
	desired.MultiplyScalar(c.maxSpeed * float32(deltaTime.Seconds()))
	// }

	steer := desired.Sub(c.velocity)

	if steer.Length() > c.maxForce {
		steer.Normalize()
		steer.MultiplyScalar(c.maxForce)
	}
	c.applyForce(steer)
}

func (c *CameraPathMover) Update(deltaTime time.Duration) {
	if !c.enabled {
		return
	}

	c.seek(deltaTime)

	c.velocity.Add(c.acceleration)
	// fmt.Println(c.velocity)
	c.location.Add(c.velocity)
	c.rotateVelocityDirection()
	c.SetPositionVec(c.location)
	c.acceleration.MultiplyScalar(0.0)

	c.cam.SetPosition(c.location.X, c.location.Y+1, c.location.Z)
	// c.cam.LookAt(c.path[c.pathId].Clone().SetY(0), &c.up)
}

func (c *CameraPathMover) rotateVelocityDirection() {
	theta := math32.Atan2(c.velocity.Y, c.velocity.X)
	c.SetRotationZ(theta)
}
