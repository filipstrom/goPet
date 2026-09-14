package pet

import (
	"math"

	"github.com/filipstrom/goPet/object"
	"github.com/jakecoffman/cp/v2"
)

type Direction int

const (
	DirectionUp Direction = iota
	DirectionDown
	DirectionLeft
	DirectionRight
)

func (ai *AI) Control(d Direction, objects []object.Object) {
	body := ai.body.Shape.Body()

	angle := body.Angle()
	force := 400.0

	forward := cp.Vector{
		X: math.Cos(angle) * force,
		Y: math.Sin(angle) * force,
	}

	switch d {
	case DirectionUp:
		body.ApplyForceAtWorldPoint(
			forward,
			body.Position(),
		)

	case DirectionDown:
		body.ApplyForceAtWorldPoint(
			cp.Vector{
				X: -forward.X,
				Y: -forward.Y,
			},
			body.Position(),
		)
	case DirectionLeft:
		body.SetTorque(-100)
	case DirectionRight:
		body.SetTorque(100)

	}
}

type AI struct {
	brain string
	body  object.Object
	world string
}

func (ai *AI) Look(space *cp.Space) cp.SegmentQueryInfo {

	leftEye := cp.Vector{X: 18, Y: -10}
	// rightEye := cp.Vector{X: 15, Y: 8}
	leftWorld := ai.GetBody().Shape.Body().LocalToWorld(leftEye)
	// rightWorld := ai.GetBody().Shape.Body().LocalToWorld(rightEye)

	start := leftWorld
	angle := ai.body.Shape.Body().Angle()
	distance := 300.0

	end := cp.Vector{
		X: start.X + math.Cos(angle)*distance,
		Y: start.Y + math.Sin(angle)*distance,
	}

	return space.SegmentQueryFirst(start, end, 0, cp.SHAPE_FILTER_ALL)
}

func (ai *AI) GetBody() object.Object {
	return ai.body
}

func NewAI(brain string, body object.Object, world string) AI {
	return AI{
		brain: brain,
		body:  body,
		world: world,
	}
}

func (ai AI) String() string {
	return ai.brain + ", " + ai.body.String() + ", " + ai.world
}
