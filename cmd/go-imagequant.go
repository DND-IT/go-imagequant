package main

import (
    "flag"
    "fmt"
    "image"
    "image/gif"
    "image/jpeg"
    "image/png"
    "io/fs"
    "log"
    "os"
    "runtime"
    "runtime/pprof"
    "time"

    _ "golang.org/x/image/webp"

    quant "github.com/DND-IT/go-imagequant/imagequant"
)

var (
    initHeapAlloc uint64
    lastHeapAlloc uint64

    memStats runtime.MemStats
)

func main() {
    imageSrcPath := flag.String("src", "", "src image path")
    imageDstPath := flag.String("dst", "", "dst image path")
    speed := flag.Uint("speed", quant.DefaultSpeed, "speed to to use")
    gamma := flag.Float64("gamma", 0, "gamma")
    minQuality := flag.Uint("min.quality", 0, "min allowed quality (default 0)")
    maxQuality := flag.Uint("max.quality", 100, "min allowed quality")

    showLibImageQuantVersion := flag.Bool("showLibImageQuantVersion", false, "show lib image quant version and exit")

    checkMem := flag.Bool("checkmem", false, "repeat the image operation and print memory usage to detect memory leaks")
    checkMemIterations := flag.Uint("iterations", 1000, "how many iterations for checkmem should by used")
    CPUProfile := flag.String("CPUProfile", "", "write cpu profile to file")
    memProfile := flag.String("MEMProfile", "", "write memory profile to file")

    flag.Parse()

    var (
        memProfileFile *os.File
        cpuProfileFile *os.File
        err            error
    )

    if *CPUProfile != "" {
        cpuProfileFile, err = os.Create(*CPUProfile)
        if err != nil {
            log.Fatal("can not create CPU Profile", err)
        }
        defer cpuProfileFile.Close()
        err = pprof.StartCPUProfile(cpuProfileFile)
        if err != nil {
            log.Fatal("can not start CPU Profile", err)
        }
        defer pprof.StopCPUProfile()

    }

    if *showLibImageQuantVersion {
        fmt.Println(quant.GetLiqVersion())
        return
    }

    if *imageSrcPath == "" {
        fmt.Println("no src image")
        flag.PrintDefaults()
        return
    }

    if *imageDstPath == "" {
        fmt.Println("no dst image")
        flag.PrintDefaults()
        return
    }

    if !*checkMem {
        *checkMemIterations = 1
    }

    imageBuff, errBuff := Read(*imageSrcPath)

    if errBuff != nil {
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

    /* switch imgType {
       case "gif":

       } */
    var (
        qImg              image.Image
        memStats          runtime.MemStats
        newLineCounter    uint64
        newLineEveryCount = uint64(100)
    )

    runtime.ReadMemStats(&memStats)
    initHeapAlloc = memStats.HeapAlloc
    lastHeapAlloc = initHeapAlloc
    if *checkMem {
        printMemStats(true)
    }

    for i := uint(0); i < *checkMemIterations; i++ {
        q, qErr := quant.New(img, *gamma, *minQuality, *maxQuality, *speed)

        if qErr != nil {
            log.Println(qErr)
            return
        }

        qImg, err = q.Run()

        if err != nil {
            log.Println(err)
            return
        }

        if *checkMem {
            newLineCounter++
            if newLineCounter > newLineEveryCount {
                newLineCounter = 0
                printMemStats(true)
            }
        }

        //// unset image
        //if i < *checkMemIterations {
        //	q = nil
        //}

    }

    out, errEncode := Encode(qImg, imgType)
    if errEncode != nil {
        log.Println(errEncode)
    }

    _ = Write(*imageDstPath, out)

    if *memProfile != "" {
        memProfileFile, err = os.Create(*memProfile)
        if err != nil {
            log.Fatal("could not create memory profile: ", err)
        }
        defer memProfileFile.Close() // error handling omitted for example

        runtime.GC() // get up-to-date statistics
        if err := pprof.WriteHeapProfile(memProfileFile); err != nil {
            log.Fatal("could not write memory profile: ", err)
        }
    }

    if *checkMem {
        printMemStats(true)
        time.Sleep(5 * time.Second)
        printMemStats(true)
        time.Sleep(5 * time.Second)
        printMemStats(true)
        fmt.Print("\ndone\n")
        printMemStats(true)
        fmt.Print("\ndone\n")
    }

}

// Read an image file.
func Read(path string) ([]byte, error) {
    return os.ReadFile(path)
}

// Write an image file.
func Write(path string, data []byte) error {
    return os.WriteFile(path, data, fs.FileMode(0640))
}

func printMemStats(newLine bool) {
    runtime.ReadMemStats(&memStats)
    var direction = "↑"
    var directionInit = "↑"
    var delta, deltaInit uint64

    if lastHeapAlloc > memStats.HeapAlloc {
        direction = "↓"
        delta = lastHeapAlloc - memStats.HeapAlloc
    } else {
        delta = memStats.HeapAlloc - lastHeapAlloc
    }

    if initHeapAlloc > memStats.HeapAlloc {
        directionInit = "↓"
        deltaInit = initHeapAlloc - memStats.HeapAlloc
    } else {
        deltaInit = memStats.HeapAlloc - initHeapAlloc
    }

    fmt.Printf("%s init heap %12d - delta %12d - current heap %12d - frees: %d - %s delta to init: %12d - numgc: %6d - heap obj: %6d    \r",
        direction,
        initHeapAlloc,
        delta,
        memStats.HeapAlloc,
        memStats.Frees,
        directionInit,
        deltaInit,
        memStats.NumGC,
        memStats.HeapObjects)
    if newLine {
        fmt.Print("\n")
    }

    lastHeapAlloc = memStats.HeapAlloc
}
