package base

import "github.com/myanagisawa/ebitest/ebitest"

// PositionImpl ...
type PositionImpl struct {
	position *ebitest.Point
}

// Position ...
func (o *PositionImpl) Position() *ebitest.Point {
	return o.position
}

// ScaleImpl ...
type ScaleImpl struct {
	scale *ebitest.Scale
}

// Scale ...
func (o *ScaleImpl) Scale() *ebitest.Scale {
	return o.scale
}

// AngleImpl ...
type AngleImpl struct {
	angle int
}

// Angle ...
func (o *AngleImpl) Angle() int {
	return o.angle
}
