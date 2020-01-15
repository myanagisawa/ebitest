package imgconv

import (
	"bytes"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"

	"golang.org/x/image/draw"

	"github.com/minodisk/go-fix-orientation/processor"

	"github.com/myanagisawa/ebitest/utils"
)

// CreateNewImages ...
func CreateNewImages(orgImgDir, imgDir string) {
	files, err := ioutil.ReadDir(orgImgDir)
	if err != nil {
		panic(err)
	}
	wg := &sync.WaitGroup{}
	for _, file := range files {
		if _, err := os.Stat(fmt.Sprintf("%s/%s", imgDir, file.Name())); os.IsNotExist(err) {
			wg.Add(1)
			fmt.Print(".")
			go func(fname string) {
				err = ri(fname, orgImgDir, imgDir)
				if err != nil {
					log.Println(err.Error())
					_ = os.Remove(fmt.Sprintf("%s/%s", imgDir, file.Name()))
				} else {
					_ = os.Remove(fmt.Sprintf("%s/%s", orgImgDir, file.Name()))
				}
				wg.Done()
			}(file.Name())
		}
	}
	wg.Wait()
}

func ri(fname, orgImgDir, imgDir string) error {
	file, err := os.Open(fmt.Sprintf("%s/%s", orgImgDir, fname))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	defer file.Close()

	// 画像読み込み
	img, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("read file error: %s", err.Error())
		return err
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
		return err
	}
	// log.Printf("image type: %s", t)

	// log.Printf("image resize")
	rctSrc := s.Bounds()
	imgDst := image.NewRGBA(image.Rect(0, 0, rctSrc.Dx()/3, rctSrc.Dy()/3))
	draw.CatmullRom.Scale(imgDst, imgDst.Bounds(), s, rctSrc, draw.Over, nil)

	var d image.Image
	if o > 0 {
		// log.Printf("image orientation")
		d = processor.ApplyOrientation(imgDst, o)
	} else {
		// log.Printf("Unnecessary orientation")
		d = imgDst
	}

	//create resized image file
	dst, err := os.Create(fmt.Sprintf("%s/%s", imgDir, fname))
	if err != nil {
		return err
	}
	defer dst.Close()

	ext := filepath.Ext(fname)[1:]
	err = utils.EncodeImage(dst, d, ext)
	if err != nil {
		return err
	}
	log.Println("finish")
	return nil
}
