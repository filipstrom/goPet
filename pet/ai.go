package pet

import (
	"fmt"
	"image/color"
	"math"

	"github.com/filipstrom/goPet/object"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
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

func (ai *AI) Eat(space *cp.Space) {
	space.ShapeQuery(ai.EatSensor, func(shape *cp.Shape, points *cp.ContactPointSet) {
		if shape.UserData == "food" {
			fmt.Print("Mums")
			space.RemoveShape(shape)
			space.RemoveBody(shape.Body())
		}
		fmt.Println("Hej")

	})

}

type AI struct {
	brain     string
	body      object.Object
	EatSensor *cp.Shape
	world     string
}

func (ai *AI) rayCast(start cp.Vector, end cp.Vector, screen *ebiten.Image) {
	vector.StrokeLine(screen, float32(start.X), float32(start.Y), float32(end.X), float32(end.Y), 2, color.Black, false)
}

func (ai *AI) Look(space *cp.Space, screen *ebiten.Image) cp.SegmentQueryInfo {

	leftEye := cp.Vector{X: 18, Y: -10}
	rightEye := cp.Vector{X: 18, Y: 10}
	leftWorld := ai.GetBody().Shape.Body().LocalToWorld(leftEye)
	rightWorld := ai.GetBody().Shape.Body().LocalToWorld(rightEye)

	startl := leftWorld
	startr := rightWorld

	angle := ai.body.Shape.Body().Angle()
	distance := 300.0

	end := cp.Vector{
		X: startl.X + math.Cos(angle)*distance,
		Y: startl.Y + math.Sin(angle)*distance,
	}

	endr := cp.Vector{
		X: startr.X + math.Cos(angle)*distance,
		Y: startr.Y + math.Sin(angle)*distance,
	}

	ai.rayCast(startr, endr, screen)
	ai.rayCast(startl, end, screen)
	return space.SegmentQueryFirst(startl, end, 0, cp.SHAPE_FILTER_ALL)
}

func (ai *AI) GetBody() object.Object {
	return ai.body
}

func NewAI(brain string, body object.Object, world string, eatSensor *cp.Shape) AI {
	return AI{
		brain:     brain,
		body:      body,
		world:     world,
		EatSensor: eatSensor,
	}
}

func (ai AI) String() string {
	return ai.brain + ", " + ai.body.String() + ", " + ai.world
}
