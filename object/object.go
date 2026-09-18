package object

import (
	"github.com/jakecoffman/cp/v2"
)

type ShapeData struct {
	Appearance []float32
	IsFood     bool
}

type Object struct {
	Appearance []int
	Texture    []int
	Shape      *cp.Shape
}

func (o Object) String() string {
	return "ja"
}
