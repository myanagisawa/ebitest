package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	"github.com/minodisk/go-fix-orientation/processor"
)

// OrientationImage 画像の向き補正処理
func OrientationImage(path string) (image.Image, error) {

	file, err := os.Open(path) //maybe file path
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}
	defer file.Close()

	// 画像読み込み
	img, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("read file error: %s", err.Error())
		return nil, err
	}

	o, err := processor.ReadOrientation(bytes.NewReader(img))
	if err != nil {
		log.Printf("ReadOrientation error: %s", err.Error())
	} else {
		log.Printf("image orientation: %d", o)
	}

	s, _, err := image.Decode(bytes.NewReader(img))
	if err != nil {
		log.Printf("image decode error: %s", err.Error())
		return nil, err
	}
	// log.Printf("image type: %s", t)

	var d image.Image
	if o > 0 {
		// log.Printf("image orientation")
		d = processor.ApplyOrientation(s, o)
	} else {
		// log.Printf("Unnecessary orientation")
		d = s
	}
	return d, err
}

// ScaleImage ...
func ScaleImage(i image.Image, w, h int) (image.Image, error) {
	rctSrc := i.Bounds()
	ratio := float64(rctSrc.Dx()) / float64(rctSrc.Dy())
	// log.Println("Width:", rctSrc.Dx())
	// log.Println("Height:", rctSrc.Dy())
	// log.Printf("ratio: %v\n", ratio)
	var rctDraw image.Rectangle
	if ratio < 1 {
		// 幅基準で縦の移動量を計算
		t := (rctSrc.Dy() - rctSrc.Dx()) / 2
		rctDraw = image.Rect(0, t, rctSrc.Dx(), rctSrc.Dy()-t)
	} else {
		// 高さ基準で横の移動量を計算
		l := (rctSrc.Dx() - rctSrc.Dy()) / 2
		rctDraw = image.Rect(l, 0, rctSrc.Dx()-l, rctSrc.Dy())
	}

	imgDst := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.CatmullRom.Scale(imgDst, imgDst.Bounds(), i, rctDraw, draw.Over, nil)
	return imgDst, nil
}

// EncodeImage ...
func EncodeImage(buf io.Writer, i image.Image, ext string) error {

	log.Printf("ext=%s", strings.ToLower(ext))
	switch strings.ToLower(ext) {
	case "jpeg", "jpg":
		if err := jpeg.Encode(buf, i, &jpeg.Options{Quality: 90}); err != nil {
			log.Printf("jpeg encode error: %s", err.Error())
			return err
		}
	case "gif":
		if err := gif.Encode(buf, i, nil); err != nil {
			log.Printf("gif encode error: %s", err.Error())
			return err
		}
	case "png":
		if err := png.Encode(buf, i); err != nil {
			log.Printf("png encode error: %s", err.Error())
			return err
		}
	default:
		err := fmt.Errorf("format error")
		log.Printf(err.Error())
		return err
	}
	return nil
}

// GetImageByPath 指定したパスからimageを取得します
func GetImageByPath(path string) (draw.Image, string) {
	//画像読み込み
	fileIn, err := os.Open(path)
	defer fileIn.Close()
	if err != nil {
		fmt.Println("error:file\n", err)
		log.Panic(err.Error())
	}

	//画像をimage型として読み込む
	img, format, err := image.Decode(fileIn)
	if err != nil {
		fmt.Println("error:decode\n", format, err)
		log.Panic(err.Error())
	}
	out := image.NewRGBA(img.Bounds())
	draw.Copy(out, image.Point{}, img, img.Bounds(), draw.Src, nil)
	return out, format
}

// MaskImage srcをmaskした結果を返します
func MaskImage(src draw.Image, mask image.Image) draw.Image {
	rmask := image.NewRGBA(src.Bounds())
	draw.CatmullRom.Scale(rmask, rmask.Bounds(), mask, mask.Bounds(), draw.Over, nil)
	// 円形maskの適用
	out := image.NewRGBA(src.Bounds())
	draw.DrawMask(out, out.Bounds(), src, image.Point{0, 0}, rmask, image.Point{0, 0}, draw.Over)
	return out
}

// DrawFont 指定imageの中心にstr文字列を描画します
func DrawFont(out draw.Image, str string, fontsize float64) {
	ft, err := truetype.Parse(fontload("/Library/Fonts/Arial Unicode.ttf"))
	if err != nil {
		fmt.Println("font", err)
		return
	}
	opt := truetype.Options{Size: fontsize}
	face := truetype.NewFace(ft, &opt)

	d := &font.Drawer{
		Dst:  out,
		Src:  image.NewUniform(color.White),
		Face: face,
	}

	// 文字を表示対象の真ん中に表示する
	size := out.Bounds().Size()
	d.Dot.X = (fixed.I(size.X) - d.MeasureString(str)) / 2
	d.Dot.Y = fixed.I((size.Y / 2) + int(fontsize/2))

	d.DrawString(str)
}

func fontload(fname string) []byte {
	file, err := os.Open(fname)
	defer file.Close()
	if err != nil {
		fmt.Println("error:file\n", err)
		return nil
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("error:fileread\n", err)
		return nil
	}

	return bytes
}

// SetTextToCenter dstの中心にtextを配置します
func SetTextToCenter(text string, src draw.Image, face font.Face, c color.Color) *draw.Image {
	out := image.NewRGBA(src.Bounds())
	draw.Copy(out, image.Point{}, src, src.Bounds(), draw.Src, nil)
	d := &font.Drawer{
		Dst:  out,
		Src:  image.NewUniform(c),
		Face: face,
	}
	d.Dot.X = (fixed.I(out.Bounds().Max.X) - d.MeasureString(text)) / 2
	d.Dot.Y = fixed.I((out.Bounds().Max.Y + face.Metrics().Height.Ceil()) / 2)
	d.DrawString(text)
	return &d.Dst
}
