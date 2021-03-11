package g

import (
	"fmt"
	"image"
	"log"
	"math"
	"path/filepath"

	"github.com/myanagisawa/ebitest/example/t7/lib/utils"
)

var (
	// Width ...
	Width int
	// Height ...
	Height int

	// Images ...
	Images map[string]image.Image
)

func init() {

	Images = make(map[string]image.Image)
	exe, err := filepath.Abs(".")
	if err != nil {
		panic(err.Error())
	}
	log.Printf("global.AppRoot: exe=%s", exe)

	Images["btnBase"], _ = utils.GetImageByPath(fmt.Sprintf("%s/%s", exe, "../../resources/system_images/button.png"))
	Images["world"], _ = utils.GetImageByPath(fmt.Sprintf("%s/%s", exe, "../../resources/system_images/world.jpg"))
}

// Point ...
type Point struct {
	x float64
	y float64
}

// NewPoint ...
func NewPoint(x, y float64) *Point {
	return &Point{x, y}
}

// DefPoint ...
func DefPoint() *Point {
	return &Point{0.0, 0.0}
}

// GoString ...
func (p Point) GoString() string {
	return fmt.Sprintf("x=%0.3f, y=%0.3f", p.x, p.y)
}

// X ...
func (p *Point) X() float64 {
	return p.x
}

// Y ...
func (p *Point) Y() float64 {
	return p.y
}

// IntX ...
func (p *Point) IntX() int {
	return int(p.x)
}

// IntY ...
func (p *Point) IntY() int {
	return int(p.y)
}

// Get ...
func (p *Point) Get() (float64, float64) {
	return p.x, p.y
}

// GetInt ...
func (p *Point) GetInt() (int, int) {
	return int(p.x), int(p.y)
}

// Set ...
func (p *Point) Set(x, y float64) {
	p.x = x
	p.y = y
}

// SetDelta ...
func (p *Point) SetDelta(dx, dy float64) *Point {
	p.x += dx
	p.y += dy
	return &Point{p.x, p.y}
}

// Scale op.scale設定値
type Scale struct {
	x float64
	y float64
}

// NewScale ...
func NewScale(x, y float64) *Scale {
	return &Scale{x, y}
}

// DefScale ...
func DefScale() *Scale {
	return &Scale{1.0, 1.0}
}

// X ...
func (s *Scale) X() float64 {
	return s.x
}

// Y ...
func (s *Scale) Y() float64 {
	return s.y
}

// Get ...
func (s *Scale) Get() (float64, float64) {
	return s.x, s.y
}

// Set ...
func (s *Scale) Set(x, y float64) {
	s.x = x
	s.y = y
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Size ...
type Size struct {
	width  int
	height int
}

// NewSize ...
func NewSize(width, height int) *Size {
	return &Size{width, height}
}

// W ...
func (s *Size) W() int {
	return s.width
}

// H ...
func (s *Size) H() int {
	return s.height
}

// Get ...
func (s *Size) Get() (int, int) {
	return s.width, s.height
}

// Set ...
func (s *Size) Set(width, height int) {
	s.width = width
	s.height = height
}

// Angle ...
type Angle struct {
	rad float64
}

// NewAngle ...
func NewAngle(rad float64) *Angle {
	return &Angle{rad: rad}
}

// NewAngleByDeg ...
func NewAngleByDeg(degree int) *Angle {
	rad := float64(degree) * (math.Pi / 180)
	return &Angle{rad: rad}
}

// DefAngle ...
func DefAngle() *Angle {
	return &Angle{rad: 0}
}

// Get ...
func (a *Angle) Get() float64 {
	return a.rad
}

// GetDeg ...
func (a *Angle) GetDeg() float64 {
	return math.Pi * a.rad / 180.0
}

// Set ...
func (a *Angle) Set(rad float64) {
	a.rad = rad
	if a.rad > 2*math.Pi {
		a.rad -= 2 * math.Pi
	} else if a.rad < 0 {
		a.rad += 2 * math.Pi
	}
}

// SetDelta ...
func (a *Angle) SetDelta(degree int) {
	rad := float64(degree) * (math.Pi / 180)
	a.rad += rad
	if a.rad > 2*math.Pi {
		a.rad -= 2 * math.Pi
	} else if a.rad < 0 {
		a.rad += 2 * math.Pi
	}
}

// Bound ...
type Bound struct {
	Min Point
	Max Point
}

// NewBound ...
func NewBound(min, max *Point) *Bound {
	return &Bound{Min: *min, Max: *max}
}

// NewBoundByPosSize ...
func NewBoundByPosSize(pos *Point, size *Size) *Bound {
	return &Bound{Min: *pos, Max: *NewPoint(pos.X()+float64(size.W()), pos.Y()+float64(size.H()))}
}

// SetDelta ...
func (b *Bound) SetDelta(min *Point, max *Point) {
	if min != nil {
		b.Min.SetDelta(min.X(), min.Y())
		b.Max.SetDelta(min.X(), min.Y())
	}
	if max != nil {
		b.Max.SetDelta(max.X(), max.Y())
	}
}

// ToPosSize ...
func (b *Bound) ToPosSize() (*Point, *Size) {
	return &b.Min, NewSize(b.Max.IntX()-b.Min.IntX(), b.Max.IntY()-b.Min.IntY())
}

// ToImageRect ...
func (b *Bound) ToImageRect() *image.Rectangle {
	return &image.Rectangle{
		Min: image.Point{b.Min.IntX(), b.Min.IntY()},
		Max: image.Point{b.Max.IntX(), b.Max.IntY()},
	}
}

// ColorScale ...
type ColorScale struct {
	r float64
	g float64
	b float64
	a float64
}

// NewCS ...
func NewCS(r, g, b, a float64) *ColorScale {
	return &ColorScale{r, g, b, a}
}

// Get ...
func (s *ColorScale) Get() (float64, float64, float64, float64) {
	return s.r, s.g, s.b, s.a
}

// Set ...
func (s *ColorScale) Set(r, g, b, a float64) {
	s.r = r
	s.g = g
	s.b = b
	s.a = a
}

// DefCS ...
func DefCS() *ColorScale {
	return &ColorScale{1.0, 1.0, 1.0, 1.0}
}

// Vector ...
type Vector struct {
	amount float64
	angle  Angle
}

// NewVector ...
func NewVector(am float64, an *Angle) *Vector {
	return &Vector{amount: am, angle: *an}
}

// GetAmount ...
func (v *Vector) GetAmount() float64 {
	return v.amount
}

// GetDelta ...
func (v *Vector) GetDelta() *Point {
	vx := math.Cos(v.angle.rad) * v.amount
	vy := math.Sin(v.angle.rad) * v.amount
	return NewPoint(vx, vy)
}

// Angle ...
func (v *Vector) Angle() *Angle {
	return &v.angle
}
