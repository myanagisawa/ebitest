package kitchen

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/myanagisawa/ebitest/utils"
)

type (
	// Sink ...
	Sink struct {
		bgImage *ebiten.Image
	}
)

// NewSink ...
func NewSink(size *Size) (*Sink, error) {
	// mask画像読み込み
	mask, _ := utils.GetImageByPath("resources/system_images/mask.png")

	// 対象画像
	path := "resources/resized_images/IMG_1212.jpg"
	img, _ := utils.GetImageByPath(path)
	utils.ScaleImage(img, size.Width, size.Height)

	// 画像をマスク
	out := utils.MaskImage(img, mask)
	// 文字列を表示
	utils.DrawFont(out, "てすと", 48.0)

	eimg, err := ebiten.NewImageFromImage(out, ebiten.FilterDefault)
	if err != nil {
		panic(err)
	}
	s := &Sink{
		bgImage: eimg,
	}

	return s, nil
}

// Update ...
func (s *Sink) Update() error {
	return nil
}

// Draw ...
func (s *Sink) Draw(r *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.ColorM.Scale(100.0/255.0, 100.0/255.0, 100.0/255.0, 1)
	r.DrawImage(s.bgImage, op)
}
