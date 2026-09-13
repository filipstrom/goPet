package pet

import (
	"math"

	"github.com/filipstrom/goPet/linalg"
	"github.com/filipstrom/goPet/object"
)

type Direction int

const (
	DirectionUp Direction = iota
	DirectionDown
	DirectionLeft
	DirectionRight
)

type Body struct {
	acc              float64
	maxVelocity      float64
	position         linalg.Vec2
	velocity         linalg.Vec2
	startPosition    linalg.Vec2
	rotation         float64
	speed            float64
	maxRotationSpeed float64
	heigt            float64
	width            float64
}

func NewBody(acc, maxVelocity, speed, maxRotationSpeed, heigt, width float64, position linalg.Vec2) Body {
	return Body{acc: acc, maxVelocity: maxVelocity, position: position, startPosition: position, speed: speed, maxRotationSpeed: maxRotationSpeed, heigt: heigt, width: width}
}

func (b Body) GetHeigt() float64 {
	return b.heigt
}

func (b Body) GetWidth() float64 {
	return b.width
}

func sign(x float64) float64 {
	if x > 0 {
		return 1
	}
	if x < 0 {
		return -1
	}
	return 0
}

func (b *Body) DeAcc(force float64) {
	velocity := b.velocity.Magnitude()

	if velocity <= force {
		b.velocity = linalg.Vec2{X: 0, Y: 0}
		return
	}

	scale := (velocity - force) / velocity

	b.velocity.ScalarMult(scale)
}

func (b Body) String() string {
	return "en kropp"
}

func (b *Body) MoveBack() {
	b.position.Subb(b.velocity)
}

func (b *Body) Move() {
	b.position.Add(b.velocity)
}

func (b *Body) Acc(force float64) {
	angleRad := b.rotation * math.Pi / 180
	movement := force * b.acc

	b.velocity.AddDirection(angleRad, movement)

	velocity := b.velocity.Magnitude()

	if velocity > b.maxVelocity {
		scale := b.maxVelocity / velocity
		b.velocity.ScalarMult(scale)
	}
}
func (b *Body) Rotate(force float64, objects []object.Object) {
	//fmt.Println(force)
	//fmt.Println(b.maxRotationSpeed)
	//fmt.Println(b.rotation)

	b.rotation += b.maxRotationSpeed * force

	//fmt.Println(b.rotation)
	if b.rotation >= 360 {
		b.rotation -= 360
	}
	if b.rotation < 0 {
		b.rotation += 360
	}

	collition, collisionPoint := b.CheckBoxCollision(objects)

	if collition {
		b.rotation -= b.maxRotationSpeed * force
		collisionPoint.Subb(b.position)
		collisionPoint.Magnitude()
		collisionPoint.ScalarMult(-0.2)
		b.position.Add(collisionPoint)
		if b.rotation >= 360 {
			b.rotation -= 360
		}
		if b.rotation < 0 {
			b.rotation += 360
		}
	}
	//fmt.Println(b.rotation)
}

func (b *Body) Reset() {
	b.velocity.ScalarMult(0)

	b.position.X = b.startPosition.X
	b.position.Y = b.startPosition.Y
}
func (ai *AI) ResetBody() {
	ai.body.Reset()
}
func (b *Body) CheckBoxCollision(objects []object.Object) (bool, linalg.Vec2) {
	angle := b.rotation * math.Pi / 180

	cos := math.Abs(math.Cos(angle))
	sin := math.Abs(math.Sin(angle))

	halfWidth := (b.width*cos + b.heigt*sin) / 2
	halfHeight := (b.width*sin + b.heigt*cos) / 2

	bodyLeft := b.position.X - halfWidth
	bodyRight := b.position.X + halfWidth
	bodyTop := b.position.Y - halfHeight
	bodyBottom := b.position.Y + halfHeight

	for _, obj := range objects {
		objLeft := obj.Position.X
		objRight := obj.Position.X + obj.Width
		objTop := obj.Position.Y
		objBottom := obj.Position.Y + obj.Height

		if bodyLeft < objRight &&
			bodyRight > objLeft &&
			bodyTop < objBottom &&
			bodyBottom > objTop {

			overlapLeft := math.Max(bodyLeft, objLeft)
			overlapRight := math.Min(bodyRight, objRight)
			overlapTop := math.Max(bodyTop, objTop)
			overlapBottom := math.Min(bodyBottom, objBottom)

			collisionPoint := linalg.Vec2{
				X: (overlapLeft + overlapRight) / 2,
				Y: (overlapTop + overlapBottom) / 2,
			}

			return true, collisionPoint
		}
	}

	return false, linalg.Vec2{X: 0, Y: 0}
}
func (ai *AI) Move(objects []object.Object) {

	ai.body.Move()
	// TODO: Fix so that it does not stop before the collition (it jumps back to far)

	collition, _ := ai.body.CheckBoxCollision(objects)

	for collition {
		ai.body.MoveBack()
		ai.body.DeAcc(0.001)
		ai.body.Move()
		collition, _ = ai.body.CheckBoxCollision(objects)

	}
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
	body  Body
	world string
}

func NewAI(brain string, body Body, world string) AI {
	return AI{
		brain: brain,
		body:  body,
		world: world,
	}
}

func (ai AI) String() string {
	return ai.brain + ", " + ai.body.String() + ", " + ai.world
}
