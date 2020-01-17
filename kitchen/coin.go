package kitchen

import (
	"image"

	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/myanagisawa/ebitest/utils"
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
)

// NewCoin ...
func NewCoin() (Coin, error) {
	// mask画像読み込み
	mask, _ := utils.GetImageByPath("resources/system_images/mask.png")

	// http://tech.nitoyon.com/ja/blog/2015/12/31/go-image-gen/
	img := image.NewRGBA(image.Rect(0, 0, 50, 50))
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			img.Set(x, y, color.RGBA{212, 215, 143, 255})
		}
	}
	// 画像をマスク
	out := utils.MaskImage(img, mask)
	eimg, err := ebiten.NewImageFromImage(out, ebiten.FilterDefault)
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
