package kitchen

import (
	"image"
	"math"
	"math/rand"
	"time"

	"image/color"

	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type (
	// Coin ...
	Coin interface {
		Scene
		Info() (int, int, int, int, int, int)
		Circle() Circle
		Collision(c *Coin)
	}

	// CoinImpl ...
	CoinImpl struct {
		label     string
		image     ebiten.Image
		circle    *Circle
		vx        float64
		vy        float64
		angle     int
		speed     int
		collision Coin
	}

	// Circle ...
	Circle struct {
		x, y float64
		r    int
	}
)

var (
	maxAngle = 360
)

// NewMyCoin ...
func NewMyCoin() (Coin, error) {
	rand.Seed(time.Now().UnixNano()) //Seed
	// mask画像読み込み
	// mask, _ := utils.GetImageByPath("resources/system_images/mask.png")
	// http://tech.nitoyon.com/ja/blog/2015/12/31/go-image-gen/
	// 座標が円に入っているか
	// http://imagawa.hatenadiary.jp/entry/2016/12/31/190000

	r := 40
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

	c.x, c.y = float64(rand.Intn(1344-c.r)), float64(rand.Intn(1008-c.r))
	if int(c.x) < c.r {
		c.x = float64(c.r)
	}
	if int(c.y) < c.r {
		c.y = float64(c.r)
	}

	// ebitenのrotateとtranslateはy軸0が最上段なので注意
	// a := rand.Intn(maxAngle)
	// s := 10

	c.x, c.y = 300, 400
	a := 23
	s := 2
	log.Printf("angle: %d, speed: %d", a, s)

	coinImpl := &CoinImpl{
		label:  "myCoin",
		image:  *eimg,
		circle: c,
		angle:  a,
		speed:  s,
	}
	coinImpl.updatePoint()

	return coinImpl, nil
}

// NewCoin ...
func NewCoin() (Coin, error) {
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

	coinImpl := &CoinImpl{
		label:  "coin2",
		image:  *eimg,
		circle: c,
		angle:  a,
		speed:  s,
	}
	coinImpl.updatePoint()

	return coinImpl, nil
}

// NewDebris ...
func NewDebris(speed int) (Coin, error) {
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

	c.x, c.y = float64(rand.Intn(1344-c.r)), float64(rand.Intn(1008-c.r))
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

	coinImpl := &CoinImpl{
		image:  *eimg,
		circle: c,
		angle:  a,
		speed:  speed,
	}
	coinImpl.updatePoint()

	return coinImpl, nil
}

// Update ...
func (s *CoinImpl) Update() error {
	// c := rand.Intn(300)
	// if c == 0 {
	// 	s.angle += 5
	// 	log.Printf("change angle: %d", s.angle)
	// } else if c == 19 {
	// 	s.angle -= 5
	// 	log.Printf("change angle: %d", s.angle)
	// }

	s.updatePoint()

	s.circle.x += s.vx
	s.circle.y -= s.vy

	if s.circle.Left() < 0 || 1344 <= s.circle.Right() {
		s.angle = 180 - s.angle
		s.updatePoint()
	}
	if s.circle.Top() < 0 || 1008 <= s.circle.Bottom() {
		s.angle = 360 - s.angle
		s.updatePoint()
	}

	return nil
}

// Draw ...
func (s *CoinImpl) Draw(r *ebiten.Image) {
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
}

// Info ...
func (s *CoinImpl) Info() (int, int, int, int, int, int) {
	return s.angle, s.speed, int(s.vx), int(s.vy), int(s.circle.x), int(s.circle.y)
}

// Circle ...
func (s *CoinImpl) Circle() Circle {
	return *s.circle
}

// Collision ...
func (s *CoinImpl) Collision(c *Coin) {
	s.collision = *c
}

// Distance ...
func (c *Circle) Distance(x, y int) float64 {
	var dx, dy int = int(c.x) - x, int(c.y) - y
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

func (s *CoinImpl) updatePoint() {
	rad := float64(s.angle) * math.Pi / 180
	s.vx = math.Cos(rad) * float64(s.speed)
	s.vy = math.Sin(rad) * float64(s.speed)
	// if s.label == "myCoin" {
	// 	log.Printf("rad=%f, vx=%f, vy=%f, sin(rad)=%f", rad, s.vx, s.vy, math.Sin(rad))
	// }
}
