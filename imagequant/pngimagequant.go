package imagequant

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L. -limagequant
#include "libimagequant.h"
#include "pngimagequant.h"
#
*/
import "C"

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	_ "image/png"
	"io/fs"
	"io/ioutil"
	"log"
	"unsafe"
)

func Ping() string {
	return "Ping will return Pong"
}

func LiqAttrCreate() {

	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)

	pngBuffer, err := Read("benchmark_image_1.png")

	if err != nil {
		log.Println(err)
		return
	}

	png, format, errImg := image.Decode(bytes.NewReader(pngBuffer))

	if errImg != nil {
		log.Println(errImg)
		return
	}

	fmt.Println("got:", format)

	pixels := imageToRGBA(png)
	fmt.Println(len(pixels.Pix))

	version := C.liq_version()
	// defer C.free(unsafe.Pointer(version))
	fmt.Println(version)

	//rawRGBAPixels := (*C.uchar)(unsafe.Pointer(&pixels.Pix[0]))
	pngPixelPtr := &pixels.Pix[0]
	rawRGBAPixels := (*C.uchar)(unsafe.Pointer(pngPixelPtr))
	fmt.Printf("%v\n", rawRGBAPixels)
	// defer C.free(unsafe.Pointer(rawRGBAPixels))
	// data := (*C.)
	quantResult := C.pngQuant{}
	quantResult = C.doQuant(rawRGBAPixels, C.uint(png.Bounds().Size().X), C.uint(png.Bounds().Size().Y), C.double(0))
	// defer C.free((C.pngQuant).unsafe.Pointer(quantResult))
	fmt.Println(quantResult.Status)
	fmt.Println(quantResult.Size)

	encodedPng := C.GoBytes(unsafe.Pointer(quantResult.Png), C.int(quantResult.Size))

	fmt.Println(len(encodedPng))
	err = Write("benchmark_image_1_quant.png", encodedPng)
	if err != nil {
		log.Println(err)
	}

}

// imageToRGBA returns RGBA having Pix []uint8
func imageToRGBA(src image.Image) *image.RGBA {

	// No conversion needed if image is an *image.RGBA.
	if dst, ok := src.(*image.RGBA); ok {
		return dst
	}

	// Use the image/draw package to convert to *image.RGBA.
	b := src.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(dst, dst.Bounds(), src, b.Min, draw.Src)
	return dst
}

// Get the bi-dimensional pixel array
func getPixels(img image.Image) ([][]Pixel, error) {

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	var pixels [][]Pixel
	for y := 0; y < height; y++ {
		var row []Pixel
		for x := 0; x < width; x++ {
			row = append(row, rgbaToPixel(img.At(x, y).RGBA()))
		}
		pixels = append(pixels, row)
	}

	return pixels, nil
}

// img.At(x, y).RGBA() returns four uint32 values; we want a Pixel
func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
	return Pixel{int(r / 257), int(g / 257), int(b / 257), int(a / 257)}
}

// Pixel struct example
type Pixel struct {
	R int
	G int
	B int
	A int
}

// Read an image file.
func Read(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func Write(path string, data []byte) error {
	return ioutil.WriteFile(path, data, fs.FileMode(0640))
}
