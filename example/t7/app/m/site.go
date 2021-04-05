package m

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"github.com/myanagisawa/ebitest/example/t7/app/enum"
	"github.com/myanagisawa/ebitest/example/t7/app/g"
	"github.com/myanagisawa/ebitest/example/t7/lib/models/char"
	"github.com/myanagisawa/ebitest/example/t7/lib/utils"
)

// Site マップ上の拠点情報
type Site struct {
	Code     string
	Type     enum.SiteTypeEnum
	Name     string
	Location *g.Point
	Image    image.Image
	Text     image.Image
}

// NewSite ...
func NewSite(code string, t enum.SiteTypeEnum, name string, loc *g.Point) *Site {
	if code == "" || loc.X() > 1 || loc.Y() > 1 {
		panic(fmt.Sprintf("Site初期化エラー: code=%#v, loc=%#v", code, loc))
	}

	// ti := utils.CreateTextImage(name, g.ScaleFonts[12], color.RGBA{32, 32, 32, 255})
	// timg := ebiten.NewImageFromImage(*ti)
	fset := char.Res.Get(12, enum.FontStyleGenShinGothicRegular)
	ti := fset.GetStringImage(name)
	ti = utils.TextColorTo(ti.(draw.Image), &color.RGBA{255, 255, 255, 255})

	return &Site{
		Code:     code,
		Type:     t,
		Name:     name,
		Location: loc,
		Image:    g.Images[fmt.Sprintf("site_%d", t)],
		Text:     ti,
	}
}

// Sites ...
type Sites []Site

// GetByCode ...
func (o *Sites) GetByCode(code string) *Site {
	for _, r := range *o {
		if r.Code == code {
			return &r
		}
	}
	return nil
}

// GetByName ...
func (o *Sites) GetByName(name string) *Site {
	for _, r := range *o {
		if r.Name == name {
			return &r
		}
	}
	return nil
}

// CreateSites ...
func CreateSites(RawSites []RawSite) *Sites {
	sites := Sites{}
	for _, site := range RawSites {
		sites = append(sites, *NewSite(site.Code, site.Type, site.Name, site.Location))
	}
	return &sites
}
