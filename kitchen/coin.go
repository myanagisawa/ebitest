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
	}

	// CoinImpl ...
	CoinImpl struct {
		image ebiten.Image
		x     int
		y     int
		vx    int
		vy    int
		angle int
		speed int
	}

	// Circle ...
	Circle struct {
		X, Y, R float64
	}
)

var (
	maxAngle = 360
)

// NewCoin ...
func NewCoin() (Coin, error) {
	rand.Seed(time.Now().UnixNano()) //Seed
	// mask画像読み込み
	// mask, _ := utils.GetImageByPath("resources/system_images/mask.png")
	// http://tech.nitoyon.com/ja/blog/2015/12/31/go-image-gen/
	// 座標が円に入っているか
	// http://imagawa.hatenadiary.jp/entry/2016/12/31/190000

	r := 40.0
	c := &Circle{r, r, r}
	m := image.NewRGBA(image.Rect(0, 0, int(r*2), int(r*2)))
	for x := 0; x < int(r*2); x++ {
		for y := 0; y < int(r*2); y++ {
			if x > int(r) && y == int(r) {
				m.Set(x, y, color.RGBA{0, 0, 0, 255})
			} else {
				d := c.Distance(float64(x), float64(y))
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
	x, y := rand.Intn(1344-w), rand.Intn(1008-h)

	// ebitenのrotateとtranslateはy軸0が最上段なので注意
	a := rand.Intn(maxAngle)
	s := rand.Intn(10) + 1
	// a := 45
	// s := 7
	log.Printf("angle: %d, speed: %d", a, s)

	coinImpl := &CoinImpl{
		image: *eimg,
		x:     x,
		y:     y,
		angle: a,
		speed: s,
	}
	coinImpl.updatePoint()

	return coinImpl, nil
}

// Update ...
func (s *CoinImpl) Update() error {
	w, h := s.image.Size()

	s.x += s.vx
	s.y -= s.vy

	if s.x < 0 || 1344 <= s.x+w {
		s.angle = 180 - s.angle
		s.updatePoint()
	}
	if s.y < 0 || 1008 <= s.y+h {
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
	op.GeoM.Translate(float64(s.x), float64(s.y))
	r.DrawImage(&s.image, op)
}

// Info ...
func (s *CoinImpl) Info() (int, int, int, int, int, int) {
	return s.angle, s.speed, s.vx, s.vy, s.x, s.y
}

// Brightness ...
func (c *Circle) Brightness(x, y float64) uint8 {
	var dx, dy float64 = c.X - x, c.Y - y
	d := math.Sqrt(dx*dx+dy*dy) / c.R
	if d > 1 {
		return 0
	}
	return uint8((1 - d) * 255)
}

// Inside ...
func (c *Circle) Inside(x, y float64) bool {
	var dx, dy float64 = c.X - x, c.Y - y
	return math.Sqrt(dx*dx+dy*dy) <= c.R
}

// Distance ...
func (c *Circle) Distance(x, y float64) float64 {
	var dx, dy float64 = c.X - x, c.Y - y
	return math.Sqrt(dx*dx+dy*dy) / c.R
}

func (s *CoinImpl) updatePoint() {
	rad := float64(s.angle) * math.Pi / 180
	s.vx = int(math.Cos(rad) * float64(s.speed))
	s.vy = int(math.Sin(rad) * float64(s.speed))
}
