package imagequant_test

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"sync"
	"testing"

	"go-imagequant/imagequant"
)

func BenchmarkRun(b *testing.B) {
	// register image formats
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
	// image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
	// image.RegisterFormat("gif", "gif", gif.Decode, gif.DecodeConfig)

	rawBenchImage, err := read("./test/benchmark_image_1.png")
	if err != nil {
		b.Fatal(err)
	}

	img, _, errDecode := image.Decode(bytes.NewReader(rawBenchImage))
	if errDecode != nil {
		b.Fatal(errDecode)
	}

	rgbaBenchImage := imagequant.ImageToRGBA(img)

	concurrencyLevels := []int{5, 10, 20, 50}
	for _, clients := range concurrencyLevels {
		b.Run(fmt.Sprintf("%d_clients", clients), func(b *testing.B) {
			// sem := make(chan struct{}, clients)
			wg := sync.WaitGroup{}
			for n := 0; n < b.N; n++ {
				wg.Add(1)
				go func() {
					// b.Log(clients, n)
					_, err := imagequant.Run(rgbaBenchImage, 0)
					if err != nil {
						b.Error(err)
					}
					// <-sem
					wg.Done()
				}()
			}
			wg.Wait()
		})
	}
}

// read an image file.
func read(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}
