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

func (ai *AI) Control(d Direction) {
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
func (ai *AI) changeState(key string, value float64) {
	switch key {
	case "hunger":
		ai.mood.hunger = clamp(ai.mood.hunger + value)
	}
}

func (ai *AI) rotate(s float32) {
	ai.body.Shape.Body().SetTorque(float64(100 * s))
}

func (ai *AI) addForce(f float32) {

	body := ai.body.Shape.Body()

	angle := body.Angle()
	force := float64(400.0 * f)

	forward := cp.Vector{
		X: math.Cos(angle) * force,
		Y: math.Sin(angle) * force,
	}

	body.ApplyForceAtWorldPoint(
		cp.Vector{
			X: forward.X,
			Y: forward.Y,
		},
		body.Position(),
	)

}

func clamp(v float64) float64 {
	if v > 1 {
		return 1
	}
	if v < -1 {
		return -1
	}
	return v
}

func (ai *AI) Eat() bool {
	eated := false
	ai.space.ShapeQuery(ai.EatSensor, func(shape *cp.Shape, points *cp.ContactPointSet) {
		data, ok := shape.UserData.(object.ShapeData)

		if !ok {
			return
		}

		if data.IsFood {
			fmt.Print("Mums")
			ai.space.RemoveShape(shape)
			ai.space.RemoveBody(shape.Body())
			ai.changeState("hunger", 0.2)
			eated = true
		}
		// fmt.Println("Hej")

	})
	return eated

}

func (ai *AI) GetMood() Mood {
	return ai.mood

}
func (ai *AI) Move(out []float32) bool {
	ai.addForce(out[0])
	ai.rotate(out[1])
	if out[2] > 0 {
		return ai.Eat()
	}
	return false
}

type Mood struct {
	hunger float64
}

func (mood Mood) String() string {
	return fmt.Sprintf("Hunger: %.3f", mood.hunger)
}

func (ai *AI) DistanceToFood() float64 {
	aiPos := ai.body.Shape.Body().Position()

	distance := math.Inf(1)

	ai.space.EachShape(func(shape *cp.Shape) {
		data, ok := shape.UserData.(object.ShapeData)
		if !ok || !data.IsFood {
			return
		}

		foodPos := shape.Body().Position()

		dx := foodPos.X - aiPos.X
		dy := foodPos.Y - aiPos.Y

		distance = math.Sqrt(dx*dx + dy*dy)
	})

	return distance
}

type AI struct {
	brain     string
	body      object.Object
	EatSensor *cp.Shape
	space     *cp.Space
	screen    *ebiten.Image
	mood      Mood
}

func (ai *AI) rayCast(start cp.Vector, end cp.Vector, screen *ebiten.Image) {
	vector.StrokeLine(screen, float32(start.X), float32(start.Y), float32(end.X), float32(end.Y), 2, color.Black, false)
}

func (ai *AI) Look() ([]float32, []float32) {

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

	//ai.rayCast(startr, endr, ai.screen)
	//ai.rayCast(startl, end, ai.screen)

	leftHit := ai.space.SegmentQueryFirst(
		startl, end, 0, cp.SHAPE_FILTER_ALL,
	)

	rightHit := ai.space.SegmentQueryFirst(
		startr, endr, 0, cp.SHAPE_FILTER_ALL,
	)

	return ai.getApperance(leftHit), ai.getApperance(rightHit)
}

func (ai *AI) getApperance(hit cp.SegmentQueryInfo) []float32 {
	if hit.Shape == nil {
		return []float32{300, 0, 0, 0}
	}

	distance := float32(hit.Alpha * 300)

	data, ok := hit.Shape.UserData.(object.ShapeData)

	if !ok {
		switch hit.Shape {
		case ai.body.Shape:
			fmt.Println("LOOK HIT OWN BODY")

		case ai.EatSensor:
			fmt.Println("LOOK HIT OWN EAT SENSOR")

		default:
			fmt.Printf("LOOK HIT UNKNOWN SHAPE: %#v\n", hit.Shape)
		}
	}

	return []float32{
		distance / 300,
		data.Appearance[0],
		data.Appearance[1],
		data.Appearance[2],
	}

}
func (ai *AI) GetBody() object.Object {
	return ai.body
}

func NewAI(brain string, body object.Object, space *cp.Space, eatSensor *cp.Shape) AI {
	return AI{
		brain:     brain,
		body:      body,
		space:     space,
		EatSensor: eatSensor,
	}
}

func (ai *AI) GetInputs() []float32 {
	left, right := ai.Look()
	body := ai.body.Shape.Body()
	velocity := body.Velocity()

	angle := body.Angle()

	forwardX := math.Cos(angle)
	forwardY := math.Sin(angle)

	speed := velocity.X*forwardX + velocity.Y*forwardY
	normalizedSpeed := clamp(speed / 137.0)

	normalizedRotationSpeed := clamp(body.AngularVelocity() / 3.5)

	inputs := []float32{
		float32(ai.mood.hunger),
		float32(normalizedSpeed),
		float32(math.Sin(angle)),
		float32(math.Cos(angle)),
		float32(normalizedRotationSpeed),
	}

	inputs = append(inputs, left...)
	inputs = append(inputs, right...)

	return inputs

}

func (ai AI) String() string {
	return ai.brain + ", " + ai.body.String() + ", "
}
