package parts

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/myanagisawa/ebitest/example/t5/ebitest"
	"github.com/myanagisawa/ebitest/example/t5/enum"
	"github.com/myanagisawa/ebitest/utils"
)

// Site マップ上の拠点情報
type Site struct {
	Code     string
	Type     enum.SiteTypeEnum
	Name     string
	Location *ebitest.Point
	Image    *ebiten.Image
	Text     *ebiten.Image
}

// NewSite ...
func NewSite(code string, t enum.SiteTypeEnum, name string, loc *ebitest.Point) *Site {
	if code == "" || loc.X() > 1 || loc.Y() > 1 {
		panic(fmt.Sprintf("Site初期化エラー: code=%#v, loc=%#v", code, loc))
	}

	img := ebiten.NewImageFromImage(ebitest.Images[fmt.Sprintf("site_%d", t)])

	ti := utils.CreateTextImage(name, ebitest.ScaleFonts[12], color.RGBA{32, 32, 32, 255})
	timg := ebiten.NewImageFromImage(*ti)

	return &Site{
		Code:     code,
		Type:     t,
		Name:     name,
		Location: loc,
		Image:    img,
		Text:     timg,
	}
}
