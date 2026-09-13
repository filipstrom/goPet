package pet

import (
	"github.com/filipstrom/goPet/object"
)

type Direction int

const (
	DirectionUp Direction = iota
	DirectionDown
	DirectionLeft
	DirectionRight
)

func (ai *AI) ResetBody() {
	ai.body.Reset()
}

func (ai *AI) Update(objects []object.Object) {
	ai.Move(objects)
	ai.body.DeAcc(0.15)

}

func (ai AI) Position() (float64, float64) {
	return ai.body.position.X, ai.body.position.Y
}

func (ai AI) GetRotation() float64 {
	return ai.body.rotation
}

func (ai AI) GetHeigt() float64 {
	return ai.body.GetHeigt()
}

func (ai AI) GetWidth() float64 {
	return ai.body.GetWidth()
}

func (ai *AI) Control(d Direction, objects []object.Object) {

	switch d {
	case DirectionUp:
		ai.body.Acc(1)
	case DirectionDown:
		ai.body.Acc(-1)
	case DirectionLeft:
		ai.body.Rotate(-1, objects)
	case DirectionRight:
		ai.body.Rotate(1, objects)

	}
}

func (ai *AI) SetPosition(x float64, y float64) {
	ai.body.position.X = x
	ai.body.position.Y = y
}

type AI struct {
	brain string
	body  object.Object
	world string
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
