package ex3

import (
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type (
	// Unit ...
	Unit interface {
		Scene
		Info() (int, int, int, int, int, int)
		Circle() Circle
		Collision(u *Unit)
	}

	// UnitImpl ...
	UnitImpl struct {
		label      string
		image      ebiten.Image
		circle     *Circle
		vx         float64
		vy         float64
		angle      int
		speed      int
		collision  Unit
		searchArea *Circle
		parent     *Game
	}

	// SearchArea ...
	SearchArea struct {
		area *Circle
	}
)

var (
	maxAngle = 360
)

// NewMyUnit ...
func NewMyUnit(parent *Game) (Unit, error) {
	rand.Seed(time.Now().UnixNano()) //Seed
	// mask画像読み込み
	// mask, _ := utils.GetImageByPath("resources/system_images/mask.png")
	// http://tech.nitoyon.com/ja/blog/2015/12/31/go-image-gen/
	// 座標が円に入っているか
	// http://imagawa.hatenadiary.jp/entry/2016/12/31/190000

	r := 10
	c := &Circle{float64(r), float64(r), r}
	m := image.NewRGBA(image.Rect(0, 0, r*2, r*2))
	for x := 0; x < r*2; x++ {
		for y := 0; y < r*2; y++ {
			if x > r && y >= r-1 && y <= r+1 {
				m.Set(x, y, color.RGBA{0, 0, 0, 255})
			} else {
				d := c.Distance(x, y)
				if d > 1 {
					m.Set(x, y, color.RGBA{0, 0, 0, 0})
				} else {
					m.Set(x, y, color.RGBA{212, 215, 143, 255})
				}
			}
		}
	}

	eimg, err := ebiten.NewImageFromImage(m, ebiten.FilterDefault)
	if err != nil {
		panic(err)
	}

	c.x, c.y = float64(rand.Intn(parent.WindowSize.Width-c.r)), float64(rand.Intn(parent.WindowSize.Height-c.r))
	if int(c.x) < c.r {
		c.x = float64(c.r)
	}
	if int(c.y) < c.r {
		c.y = float64(c.r)
	}

	c.x, c.y = 300, 400
	a := 23
	s := 2
	log.Printf("angle: %d, speed: %d", a, s)

	unitImpl := &UnitImpl{
		label:  "myUnit",
		image:  *eimg,
		circle: c,
		angle:  a,
		speed:  s,
		parent: parent,
	}
	unitImpl.updatePoint()

	return unitImpl, nil
}

// NewUnit ...
func NewUnit(parent *Game) (Unit, error) {
	rand.Seed(time.Now().UnixNano()) //Seed

	r := 100
	c := &Circle{float64(r), float64(r), r}
	m := image.NewRGBA(image.Rect(0, 0, r*2, r*2))
	for x := 0; x < r*2; x++ {
		for y := 0; y < r*2; y++ {
			if x > r && y >= r-1 && y <= r+1 {
				m.Set(x, y, color.RGBA{0, 0, 0, 255})
			} else {
				d := c.Distance(x, y)
				if d > 1 {
					m.Set(x, y, color.RGBA{0, 0, 0, 0})
				} else {
					m.Set(x, y, color.RGBA{212, 0, 0, 255})
				}
			}
		}
	}

	eimg, err := ebiten.NewImageFromImage(m, ebiten.FilterDefault)
	if err != nil {
		panic(err)
	}

	c.x, c.y = 600, 200
	a := 0
	s := 0
	log.Printf("angle: %d, speed: %d", a, s)

	unitImpl := &UnitImpl{
		label:  "coin2",
		image:  *eimg,
		circle: c,
		angle:  a,
		speed:  s,
		parent: parent,
	}
	unitImpl.updatePoint()

	return unitImpl, nil
}

// NewDebris ...
func NewDebris(speed int, parent *Game) (Unit, error) {
	rand.Seed(time.Now().UnixNano()) //Seed

	rd, gr, bl := uint8(rand.Intn(55)+200), uint8(rand.Intn(55)+200), uint8(rand.Intn(55)+200)

	r := rand.Intn(80) + 20
	c := &Circle{float64(r), float64(r), r}
	m := image.NewRGBA(image.Rect(0, 0, r*2, r*2))
	for x := 0; x < r*2; x++ {
		for y := 0; y < r*2; y++ {
			if x > r && y >= r-1 && y <= r+1 {
				m.Set(x, y, color.RGBA{0, 0, 0, 255})
			} else {
				d := c.Distance(x, y)
				if d > 1 {
					m.Set(x, y, color.RGBA{0, 0, 0, 0})
				} else {
					m.Set(x, y, color.RGBA{rd, gr, bl, 255})
				}
			}
		}
	}

	eimg, err := ebiten.NewImageFromImage(m, ebiten.FilterDefault)
	if err != nil {
		panic(err)
	}

	c.x, c.y = float64(rand.Intn(parent.WindowSize.Width-c.r)), float64(rand.Intn(parent.WindowSize.Height-c.r))
	if int(c.x) < c.r {
		c.x = float64(c.r)
	}
	if int(c.y) < c.r {
		c.y = float64(c.r)
	}

	// ebitenのrotateとtranslateはy軸0が最上段なので注意
	a := rand.Intn(maxAngle)
	// a := 45
	log.Printf("angle: %d, speed: %d", a, speed)

	unitImpl := &UnitImpl{
		image:  *eimg,
		circle: c,
		angle:  a,
		speed:  speed,
		parent: parent,
	}
	unitImpl.updatePoint()

	return unitImpl, nil
}

// Update ...
func (s *UnitImpl) Update() error {
	s.updatePoint()

	s.circle.x += s.vx
	s.circle.y -= s.vy

	w := s.parent.WindowSize.Width
	if s.circle.Left() < 0 || w <= s.circle.Right() {
		s.angle = 180 - s.angle
		s.updatePoint()
	}
	h := s.parent.WindowSize.Height
	if s.circle.Top() < 0 || h <= s.circle.Bottom() {
		s.angle = 360 - s.angle
		s.updatePoint()
	}

	// 索敵範囲
	r := 100
	area := &Circle{
		x: s.circle.x,
		y: s.circle.y,
		r: r,
	}
	s.searchArea = area
	// if s.label == "myUnit" {
	// 	log.Printf("s.circle: %#v, s.searchArea: %#v", s.circle, s.searchArea)
	// }

	return nil
}

// Draw ...
func (s *UnitImpl) Draw(r *ebiten.Image) {

	w, h := s.image.Size()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Reset()
	op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	op.GeoM.Rotate(-2 * math.Pi * float64(s.angle) / float64(maxAngle))
	op.GeoM.Translate(float64(w)/2, float64(h)/2)
	op.GeoM.Translate(float64(s.circle.Left()), float64(s.circle.Top()))
	if s.collision != nil {
		op.ColorM.Scale(0.5, 0.5, 0.5, 1.0)
	}
	r.DrawImage(&s.image, op)

	if s.collision != nil {
		// draw line
		t := s.collision

		x1, x2, y1, y2 := s.circle.x, t.Circle().x, s.circle.y, t.Circle().y
		sx, lx, sy, ly := math.Min(x1, x2), math.Max(x1, x2), math.Min(y1, y2), math.Max(y1, y2)
		x, y := lx-sx, ly-sy
		n := math.Atan2(y, x)
		d := n * 180 / math.Pi
		log.Printf("atan2(%f, %f)=%f, deg=%f", y, x, n, d)
		log.Printf("  sx=%f, lx=%f, sy=%f, ly=%f, ", sx, lx, sy, ly)

		ebitenutil.DrawLine(r, x1, y1, x2, y2, color.RGBA{0, 255, 0, 255})
		// draw line

		s.collision = nil
	}

	// 索敵範囲を描画
	if s.label == "myUnit" {
		drawSearchArea(s, r)
	}
}

func drawSearchArea(s *UnitImpl, r *ebiten.Image) {
	sr := s.searchArea.r
	m := image.NewRGBA(image.Rect(0, 0, sr*2, sr*2))
	for x := 0; x < sr*2; x++ {
		for y := 0; y < sr*2; y++ {
			d := s.searchArea.Distance(x, y)
			if d > 1 {
				m.Set(x, y, color.RGBA{0, 0, 0, 0})
			} else {
				m.Set(x, y, color.RGBA{143, 215, 212, 128})
			}
		}
	}

	eimg, err := ebiten.NewImageFromImage(m, ebiten.FilterDefault)
	if err != nil {
		panic(err)
	}
	w, h := eimg.Size()
	op := &ebiten.DrawImageOptions{}
	log.Printf("w=%d, h=%d, left=%d, top=%d", w, h, s.searchArea.Left(), s.searchArea.Top())
	op.GeoM.Reset()
	op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	op.GeoM.Translate(float64(w)/2, float64(h)/2)
	op.GeoM.Translate(float64(s.searchArea.Left()), float64(s.searchArea.Top()))
	r.DrawImage(eimg, op)
}

// GetSize ...
func (s *UnitImpl) GetSize() (int, int) {
	return s.circle.r, s.circle.r
}

// Info ...
func (s *UnitImpl) Info() (int, int, int, int, int, int) {
	return s.angle, s.speed, int(s.vx), int(s.vy), int(s.circle.x), int(s.circle.y)
}

// Circle ...
func (s *UnitImpl) Circle() Circle {
	return *s.circle
}

// Collision ...
func (s *UnitImpl) Collision(c *Unit) {
	s.collision = *c
}

// Distance ...
func (c *Circle) Distance(x, y int) float64 {
	var dx, dy int = int(c.r) - x, int(c.r) - y
	return math.Sqrt(float64(dx*dx+dy*dy)) / float64(c.r)
}

// Left ...
func (c *Circle) Left() int {
	return int(c.x) - c.r
}

// Top ...
func (c *Circle) Top() int {
	return int(c.y) - c.r
}

// Right ...
func (c *Circle) Right() int {
	return int(c.x) + c.r
}

// Bottom ...
func (c *Circle) Bottom() int {
	return int(c.y) + c.r
}

// Width ...
func (c *Circle) Width() int {
	return 2 * c.r
}

// Height ...
func (c *Circle) Height() int {
	return 2 * c.r
}

func (s *UnitImpl) updatePoint() {
	rad := float64(s.angle) * math.Pi / 180
	s.vx = math.Cos(rad) * float64(s.speed)
	s.vy = math.Sin(rad) * float64(s.speed)
	// if s.label == "myCoin" {
	// 	log.Printf("rad=%f, vx=%f, vy=%f, sin(rad)=%f", rad, s.vx, s.vy, math.Sin(rad))
	// }
}
