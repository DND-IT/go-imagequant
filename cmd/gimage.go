package main

import (
	"bytes"
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
)

var (
	ErrImageBufferPtrIsNil = errors.New("image byte buffer pointer can not be nil")
	// ErrImagePointerIsNil          = errors.New("image struct pointer can not be nil")
	// ErrImageOperationNotSupported = errors.New("image operation not supported")
)

// DecodeConfig is a wrapper for image.DecodeConfig().
func DecodeConfig(buffer *[]byte) (image.Config, string, error) {
	if buffer == nil {
		return image.Config{}, "", ErrImageBufferPtrIsNil
	}
	r := bytes.NewReader(*buffer)

	return image.DecodeConfig(r)
}

func Decode(buffer *[]byte) (image.Image, string, error) {
	if buffer == nil {
		return nil, "", ErrImageBufferPtrIsNil
	}
	r := bytes.NewReader(*buffer)
	return image.Decode(r)
}

/*
func DecodeAllGIF(buffer *[]byte) (*gif.GIF, error) {
	if buffer == nil {
		return nil, ErrImageBufferPtrIsNil
	}
	r := bytes.NewReader(*buffer)
	g, err := gif.DecodeAll(r)
	if err != nil {
		return nil, err
	}
	return g, nil
} */

// Encode is a wrapper for various image encoder.
func Encode(img image.Image, imageTypeName string) ([]byte, error) {
	var (
		buff = new(bytes.Buffer)
		err  error
	)

	switch imageTypeName {
	case "png":
		err = png.Encode(buff, img)
	case "jpeg":
		err = jpeg.Encode(buff, img, &jpeg.Options{Quality: jpeg.DefaultQuality})
	case "gif":
		err = gif.Encode(buff, img, nil)
	//case "webp":
	//	var options *encoder.Options
	//
	//	options, err = encoder.NewLossyEncoderOptions(encoder.PresetDefault, cfg.ImageOptions.WebpQuality)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	err = webp.Encode(buff, dst, options)
	default:
		var encoder png.Encoder
		encoder.CompressionLevel = png.BestSpeed
		err = encoder.Encode(buff, img)
	}

	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}
