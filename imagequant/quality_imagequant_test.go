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

	var deltas []float64

	// set to true for visual inspection
	var keepAllFiles = false

	var wg sync.WaitGroup
	wg.Add(urlListLength)

	for i := 0; i < urlListLength; i++ {
		go func(i int) {
			defer wg.Done()

			r, err := processFile(urlList[i], i)
			if err != nil {
				t.Log(err)
				return
			}
			delta := r.calculateDelta()

			t.Logf("processing %d: %s", i, urlList[i])

			// keep critical images - otherwise remove
			if delta > 0 {
				t.Logf("file %s got bigger: %f%%", r.QuantitizedFilename, delta)
			} else {
				if !keepAllFiles {
					_ = os.Remove(r.OriginalFilename)
					_ = os.Remove(r.QuantitizedFilename)
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

func processFile(url string, index int) (r *processFileResult, err error) {

	// download image to local temp file and decode into image
	var fileIn, fileOut *os.File
	fileIn, err = ioutil.TempFile("/tmp", fmt.Sprintf("image%v-in-*.png", index))
	if err != nil {
		return
	}

	err = downloadFile(fileIn.Name(), url)
	if err != nil {
		return
	}

	inputImage, err := ioutil.ReadFile(fileIn.Name())
	if err != nil {
		return
	}

	img, _, err := image.Decode(bytes.NewReader(inputImage))
	if err != nil {
		return
	}

	// quantisation
	q, errQ := imagequant.New(img, 0, 0, 100, imagequant.DefaultSpeed)
	if errQ != nil {
		return
	}

	quantImg, errR := q.Run()
	if errR != nil {
		return
	}

	buff := bytes.NewBuffer([]byte{})

	// encode as png and write to temp file
	err = png.Encode(buff, quantImg)
	if err != nil {
		return
	}

	fileOut, err = ioutil.TempFile("/tmp", fmt.Sprintf("image%v-out-*.png", index))
	if err != nil {
		return
	}

	fw := bufio.NewWriter(fileOut)
	_, err = fw.Write(buff.Bytes())
	if err != nil {
		return
	}

	r = &processFileResult{OriginalSize: len(inputImage), QuantitizedSize: len(buff.Bytes()), OriginalFilename: fileIn.Name(), QuantitizedFilename: fileOut.Name()}

	return
}

func readUrls() (urls []string, err error) {
	f, err := os.OpenFile("./test/urls.txt", os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("open file error: %v", err)
		return
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		urls = append(urls, sc.Text())
	}
	if err = sc.Err(); err != nil {
		log.Fatalf("scan file error: %v", err)
		return
	}

	return

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

type processFileResult struct {
	OriginalSize        int
	OriginalFilename    string
	QuantitizedFilename string
	QuantitizedSize     int
}

func (m processFileResult) calculateDelta() (delta float64) {

	diff := float64(m.QuantitizedSize - m.OriginalSize)
	delta = (diff / float64(m.OriginalSize)) * 100
	return
}
