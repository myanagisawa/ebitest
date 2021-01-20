package base

import "github.com/myanagisawa/ebitest/app/g"

// PositionImpl ...
type PositionImpl struct {
	position *g.Point
}

// Position ...
func (o *PositionImpl) Position() *g.Point {
	return o.position
}

// ScaleImpl ...
type ScaleImpl struct {
	scale *g.Scale
}

// Scale ...
func (o *ScaleImpl) Scale() *g.Scale {
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
