package main

import (
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten"
)

type (
	// Object ...
	Object struct {
		label string
		image image.Image
		ei    *ebiten.Image
		x, y  float64
		angle int
		speed int
	}
	// Circle ...
	Circle struct {
		Object
		r int
	}
)

// NewCircle ...
func NewCircle(r int) Circle {
	c := Circle{
		r: r,
	}
	c.Object = Object{
		x:     0,
		y:     0,
		angle: 0,
		speed: 0,
	}
	c.SetImage()

	return c
}

// SetOption ...
func (o *Object) SetOption(op *ebiten.DrawImageOptions) {
	w, h := o.ei.Size()
	op.GeoM.Reset()
	op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	op.GeoM.Rotate(-2 * math.Pi * float64(o.angle) / float64(maxAngle))
	op.GeoM.Translate(float64(w)/2, float64(h)/2)
	op.GeoM.Translate(float64(o.x), float64(o.y))
}

// EImage ...
func (c *Circle) EImage() *ebiten.Image {
	if c.ei != nil {
		return c.ei
	}
	eimg, err := ebiten.NewImageFromImage(c.image, ebiten.FilterDefault)
	if err != nil {
		panic(err)
	}
	c.ei = eimg
	return eimg
}

// SetImage ...
func (c *Circle) SetImage() {
	w, h := c.r*2, c.r*2
	m := image.NewRGBA(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			d := c.Distance(x, y)
			if d <= 1.0 && d > 0.5 {
				m.Set(x, y, color.RGBA{212, 215, 143, 255})
			} else {
				m.Set(x, y, color.RGBA{0, 0, 0, 0})
			}
		}
	}
	c.image = m
	c.EImage()
}

// Distance ...
func (c *Circle) Distance(x, y int) float64 {
	var dx, dy int = c.r - x, c.r - y
	return math.Sqrt(float64(dx*dx+dy*dy)) / float64(c.r)
}

var (
	c1 Circle

	maxAngle = 360
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func update(r *ebiten.Image) error {
	c1.angle += 2
	if c1.angle > maxAngle {
		c1.angle = 0
	}

	op := &ebiten.DrawImageOptions{}
	c1.SetOption(op)
	r.DrawImage(c1.EImage(), op)
	return nil
}

func main() {
	c1 = NewCircle(10)
	if err := ebiten.Run(update, 200, 100, 1.0, "t2"); err != nil {
		log.Fatal(err)
	}
}
