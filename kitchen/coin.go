package kitchen

import (
	"image"
	"math"

	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type (
	// Coin ...
	Coin interface {
		Scene
	}

	// CoinImpl ...
	CoinImpl struct {
		image  ebiten.Image
		x      int
		y      int
		width  int
		height int
	}

	// Circle ...
	Circle struct {
		X, Y, R float64
	}
)

// NewCoin ...
func NewCoin() (Coin, error) {
	// mask画像読み込み
	// mask, _ := utils.GetImageByPath("resources/system_images/mask.png")
	// http://tech.nitoyon.com/ja/blog/2015/12/31/go-image-gen/
	// 座標が円に入っているか
	// http://imagawa.hatenadiary.jp/entry/2016/12/31/190000
	// img := image.NewRGBA(image.Rect(0, 0, 50, 50))
	// bounds := img.Bounds()
	// for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
	// 	for x := bounds.Min.X; x < bounds.Max.X; x++ {
	// 		img.Set(x, y, color.RGBA{212, 215, 143, 255})
	// 	}
	// }
	// // 画像をマスク
	// out := utils.MaskImage(img, mask)

	var w, h int = 280, 240
	var hw, hh float64 = float64(w / 2), float64(h / 2)
	r := 40.0
	θ := 2 * math.Pi / 3
	cr := &Circle{hw - r*math.Sin(0), hh - r*math.Cos(0), 60}
	cg := &Circle{hw - r*math.Sin(θ), hh - r*math.Cos(θ), 60}
	cb := &Circle{hw - r*math.Sin(-θ), hh - r*math.Cos(-θ), 60}

	m := image.NewRGBA(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			r := cr.Brightness(float64(x), float64(y))
			g := cg.Brightness(float64(x), float64(y))
			b := cb.Brightness(float64(x), float64(y))
			a := r
			if a < g {
				a = g
			}
			if a < b {
				a = b
			}
			c := color.RGBA{r, g, b, a}
			m.Set(x, y, c)
		}
	}

	eimg, err := ebiten.NewImageFromImage(m, ebiten.FilterDefault)
	if err != nil {
		panic(err)
	}

	return &CoinImpl{
		image: *eimg,
	}, nil
}

// Update ...
func (s *CoinImpl) Update() error {
	return nil
}

// Draw ...
func (s *CoinImpl) Draw(r *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	r.DrawImage(&s.image, op)
}

// Brightness ...
func (c *Circle) Brightness(x, y float64) uint8 {
	var dx, dy float64 = c.X - x, c.Y - y
	d := math.Sqrt(dx*dx+dy*dy) / c.R
	if d > 1 {
		return 0
	}
	return uint8((1 - math.Pow(d, 5)) * 255)
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
