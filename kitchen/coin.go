package kitchen

import (
	"image"
	"math"
	"math/rand"

	"image/color"

	"github.com/hajimehoshi/ebiten"
)

type (
	// Coin ...
	Coin interface {
		Scene
	}

	// CoinImpl ...
	CoinImpl struct {
		image ebiten.Image
		x     int
		y     int
		vx    int
		vy    int
		angle int
	}

	// Circle ...
	Circle struct {
		X, Y, R float64
	}
)

var (
	maxAngle = 256
)

// NewCoin ...
func NewCoin() (Coin, error) {
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
			if y == int(r) {
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
	vx, vy := 2*rand.Intn(2)-1, 2*rand.Intn(2)-1
	a := rand.Intn(maxAngle)

	return &CoinImpl{
		image: *eimg,
		x:     x,
		y:     y,
		vx:    vx,
		vy:    vy,
		angle: a,
	}, nil
}

// Update ...
func (s *CoinImpl) Update() error {
	w, h := s.image.Size()

	s.x += s.vx
	s.y += s.vy
	if s.x < 0 {
		s.x = -s.x
		s.vx = -s.vx
	} else if 1344 <= s.x+w {
		s.x = 2*(1344-w) - s.x
		s.vx = -s.vx
	}
	if s.y < 0 {
		s.y = -s.y
		s.vy = -s.vy
	} else if 1008 <= s.y+h {
		s.y = 2*(1008-h) - s.y
		s.vy = -s.vy
	}
	s.angle++
	s.angle %= maxAngle

	return nil
}

// Draw ...
func (s *CoinImpl) Draw(r *ebiten.Image) {
	w, h := s.image.Size()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Reset()
	op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	op.GeoM.Rotate(2 * math.Pi * float64(s.angle) / float64(maxAngle))
	op.GeoM.Translate(float64(w)/2, float64(h)/2)
	op.GeoM.Translate(float64(s.x), float64(s.y))
	r.DrawImage(&s.image, op)
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
