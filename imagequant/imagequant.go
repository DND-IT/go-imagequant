package imagequant

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"unsafe"
)

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L. -limagequant
#include "stdlib.h"
#include "libimagequant.h"
#include "go-imagequant.h"
*/
import "C"

import (
	"image/draw"
)

const (
	DefaultSpeed      = 4
	DefaultMinQuality = 0
	DefaultMaxQuality = 100
)

type QImg struct {
	Img        image.Image
	ImgRGBA    *image.RGBA
	Gamma      float64
	MinQuality uint // default 0
	MaxQuality uint // default 100
	Speed      uint // range allowed between 1 and 10, default is 4

}

// GetLiqVersion returns version if libimagequant used.
func GetLiqVersion() int {
	return int(C.int(C.liq_version()))
}

func New(img image.Image, gamma float64, minQuality, maxQuality, speed uint) (*QImg, error) {

	// validate
	if minQuality > maxQuality {
		return nil, errors.New("min quality can not be bigger than max quality")
	}

	if maxQuality > 100 {
		return nil, errors.New("max quality allowed is 100")
	}

	if speed < 1 {
		return nil, errors.New("min allowed speed is 1")
	}

	if speed > 10 {
		return nil, errors.New("max allowed speed is 10")
	}

	q := QImg{
		Img:        img,
		ImgRGBA:    nil,
		Gamma:      gamma,
		MinQuality: minQuality,
		MaxQuality: maxQuality,
		Speed:      speed,
	}

	q.ImgRGBA = ImageToRGBA(img)

	return &q, nil
}

// Run call c lib imagequant functions needed for quantize an RGBA image.
func (q *QImg) Run() (image.Image, error) {

	if q.ImgRGBA == nil {
		return nil, fmt.Errorf("can not quant nil image")
	}

	// get ptr for first slice item
	pixelPtr := &q.ImgRGBA.Pix[0]
	// create unsafe unsigned char pointer needed for C
	ptrToRawRGBAPixels := (*C.uchar)(unsafe.Pointer(pixelPtr))

	handle := C.liq_attr_create()
	defer C.liq_attr_destroy_wrapper(handle)

	liqError := C.liq_set_speed(handle, C.int(q.Speed))
	if liqError != C.LIQ_OK {
		return nil, fmt.Errorf("c call to liq_set_speed() failed with code %v", liqError)
	}

	if q.MaxQuality != DefaultMaxQuality || q.MinQuality != DefaultMinQuality {
		liqError = C.liq_set_quality(handle, C.int(q.MinQuality), C.int(q.MaxQuality))
		if liqError != C.LIQ_OK {
			return nil, fmt.Errorf("c call to liq_set_speed() failed with code %v", liqError)
		}
	}

	// liq_image *input_image = liq_image_create_rgba(handle, raw_rgba_pixels, (int) width, (int) height, gamma);
	cWidth := C.int(q.ImgRGBA.Bounds().Size().X)
	cHeight := C.int(q.ImgRGBA.Bounds().Size().Y)
	cGamma := C.double(q.Gamma)

	inputImage := C.liq_image_create_rgba_wrapper(handle, ptrToRawRGBAPixels, cWidth, cHeight, cGamma)
	defer C.liq_image_destroy(inputImage)

	var liqResult *C.liq_result
	defer C.liq_result_destroy(liqResult)

	liqError = C.liq_image_quantize(inputImage, handle, &liqResult)

	if liqError != C.LIQ_OK {
		return nil, fmt.Errorf("c call to liq_image_quantize() failed with code %v", liqError)
	}

	pixelSize := q.ImgRGBA.Bounds().Size().X * q.ImgRGBA.Bounds().Size().Y

	// alloc memory needed to liq_write_remapped_image
	cRaw8BitPixels := C.CBytes(make([]uint8, pixelSize))
	defer C.free(cRaw8BitPixels) // be sure to release C alloc memory

	// call c lib to write the new
	C.liq_write_remapped_image(liqResult, inputImage, cRaw8BitPixels, C.ulong(pixelSize))

	// get palette
	cPalette := C.liq_get_palette(liqResult)

	var outPalette color.Palette

	palCount := uint(C.int(cPalette.count))

	// iterate the palette received from c lib imagequant and
	// create a go color palette.
	for i := uint(0); i < palCount; i++ {
		col := color.RGBA{
			R: uint8(C.uint(cPalette.entries[i].r)),
			G: uint8(C.uint(cPalette.entries[i].g)),
			B: uint8(C.uint(cPalette.entries[i].b)),
			A: uint8(C.uint(cPalette.entries[i].a)),
		}
		outPalette = append(outPalette, col)
	}

	// create new go palette image
	qImg := image.NewPaletted(q.ImgRGBA.Rect, outPalette)

	// copy unsigned chars from c lib imagequant alloc memory into go []uint8
	qImg.Pix = C.GoBytes(cRaw8BitPixels, C.int(pixelSize))

	return qImg, nil
}

// ImageToRGBA returns RGBA image.
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

/**
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
*/
