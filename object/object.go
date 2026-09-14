package object

import (
	"github.com/jakecoffman/cp/v2"
)

type Object struct {
	Appearance []int
	Texture    []int
	Shape      *cp.Shape
}

func (o Object) String() string {
	return "ja"
}
