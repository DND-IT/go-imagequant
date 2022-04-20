package imagequant_test

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"testing"

	"github.com/DND-IT/go-imagequant/imagequant"
)

func TestQualityRun(t *testing.T) {
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)

	urlList, err := readUrls()
	urlListLength := len(urlList)

	if err != nil {
		t.Fatal(err)
	}

	var deltas = []float64{}

	// set to true for visual inspection
	var keepAllFiles = false

	var wg sync.WaitGroup
	wg.Add(urlListLength)

	for i := 0; i < urlListLength; i++ {
		go func(i int) {
			defer wg.Done()

			r, err := processFile(urlList, i)
			if err != nil {
				return
			}
			delta := r.calulateDelta()

			t.Logf("processing %s", urlList[i])

			// keep critical images - otherwise remove
			if delta > 0 {
				t.Logf("file %s got bigger: %f%%", r.QuantitizedFilename, delta)
			} else {
				if !keepAllFiles {
					os.Remove(r.OriginalFilename)
					os.Remove(r.QuantitizedFilename)
				}
			}
			deltas = append(deltas, delta)
		}(i)
	}

	wg.Wait()

	var totalDelta float64 = 0

	for j := 0; j < len(deltas); j++ {
		totalDelta += deltas[j]
	}
	var averageDelta float64 = totalDelta / float64(len(deltas))

	t.Logf("average reduction in percent %f%%", averageDelta)

}

func processFile(urlList []string, i int) (*encodeResult, error) {
	url := urlList[i]

	// download image to local temp file and decode into image
	fileIn, err := ioutil.TempFile("/tmp", fmt.Sprintf("image%v-in-*.png", i))

	if err != nil {
		return nil, err
	}
	err = downloadFile(fileIn.Name(), url)

	if err != nil {
		return nil, err
	}

	inputImage, err := ioutil.ReadFile(fileIn.Name())
	if err != nil {
		return nil, err
	}

	img, _, errDecode := image.Decode(bytes.NewReader(inputImage))
	if errDecode != nil {
		return nil, errDecode
	}

	// quantisation

	q, errQ := imagequant.New(img, 0, 0, 100, imagequant.DefaultSpeed)
	if errQ != nil {
		return nil, errQ
	}

	quantImg, errR := q.Run()
	if errR != nil {
		return nil, errR
	}

	buff := bytes.NewBuffer([]byte{})

	// encode as png and write to temp file
	err = png.Encode(buff, quantImg)
	if err != nil {
		return nil, err
	}

	fileOut, err := ioutil.TempFile("/tmp", fmt.Sprintf("image%v-out-*.png", i))
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	fw := bufio.NewWriter(fileOut)

	fw.Write(buff.Bytes())

	r := &encodeResult{OriginalSize: len(inputImage), QuantitizedSize: len(buff.Bytes()), OriginalFilename: fileIn.Name(), QuantitizedFilename: fileOut.Name()}

	return r, nil
}

func readUrls() ([]string, error) {
	f, err := os.OpenFile("./test/urls.txt", os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("open file error: %v", err)
		return nil, err
	}
	defer f.Close()

	var urls []string

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		urls = append(urls, sc.Text())
	}
	if err := sc.Err(); err != nil {
		log.Fatalf("scan file error: %v", err)
		return nil, err
	}

	return urls, nil

}

func downloadFile(filepath string, url string) (err error) {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return
}

type encodeResult struct {
	OriginalSize        int
	OriginalFilename    string
	QuantitizedFilename string
	QuantitizedSize     int
}

func (m encodeResult) calulateDelta() (delta float64) {

	diff := float64(m.QuantitizedSize - m.OriginalSize)
	delta = (diff / float64(m.OriginalSize)) * 100
	return
}
