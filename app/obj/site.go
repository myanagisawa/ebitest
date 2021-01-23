package obj

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/myanagisawa/ebitest/app/g"
	"github.com/myanagisawa/ebitest/enum"
)

// Site マップ上の拠点情報
type Site struct {
	Code     string
	Type     enum.SiteTypeEnum
	Name     string
	Location *g.Point
	Image    *ebiten.Image
	Text     *ebiten.Image
}

// NewSite ...
func NewSite(code string, t enum.SiteTypeEnum, name string, loc *g.Point) *Site {
	if code == "" || loc.X() > 1 || loc.Y() > 1 {
		panic(fmt.Sprintf("Site初期化エラー: code=%#v, loc=%#v", code, loc))
	}

	img := ebiten.NewImageFromImage(g.Images[fmt.Sprintf("site_%d", t)])

	// ti := utils.CreateTextImage(name, g.ScaleFonts[12], color.RGBA{32, 32, 32, 255})
	// timg := ebiten.NewImageFromImage(*ti)

	return &Site{
		Code:     code,
		Type:     t,
		Name:     name,
		Location: loc,
		Image:    img,
		// Text:     timg,
	}
}

// Sites ...
type Sites []Site

// GetSiteByCode ...
func (o *Sites) GetSiteByCode(code string) *Site {
	for _, site := range *o {
		if site.Code == code {
			return &site
		}
	}
	return nil
}

// CreateSites ...
func CreateSites(siteMaster []g.MSite) Sites {
	sites := Sites{}
	for _, site := range siteMaster {
		sites = append(sites, *NewSite(site.Code, site.Type, site.Name, site.Location))
	}
	return sites
}
