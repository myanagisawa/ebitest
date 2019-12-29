package main

import (
	"bytes"
	"flag"
	"image"
	_ "image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime/pprof"
	"sync"

	"fmt"

	"golang.org/x/image/draw"

	"github.com/minodisk/go-fix-orientation/processor"

	"github.com/hajimehoshi/ebiten"
	"github.com/myanagisawa/ebitest/ebitest"
	"github.com/myanagisawa/ebitest/utils"
)

const (
	orgImgDir = "images"
	imgDir    = "resized_images"
)

var cpuProfile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			log.Fatal(err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal(err)
		}
		defer pprof.StopCPUProfile()
	}

	// 新着イメージの変換
	createNewImages()

	ebiten.SetRunnableInBackground(true)

	game, _ := ebitest.NewGame(getResourceNames())
	update := game.Update
	if err := ebiten.Run(update, ebitest.ScreenWidth, ebitest.ScreenHeight, 1, "ebitest"); err != nil {
		log.Fatal(err)
	}
}

func getImages(dir string) []string {
	files, _ := ioutil.ReadDir(dir)
	fnames := make([]string, len(files))
	for i, f := range files {
		fnames[i] = f.Name()
	}
	return fnames
}

func getResourceNames() []string {
	fnames := getImages(imgDir)
	paths := make([]string, len(fnames))
	for i, f := range fnames {
		paths[i] = fmt.Sprintf("%s/%s", imgDir, f)
	}
	return paths
}

func createNewImages() {
	images := getImages(orgImgDir)
	wg := &sync.WaitGroup{}
	for _, image := range images {
		if _, err := os.Stat(fmt.Sprintf("%s/%s", imgDir, image)); os.IsNotExist(err) {
			wg.Add(1)
			fmt.Print(".")
			go func(fname string) {
				err = ri(fname)
				if err != nil {
					log.Println(err.Error())
					_ = os.Remove(fmt.Sprintf("%s/%s", imgDir, image))
				} else {
					_ = os.Remove(fmt.Sprintf("%s/%s", orgImgDir, image))
				}
				wg.Done()
			}(image)
		}
	}
	wg.Wait()
}

func ri(fname string) error {
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

// func resizeImages(paths []string) {
// 	wg := &sync.WaitGroup{}
// 	for _, path := range paths {
// 		wg.Add(1) // wgをインクリメント
// 		go func(p string) {
// 			file, err := os.Open(p) //maybe file path
// 			if err != nil {
// 				fmt.Fprintln(os.Stderr, err)
// 				panic(err)
// 			}
// 			defer file.Close()

// 			// 画像読み込み
// 			img, err := ioutil.ReadAll(file)
// 			if err != nil {
// 				log.Printf("read file error: %s", err.Error())
// 				panic(err)
// 			}

// 			s, _, err := image.Decode(bytes.NewReader(img))
// 			if err != nil {
// 				log.Printf("image decode error: %s", err.Error())
// 				panic(err)
// 			}

// 			rctSrc := s.Bounds()
// 			imgDst := image.NewRGBA(image.Rect(0, 0, rctSrc.Dx()/3, rctSrc.Dy()/3))
// 			draw.CatmullRom.Scale(imgDst, imgDst.Bounds(), s, rctSrc, draw.Over, nil)

// 			//create resized image file
// 			dst, err := os.Create(fmt.Sprintf("resized_%s", p))
// 			if err != nil {
// 				panic(err)
// 			}
// 			defer dst.Close()

// 			err = utils.EncodeImage(dst, imgDst, "jpeg")
// 			if err != nil {
// 				panic(err)
// 			}
// 			log.Println("finish")
// 			wg.Done() // 完了したのでwgをデクリメント
// 		}(path)
// 	}
// 	wg.Wait()
// }

// // 画像のリサイズ処理
// func resizeImage(path string, w, h int) (image.Image, error) {

// 	file, err := os.Open(path) //maybe file path
// 	if err != nil {
// 		fmt.Fprintln(os.Stderr, err)
// 		return nil, err
// 	}
// 	defer file.Close()

// 	// 画像読み込み
// 	img, err := ioutil.ReadAll(file)
// 	if err != nil {
// 		log.Printf("read file error: %s", err.Error())
// 		return nil, err
// 	}

// 	o, err := processor.ReadOrientation(bytes.NewReader(img))
// 	if err != nil {
// 		log.Printf("ReadOrientation error: %s", err.Error())
// 	} else {
// 		log.Printf("image orientation: %d", o)
// 	}

// 	s, _, err := image.Decode(bytes.NewReader(img))
// 	if err != nil {
// 		log.Printf("image decode error: %s", err.Error())
// 		return nil, err
// 	}
// 	// log.Printf("image type: %s", t)

// 	// log.Printf("image resize")
// 	s, err = utils.ScaleImage(s, w, h)
// 	if err != nil {
// 		log.Printf("image scale error: %s", err.Error())
// 		return nil, err
// 	}

// 	var d image.Image
// 	if o > 0 {
// 		// log.Printf("image orientation")
// 		d = processor.ApplyOrientation(s, o)
// 	} else {
// 		// log.Printf("Unnecessary orientation")
// 		d = s
// 	}
// 	return d, err
// }
