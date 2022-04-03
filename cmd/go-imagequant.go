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
	showLibImageQuantVersion := flag.Bool("showLibImageQuantVersion", false, "show lib image quant version and exit")

	flag.Parse()

	if *showLibImageQuantVersion {
		fmt.Println(quant.GetLiqVersion())
		return
	}

	if *imageSrcPath == "" {
		fmt.Println("no src image")
		return
	}

	if *imageDstPath == "" {
		fmt.Println("no dst image")
		return
	}

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

	img, imgType, errDecode := Decode(&imageBuff)

	if errDecode != nil {
		log.Println(err)
		return
	}

	switch imgType {
	case "gif":

	}

	rgbaImg := quant.ImageToRGBA(img)
	qImg, errQuant := quant.Run(rgbaImg, 0)
	// fmt.Printf("raw pixels size:%d\n", len(qImg.Pix))
	//rgbaQuantImg, errQuant := quant.Quant(rgbaImg, 0)
	if errQuant != nil {
		log.Println(errQuant)
		return
	}

	out, errEncode := Encode(qImg, imgType)
	if errEncode != nil {
		log.Println(errEncode)
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
