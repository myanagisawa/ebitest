package char

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"log"

	"github.com/golang/freetype/truetype"
	"github.com/myanagisawa/ebitest/example/t7/app/enum"
	"github.com/myanagisawa/ebitest/example/t7/lib/utils"
	"golang.org/x/image/font"
)

const (
	fontFilePath = "../../resources/fonts/"
)

var (
	// Res 文字リソース管理変数
	Res   *Resources
	fonts map[enum.FontStyleEnum]*truetype.Font
	faces map[enum.FontStyleEnum]map[int]font.Face
)

func init() {
	fonts = make(map[enum.FontStyleEnum]*truetype.Font)
	faces = make(map[enum.FontStyleEnum]map[int]font.Face)
	Res = &Resources{
		list: []*Resource{},
	}
}

// FontLoad ...
func FontLoad(style enum.FontStyleEnum, size int) font.Face {
	if face, ok := faces[style][size]; ok {
		// すでにロード済ならそれを返す
		return face
	}
	var fd *truetype.Font
	if f, ok := fonts[style]; ok {
		fd = f
	} else {
		// フォント読み込み
		ftBinary, err := ioutil.ReadFile(fmt.Sprintf("%s%s", fontFilePath, style.Name()))
		if err != nil {
			panic(err)
		}

		tt, err := truetype.Parse(ftBinary)
		if err != nil {
			panic(err)
		}
		fd = tt
	}
	face := truetype.NewFace(fd, &truetype.Options{
		Size:    float64(size),
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if _, ok := faces[style]; !ok {
		faces[style] = make(map[int]font.Face)
	}
	faces[style][size] = face
	log.Printf("faces: %#v", faces)
	return faces[style][size]
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
		m:         make(map[rune]image.Image),
	}
	o.list = append(o.list, r)
	return r
}

// Resource フォントごとの文字リソース
type Resource struct {
	fontSize  int
	fontStyle enum.FontStyleEnum
	font      font.Face
	m         map[rune]image.Image
}

// GetByRune 指定文字１文字の文字画像を返します
func (o *Resource) GetByRune(r rune) image.Image {
	if i, ok := o.m[r]; ok {
		return i
	}
	// なかったら作って返す
	img := utils.CreateTextImage(string(r), o.font, color.RGBA{255, 255, 255, 255})
	o.m[r] = *img
	return o.m[r]
}

// GetByString 指定文字列の文字画像を返します
func (o *Resource) GetByString(s string) []image.Image {
	ret := make([]image.Image, len([]rune(s)))
	for i, r := range []rune(s) {
		ret[i] = o.GetByRune(r)
	}
	return ret
}

// GetStringImage 指定文字列の文字画像を返します
func (o *Resource) GetStringImage(s string) image.Image {
	images := make([]image.Image, len([]rune(s)))
	for i, r := range []rune(s) {
		images[i] = o.GetByRune(r)
	}
	// テキスト画像の幅を取得
	w, h := 0, 0
	for i := range images {
		image := images[i]
		size := image.Bounds().Size()
		w += size.X
		if h < size.Y {
			h = size.Y
		}
	}
	base := utils.CreateRectImage(w, h, &color.RGBA{0, 0, 0, 0}).(draw.Image)
	tx := 0
	for i := range images {
		img := images[i]
		base = utils.StackImage(base, img, image.Point{tx, 0})
		tx += img.Bounds().Size().X
	}

	return base
}
