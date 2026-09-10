package object

import "github.com/filipstrom/goPet/linalg"

type Object struct {
	Appearance []int
	Texture    []int
	Position   linalg.Vec2
}
