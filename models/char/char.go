package char

import (
	"fmt"
	"image/color"
	"io/ioutil"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/myanagisawa/ebitest/enum"
	"github.com/myanagisawa/ebitest/utils"
	"golang.org/x/image/font"
)

const (
	fontFilePath = "resources/fonts/"
)

var (
	// Res 文字リソース管理変数
	Res   *Resources
	fonts map[enum.FontStyleEnum]font.Face
)

func init() {
	fonts = make(map[enum.FontStyleEnum]font.Face)
	Res = &Resources{
		list: []*Resource{},
	}
}

// FontLoad ...
func FontLoad(style enum.FontStyleEnum, size int) font.Face {
	// すでにロード済ならそれを返す
	if face, ok := fonts[style]; ok {
		return face
	}
	// フォント読み込み
	ftBinary, err := ioutil.ReadFile(fmt.Sprintf("%s%s", fontFilePath, style.Name()))
	if err != nil {
		panic(err)
	}

	tt, err := truetype.Parse(ftBinary)
	if err != nil {
		panic(err)
	}
	face := truetype.NewFace(tt, &truetype.Options{
		Size:    float64(size),
		DPI:     72,
		Hinting: font.HintingFull,
	})
	fonts[style] = face
	return fonts[style]
}

// Resources 文字リソースの管理構造体
type Resources struct {
	list []*Resource
}

// Get 指定のリソースを取得
func (o *Resources) Get(size int, style enum.FontStyleEnum) *Resource {
	for _, r := range o.list {
		if r.fontSize == size && r.fontStyle == style {
			return r
		}
	}
	// リソース未作成なら作成
	r := &Resource{
		fontSize:  size,
		fontStyle: style,
		font:      FontLoad(style, size),
		m:         make(map[rune]*ebiten.Image),
	}
	o.list = append(o.list, r)
	return r
}

// Resource フォントごとの文字リソース
type Resource struct {
	fontSize  int
	fontStyle enum.FontStyleEnum
	font      font.Face
	m         map[rune]*ebiten.Image
}

// GetByRune 指定文字１文字の文字画像を返します
func (o *Resource) GetByRune(r rune) *ebiten.Image {
	if i, ok := o.m[r]; ok {
		return i
	}
	// なかったら作って返す
	img := utils.CreateTextImage(string(r), o.font, color.RGBA{255, 255, 255, 255})
	eimg := ebiten.NewImageFromImage(*img)
	o.m[r] = eimg
	return o.m[r]
}

// GetByString 指定文字列の文字画像を返します
func (o *Resource) GetByString(s string) []*ebiten.Image {
	ret := make([]*ebiten.Image, len([]rune(s)))
	for i, r := range []rune(s) {
		ret[i] = o.GetByRune(r)
	}
	return ret
}
