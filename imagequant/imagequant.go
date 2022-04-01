package imagequant

import "C"
import (
	"fmt"
	"image"
	"image/color"
	"log"
	"unsafe"
)

//
// https://gist.github.com/zchee/b9c99695463d8902cd33

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L. -limagequant
#include "libimagequant.h"
#include "go-imagequant.h"
*/
import "C"

import (
	"image/draw"
	_ "image/png"
	"io/fs"
	"io/ioutil"
)

func Ping() string {
	return "Ping will return Pong"
}

// GetLiqVersion returns version if libimagequant used.
func GetLiqVersion() int {
	return int(C.int(C.liq_version()))
}

// Quant call c lib.
func Quant(imgRGBA *image.RGBA, gamma float64) (palImg image.Image, err error) {

	if imgRGBA == nil {
		return nil, fmt.Errorf("can not quant nil image")
	}

	// get ptr for first slice item
	pixelPtr := &imgRGBA.Pix[0]
	// create unsafe unsigned char pointer needed for C
	rawRGBAPixels := (*C.uchar)(unsafe.Pointer(pixelPtr))
	// defer C.free(unsafe.Pointer(rawRGBAPixels))

	quantResult := C.imgQuant{}
	quantResult = C.doQuant(rawRGBAPixels,
		C.uint(imgRGBA.Bounds().Size().X),
		C.uint(imgRGBA.Bounds().Size().Y),
		C.double(gamma),
	)
	// defer C.free((C.pngQuant).unsafe.Pointer(quantResult))
	// fmt.Println(quantResult.Status)
	// fmt.Println(quantResult.Size)

	// copy unsigned char from c into go []uint8

	// create a new palette
	var outPalette color.Palette
	// get liq_palette struct from c

	// C.free(unsafe.Pointer(quantResult.Pixels))
	palCount := uint(C.int(quantResult.palette.count))
	log.Println(palCount)
	// iterate the palette received from lib imagequant
	for i := uint(0); i < palCount; i++ {
		col := color.RGBA{
			R: uint8(C.uint(quantResult.palette.entries[i].r)),
			G: uint8(C.uint(quantResult.palette.entries[i].g)),
			B: uint8(C.uint(quantResult.palette.entries[i].b)),
			A: uint8(C.uint(quantResult.palette.entries[i].a)),
		}
		outPalette = append(outPalette, col)
	}

	// create new palette image
	pImg := image.NewPaletted(imgRGBA.Rect, outPalette)

	// copy unsigned chars from lib imagequant into go []uint8
	pImg.Pix = C.GoBytes(unsafe.Pointer(quantResult.pixels), C.int(quantResult.size))

	// free C alloc
	// C.destroyImgQuant(quantResult)

	return pImg, nil
}

// ImageToRGBA returns RGBA having Pix []uint8
func ImageToRGBA(src image.Image) *image.RGBA {

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
