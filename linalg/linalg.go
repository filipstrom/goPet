package linalg

import "math"

type Vec2 struct {
	X float64
	Y float64
}

func (v *Vec2) Magnitude() float64 {
	return math.Hypot(v.X, v.Y)
}

func (v *Vec2) ScalarMult(scalar float64) {
	v.X *= scalar
	v.Y *= scalar
}

func (v *Vec2) Add(vec2 Vec2) {
	v.X += vec2.X
	v.Y += vec2.Y
}

func (v *Vec2) AddDirection(angleRad, magnitude float64) {
	v.X += math.Cos(angleRad) * magnitude
	v.Y += math.Sin(angleRad) * magnitude
}
