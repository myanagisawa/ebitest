package kitchen

import (
	"image"
	"math"
	"math/rand"
	"time"

	"image/color"

	"log"

	"github.com/hajimehoshi/ebiten"
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
		image     ebiten.Image
		circle    *Circle
		vx        int
		vy        int
		angle     int
		speed     int
		collision Coin
	}

	// Circle ...
	Circle struct {
		x, y, r int
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
	c := &Circle{r, r, r}
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

	w, h := eimg.Size()
	c.x, c.y = rand.Intn(1344-w), rand.Intn(1008-h)

	// ebitenのrotateとtranslateはy軸0が最上段なので注意
	a := rand.Intn(maxAngle)
	s := rand.Intn(10) + 1
	// a := 45
	// s := 7
	log.Printf("angle: %d, speed: %d", a, s)

	coinImpl := &CoinImpl{
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
	c := &Circle{r, r, r}
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

	w, h := eimg.Size()
	c.x, c.y = rand.Intn(1344-w), rand.Intn(1008-h)

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
	c := rand.Intn(300)
	if c == 0 {
		s.angle += 5
		log.Printf("change angle: %d", s.angle)
	} else if c == 19 {
		s.angle -= 5
		log.Printf("change angle: %d", s.angle)
	}

	s.updatePoint()

	s.circle.x += s.vx
	s.circle.y -= s.vy

	w, h := s.image.Size()
	if s.circle.x < 0 || 1344 <= s.circle.x+w {
		s.angle = 180 - s.angle
		s.updatePoint()
	}
	if s.circle.y < 0 || 1008 <= s.circle.y+h {
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
	// op.GeoM.Rotate(float64(s.angle))
	op.GeoM.Translate(float64(w)/2, float64(h)/2)
	op.GeoM.Translate(float64(s.circle.x), float64(s.circle.y))
	if s.collision != nil {
		op.ColorM.Scale(0.5, 0.5, 0.5, 1.0)
		s.collision = nil
	}
	r.DrawImage(&s.image, op)
}

// Info ...
func (s *CoinImpl) Info() (int, int, int, int, int, int) {
	return s.angle, s.speed, s.vx, s.vy, s.circle.x, s.circle.y
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
	var dx, dy int = c.x - x, c.y - y
	return math.Sqrt(float64(dx*dx+dy*dy)) / float64(c.r)
}

func (s *CoinImpl) updatePoint() {
	rad := float64(s.angle) * math.Pi / 180
	s.vx = int(math.Cos(rad) * float64(s.speed))
	s.vy = int(math.Sin(rad) * float64(s.speed))
}
