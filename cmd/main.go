package main

import (
	"flag"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/fs"
	"io/ioutil"
	"log"

	_ "golang.org/x/image/webp"

	quant "go-imagequant/imagequant"
)

func main() {
	imageSrcPath := flag.String("src", "", "src image path")
	imageDstPath := flag.String("dst", "", "dst image path")

	flag.Parse()

	if *imageSrcPath == "" {
		fmt.Println("no src image")
		return
	}

	if *imageDstPath == "" {
		fmt.Println("no dst image")
		return
	}

	// fmt.Println(quant.Ping())
	liqVersion := quant.GetLiqVersion()
	log.Printf("using libimagequant version %d\n", liqVersion)

	imageBuff, err := Read(*imageSrcPath)

	if err != nil {
		fmt.Println(err)
		return
	}

	// register image formats
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
	image.RegisterFormat("gif", "gif", gif.Decode, gif.DecodeConfig)
	// image.RegisterFormat("webp", webp, webp.Decode, webp.DecodeConfig)

	// try to decode image
	_, _, err = DecodeConfig(&imageBuff)
	if err != nil {
		fmt.Println(err)
		return
	}

	img, _, errDecode := Decode(&imageBuff)

	if errDecode != nil {
		log.Println(err)
		return
	}

	rgbaImg := quant.ImageToRGBA(img)
	fmt.Printf("raw pixels size:%d\n", len(rgbaImg.Pix))
	rgbaQuantImg, errQuant := quant.Quant(rgbaImg, 0)
	if errQuant != nil {
		log.Println(errQuant)
		return
	}

	out, errPng := Encode(rgbaQuantImg, "png")
	if errPng != nil {
		log.Println(errPng)
	}

	_ = Write(*imageDstPath, out)

}

// Read an image file.
func Read(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

// Write an image file.
func Write(path string, data []byte) error {
	return ioutil.WriteFile(path, data, fs.FileMode(0640))
}
