package wmap

import (
	"image/color"
	"log"

	"github.com/myanagisawa/ebitest/app/g"
	"github.com/myanagisawa/ebitest/app/obj"
	"github.com/myanagisawa/ebitest/enum"
	"github.com/myanagisawa/ebitest/interfaces"
	"github.com/myanagisawa/ebitest/models/control"
	"github.com/myanagisawa/ebitest/models/layer"
)

var (
	siteInfo   interfaces.UIDialog
	siteInfoDs *obj.Site
)

// infoLayer ...
type infoLayer struct {
	layer.Base
}

// newInfoLayer ...
func newInfoLayer(w, h int) *infoLayer {
	log.Printf("newInfoLayer(%d, %d)", w, h)
	scrollProg = nil

	l := layer.NewLayerBase("info", g.NewPoint(0, 0), g.NewSize(w, h), &color.RGBA{0, 0, 0, 0}, false).(*layer.Base)
	il := &infoLayer{
		Base: *l,
	}

	newSiteInfo(il)

	siSize := siteInfo.Size(enum.TypeScaled)
	x := w - siSize.W() - 5
	y := h - siSize.H() - 5
	siteInfo.SetPosition(float64(x), float64(y))

	// siteInfo.SetVisible(true)

	return il
}

func newSiteInfo(parent *infoLayer) {
	siteInfo = control.NewDialog("site info", g.NewSize(300, 500), false)
	siteInfo.AddItem(control.NewHeaderBar("site info header", siteInfo, true, true))

	siteInfo.SetLayer(parent)
}

func (o *infoLayer) ShowSiteInfo(obj *obj.Site) {
	siteInfoDs = obj
	log.Printf("siteinfo(%0.1f, %0.1f):  datasource: %#v", siteInfo.Position(enum.TypeLocal).X(), siteInfo.Position(enum.TypeLocal).Y(), siteInfoDs)
	siteInfo.SetVisible(true)
}

func (o *infoLayer) HideSiteInfo() {
	siteInfo.SetVisible(false)
}

// GetObjects ...
func (o *infoLayer) GetObjects(x, y int) []interfaces.EbiObject {
	objs := []interfaces.EbiObject{}
	objs = append(objs, siteInfo.GetObjects(x, y)...)
	return objs
}
