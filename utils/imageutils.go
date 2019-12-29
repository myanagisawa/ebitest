package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/minodisk/go-fix-orientation/processor"

	"golang.org/x/image/draw"
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
